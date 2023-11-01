# Stage 1: Build the Go microservice
# Use golang alpine latest
FROM golang:alpine
WORKDIR /usr/src/app

# Copy local code to the container image.
COPY . ./

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod download && go mod verify

RUN go build
RUN swag init

# RUN echo echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/main' >> /etc/apk/repositories
# RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/community' >> /etc/apk/repositories
# RUN apk add --no-cache mongodb-tools mongodb

# Run the Go service on startup
CMD ["go", "run", "ccu"]
