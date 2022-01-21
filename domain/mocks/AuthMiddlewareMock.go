package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MiddlewareMock mocks the middleware for testing
type MiddlewareMock struct {
	mock.Mock
}

// AuthMiddleware doesn't check if the user has credentials. It simply assigns the id from the header
// to the context that can be used further for evaluation. Without being authenticated
func (m *MiddlewareMock) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("id")
		c.Set("loggedID", id)
		c.Next()
	}
}