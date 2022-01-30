package app

import (
	"chat/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Middleware defines the contracts
type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type middleware struct{}

// NewMiddleware is a constructor
func NewMiddleware() Middleware {
	return &middleware{}
}

// AuthMiddleware checks if it has a jwt token and then requests the auth service to verify it for us
func (h *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" || authToken[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			c.Abort()
			return
		}
		client := &http.Client{}
		url := os.Getenv("AUTH_URL")
		circuitBeakerSettings := gobreaker.Settings{
			Name:          "auth",
			MaxRequests:   5,
			Interval:      time.Millisecond*50,
			Timeout:       time.Second,
			ReadyToTrip:   nil,
			OnStateChange: nil,
			IsSuccessful:  nil,
		}
		cb := gobreaker.NewCircuitBreaker(circuitBeakerSettings)
		cbResponse, err := cb.Execute(func() (interface{}, error) {
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, err
			}
			request.Header.Set("Authorization", authToken)
			response, err := client.Do(request)
			if err != nil {
				return nil, err
			}
			return response, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "some error occurred while making a request",
			})
			log.Println(err.Error())
			c.Abort()
			return
		}

		response, ok := cbResponse.(*http.Response)
		if !ok {
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError("some error occurred with cb"))
			return
		}

		if response.StatusCode == 200 {
			jwtToken := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)

			token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
			if err != nil {
				c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid token. I think"))
				return
			}
			var loggedID string
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				iss := claims["iss"]
				if f, ok := iss.(string); ok {
					loggedID = f
				}
			} else {
				c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid token. I think"))
				return
			}
			c.Set("loggedID", loggedID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not authorized",
			})
			c.Abort()
			return
		}
	}
}
