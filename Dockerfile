FROM golang:1.17-alpine

ARG ROOT_DIR=$GOPATH/github.com/Airbenders-490/chat

# Current working directory inside the container
WORKDIR "$ROOT_DIR"

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build (compile code with dependencies) the app then leave the result in the output main directory
RUN go build -o main .

# Move to /bin directory as the place for resulting binary
WORKDIR /bin

COPY static .

# Copy binary from rootDir/main to current /bin folder
RUN cp /"$ROOT_DIR"/main .

# This container exposes port 8080 to the outside world & listens on this port
EXPOSE 8080

COPY ./wait-for-cassandra.sh /wait-for-cassandra.sh
RUN chmod +x /wait-for-cassandra.sh

# Command to run when starting the container
CMD [ "/wait-for-cassandra.sh" ]
