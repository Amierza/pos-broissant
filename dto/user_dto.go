package dto

import (
	"errors"

	"github.com/Amierza/pos-broissant/entity"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER      = "failed register user"
	MESSAGE_FAILED_LOGIN_USER         = "failed login user"
	MESSAGE_FAILED_GET_LIST_USER      = "failed get list user"
	MESSAGE_FAILED_PROSES_REQUEST     = "failed proses request"
	MESSAGE_FAILED_TOKEN_NOT_FOUND    = "failed token not found"

	// Success
	MESSAGE_SUCCESS_GET_DATA_FROM_BODY = "success get data from body"
	MESSAGE_SUCCESS_REGISTER_USER      = "success register user"
	MESSAGE_SUCCESS_LOGIN_USER         = "success login user"
	MESSAGE_SUCCESS_GET_LIST_USER      = "success get list user"
)

var (
	ErrEmailOrPhoneNumberAlreadyExists = errors.New("email or phone number is already exists")
	ErrInvalidEmail                    = errors.New("invalid email")
	ErrPasswordLessThanEight           = errors.New("length password must be more than or equal to eight")
	ErrLengthPinMustBeSix              = errors.New("length pin must be six")
	ErrRegisterUser                    = errors.New("failed to register the user")
	ErrEmailNotFound                   = errors.New("email not found")
	ErrPasswordDoesntMatch             = errors.New("password doesnt match")
	ErrGenerateToken                   = errors.New("failed to generate token")
)

type (
	UserRegisterRequest struct {
		FirstName   string `json:"first_name" form:"first_name" validate:"required"`
		LastName    string `json:"last_name" form:"last_name" validate:"required"`
		Email       string `json:"email" form:"email" validate:"required,email"`
		Password    string `json:"password" form:"password" validate:"required,min=8"`
		PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	}

	UserRegisterResponse struct {
		ID          string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
		entity.Timestamp
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	UserLoginResponse struct {
		ID           string `json:"user_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		PhoneNumber  string `json:"phone_number"`
		Pin          string `json:"pin"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		entity.Timestamp
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User
		PaginationResponse
	}

	AllUserResponse struct {
		ID          string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
		entity.Timestamp
	}

	UserPaginationResponse struct {
		Data []AllUserResponse `json:"data"`
		PaginationResponse
	}
)
