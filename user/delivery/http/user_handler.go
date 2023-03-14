package delivery

import (
	"api-mygram-go/domain"
	"api-mygram-go/helpers"
	"api-mygram-go/user/delivery/http/middleware"
	"api-mygram-go/user/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(routers *gin.Engine, userUseCase domain.UserUseCase) {
	handler := &userHandler{userUseCase}

	router := routers.Group("/users")
	{
		router.POST("/register", handler.Register)
		router.POST("/login", handler.Login)
		router.GET("", middleware.Authentication(), handler.GetAllUsers)
		router.PUT("", middleware.Authentication(), handler.Update)
		router.DELETE("", middleware.Authentication(), handler.Delete)
	}
}

// Register godoc
// @Summary			Register a user
// @Description		create and store a user
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json	body		utils.RegisterUser	true	"Register User"
// @Success			201		{object}	utils.ResponseDataRegisteredUser
// @Failure			400  	{object}	utils.ResponseMessage
// @Failure			409  	{object}	utils.ResponseMessage
// @Router			/users/register	[post]
func (handler *userHandler) Register(ctx *gin.Context) {
	var (
		user domain.User
		err  error
	)

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	if err = handler.userUseCase.Register(ctx.Request.Context(), &user); err != nil {
		if strings.Contains(err.Error(), "idx_users_username") {
			ctx.AbortWithStatusJSON(http.StatusConflict, helpers.ResponseMessage{
				Status:  "fail",
				Message: "the username you entered has been used",
			})

			return
		}

		if strings.Contains(err.Error(), "idx_users_email") {
			ctx.AbortWithStatusJSON(http.StatusConflict, helpers.ResponseMessage{
				Status:  "fail",
				Message: "the email you entered has been used",
			})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, helpers.ResponseData{
		Status: "success",
		Data: utils.RegisteredUser{
			Age:      user.Age,
			Email:    user.Email,
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// Login godoc
// @Summary			Login a user
// @Description		Authentication a user and retrieve a token
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json	body		utils.LoginUser	true	"Login User"
// @Success			200		{object}	utils.ResponseDataLoggedinUser
// @Failure			400		{object}	utils.ResponseMessage
// @Failure			401		{object}	utils.ResponseMessage
// @Router			/users/login		[post]
func (handler *userHandler) Login(ctx *gin.Context) {
	var (
		user  domain.User
		err   error
		token string
	)

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: fmt.Sprint("bad payload, enter the valid payload; error: ", err.Error()),
		})

		return
	}

	if err = handler.userUseCase.Login(ctx.Request.Context(), &user); err != nil {
		if strings.Contains(err.Error(), "the credential you entered are wrong") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
				Status:  "unauthenticated",
				Message: err.Error(),
			})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "unauthenticated",
			Message: err.Error(),
		})

		return
	}

	if token = helpers.GenerateToken(user.ID, user.Email); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "unauthenticated",
			"message": err.Error(),
		})

		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 300, "", "", false, true)

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: utils.LoggedinUser{
			Username: user.Email,
			Token:    token,
		},
	})
}

// Get Users godoc
// @Summary			Get all users
// @Description		Get all users data with authenticated user
// @Tags			users
// @Accept			json
// @Produce			json
// @Success			200		{object}	utils.ResponseGetAllUsers
// @Failure			400		{object}	utils.ResponseMessage
// @Failure			401		{object}	utils.ResponseMessage
// @Security    	Bearer
// @Router			/users	[get]
func (handler *userHandler) GetAllUsers(ctx *gin.Context) {
	var (
		users []domain.User
		err   error
	)

	if err = handler.userUseCase.GetAllUsers(ctx.Request.Context(), &users); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: utils.FetchedUsers{
			// Users: users,
			Users: users,
		},
	})
}

// Update godoc
// @Summary			Update a user
// @Description		Update a user with authentication user
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			json		body		utils.UpdateUser   true  "Update User"
// @Success			200			{object}  	utils.ResponseDataUpdatedUser
// @Failure			400			{object}	utils.ResponseMessage
// @Failure			401			{object}	utils.ResponseMessage
// @Failure			409			{object}	utils.ResponseMessage
// @Security		Bearer
// @Router			/users	[put]
func (handler *userHandler) Update(ctx *gin.Context) {
	var (
		user domain.User
		err  error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedUser := domain.User{
		Username:        user.Username,
		Email:           user.Email,
		Age:             user.Age,
		ProfileImageUrl: user.ProfileImageUrl,
	}

	if user, err = handler.userUseCase.Update(ctx.Request.Context(), updatedUser, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: utils.UpdatedUser{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Age:       user.Age,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Delete godoc
// @Summary			Delete a user
// @Description	Delete a user with authentication user
// @Tags				users
// @Accept			json
// @Produce			json
// @Success			200			{object}	utils.ResponseMessageDeletedUser
// @Failure			400			{object}	utils.ResponseMessage
// @Failure			401			{object}	utils.ResponseMessage
// @Failure			404			{object}	utils.ResponseMessage
// @Security		Bearer
// @Router			/users	[delete]
func (handler *userHandler) Delete(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err := handler.userUseCase.Delete(ctx, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: "account not found",
		})

		return
	}

	ctx.JSON(
		http.StatusOK,
		helpers.ResponseMessage{
			Status:  "success",
			Message: "your account has been successfully deleted",
		},
	)
}
