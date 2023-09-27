# Login API Service
This service is responsible for login logic. Handles login events and account creation

API Documentation is created using the [Swagger](https://swagger.io/). The url for the Swagger UI is on the same port as the `APP_PORT` in the `.env` file at 

```
http://<ip>:<APP_PORT>/swagger/index.html
``` 

If you are running locally it would be at [http://127.0.0.1:8088/swagger/index.html](http://127.0.0.1:8088/swagger/index.html)

# How to run

## Create
### `.env` Creation
Create a `.env` file
```bash
touch .env
```
Edit the `.env` file with any editor of your choice
```bash
vim .env
```

### `.env` Template
```
APP_PORT=:8088 // Standard port for this microservice
LOG_LEVEL=trace
METHOD_LOGGING=false
```
> Additional fields will also be required in the `.env` file to run the microservice successfully. Here is a basic template of the `.env`. Customize to your liking. This template will change as the microservice matures and implements new features.

## Build
```bash
go build
```
## Run
```bash
go run ccu
```
or if you dont want to build
```bash
go run main.go
```
## (Optional) Update package checksums and download dependencies
```
go mod tidy
``` 
