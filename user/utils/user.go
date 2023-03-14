package utils

import (
	"time"
)

type RegisterUser struct {
	Age      uint   `json:"age" example:"8"`
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secret"`
	Email    string `json:"email" example:"johndoe@example.com"`
}

type RegisteredUser struct {
	ID       string `json:"id" example:"the user id generated here"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Username string `json:"username" example:"johndoe"`
	Age      uint   `json:"age" example:"8"`
}

type ResponseDataRegisteredUser struct {
	Status string         `json:"status" example:"success"`
	Data   RegisteredUser `json:"data"`
}

type LoginUser struct {
	Email    string `json:"email" example:"johndoe@example.com"`
	Password string `json:"password" example:"secret"`
}

type LoggedinUser struct {
	Username string `json:"username" example:"johndoe"`
	Token    string `json:"token" example:"the token generated here"`
}

type ResponseDataLoggedinUser struct {
	Status string       `json:"status" example:"success"`
	Data   LoggedinUser `json:"data"`
}

type Users struct {
	ID       string `json:"id" example:"the user id generated here"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Username string `json:"username" example:"johndoe"`
	Age      uint   `json:"age" example:"8"`
}

type FetchedUsers struct {
	Users interface{} `json:"users"`
}

type ResponseGetAllUsers struct {
	Status string `json:"status" example:"success"`
	Data   Users  `json:"data"`
}

type UpdateUser struct {
	Email    string `json:"email" example:"newjohndoe@example.com"`
	Username string `json:"username" example:"newjohndoe"`
}

type UpdatedUser struct {
	ID        string     `json:"id" example:"here is the generated user id"`
	Email     string     `json:"email" example:"newjohndoe@example.com"`
	Username  string     `json:"username" example:"newjohndoe"`
	Age       uint       `json:"age" example:"8"`
	UpdatedAt *time.Time `json:"updated_at" example:"the updated at generated here"`
}

type ResponseDataUpdatedUser struct {
	Status string      `json:"status" example:"success"`
	Data   UpdatedUser `json:"data"`
}

type ResponseMessageDeletedUser struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your account has been successfully deleted"`
}

type ResponseMessage struct {
	Status string `json:"status" example:"fail"`
	Data   string `json:"data" example:"the error explained here"`
}
