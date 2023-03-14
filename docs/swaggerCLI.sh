# TUTORIAL
# https://lemoncode21.medium.com/how-to-add-swagger-in-golang-gin-6932e8076ec0

# install CLI
go get -u github.com/swaggo/swag/cmd/swag

go install github.com/swaggo/swag/cmd/swag@latest

# installation checking
swag -h

# instal Gin-Swagger
go get -u github.com/swaggo/gin-swagger

# generate swagger docs every time you modify doc string
swag init

# http://localhost:8080/swagger/index.html