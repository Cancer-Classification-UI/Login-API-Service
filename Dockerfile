
FROM golang:alpine
WORKDIR /usr/src/app

# Copy local code to the container image.
COPY . ./

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod download && go mod verify

RUN go build
RUN swag init

# Run the the service on startup
CMD ["go","run","ccu"]
