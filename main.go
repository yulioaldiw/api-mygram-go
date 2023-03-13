package main

import (
	"log"
	"os"

	commentDelivery "api-mygram-go/comment/delivery/http"
	commentRepository "api-mygram-go/comment/repository/postgres"
	commentUseCase "api-mygram-go/comment/usecase"
	"api-mygram-go/config/database"
	photoDelivery "api-mygram-go/photo/delivery/http"
	photoRepository "api-mygram-go/photo/repository/postgres"
	photoUseCase "api-mygram-go/photo/usecase"
	socialMediaDelivery "api-mygram-go/socialmedia/delivery/http"
	socialMediaRepository "api-mygram-go/socialmedia/repository/postgres"
	socialMediaUseCase "api-mygram-go/socialmedia/usecase"
	userDelivery "api-mygram-go/user/delivery/http"
	userRepository "api-mygram-go/user/repository/postgres"
	userUseCase "api-mygram-go/user/usecase"

	_ "api-mygram-go/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title MyGram API
// @version 1.0
// @description MyGram is a free photo sharing app written in Go. People can share, view, and comment photos by everyone. Anyone can create an account by registering an email address and creating a username.
// @termOfService http://swagger.io/terms/
// @contact.name yulioaldiw
// @contact.email yulioaldiw@gmail.com
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used
func main() {
	if err := godotenv.Load("../api-mygram-go/config/env/.env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	db := database.StartDB()

	routers := gin.Default()

	routers.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	})

	userRepository := userRepository.NewUserRepository(db)
	userUseCase := userUseCase.NewUserUseCase(userRepository)

	userDelivery.NewUserHandler(routers, userUseCase)

	photoRepository := photoRepository.NewPhotoRepository(db)
	photoUseCase := photoUseCase.NewPhotoUseCase(photoRepository)

	photoDelivery.NewPhotoHandler(routers, photoUseCase)

	commentRepository := commentRepository.NewCommentRepository(db)
	commentUseCase := commentUseCase.NewCommentUseCase(commentRepository)

	commentDelivery.NewCommentHandler(routers, commentUseCase, photoUseCase)

	socialMediaRepository := socialMediaRepository.NewSocialMediaRepository(db)
	socialMediaUseCase := socialMediaUseCase.NewSocialMediaUseCase(socialMediaRepository)

	socialMediaDelivery.NewSocialMediaHandler(routers, socialMediaUseCase)

	routers.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := godotenv.Load("../api-mygram-go/config/env/.env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	port := os.Getenv("PORT")

	if len(os.Args) > 1 {
		reqPort := os.Args[1]

		if reqPort != "" {
			port = reqPort
		}
	}

	if port == "" {
		port = "8080"
	}

	routers.Run(":" + port)
}
