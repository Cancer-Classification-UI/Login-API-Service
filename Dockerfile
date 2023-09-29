# Docker container for the login-api service
# Will download needed dependencies and build the service and run it

FROM golang:alpine
WORKDIR /usr/src/app

COPY go.mod go.sum .env ./

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod download && go mod verify

COPY . .
RUN go build
RUN swag init


CMD ["go","run","ccu"]
