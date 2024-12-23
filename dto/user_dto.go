package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER      = "failed register user"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER = "success register user"
)

var (
	ErrEmailOrPhoneNumberAlreadyExists = errors.New("email or phone number is already exists")
	ErrInvalidEmail                    = errors.New("invalid email")
	ErrPasswordLessThanEight           = errors.New("length password must be more than or equal to eight")
	ErrLengthPinMustBeSix              = errors.New("length pin must be six")
	ErrRegisterUser                    = errors.New("failed to register the user")
)

type (
	UserRegisterRequest struct {
		FirstName   string `json:"first_name" form:"first_name"`
		LastName    string `json:"last_name" form:"last_name"`
		Email       string `json:"email" form:"email"`
		Password    string `json:"password" form:"password"`
		PhoneNumber string `json:"phone_number" form:"phone_number"`
		Pin         string `json:"pin" form:"pin"`
	}

	UserRegisterResponse struct {
		ID          string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
	}
)
