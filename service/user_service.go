package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/Amierza/pos-broissant/dto"
	"github.com/Amierza/pos-broissant/entity"
	"github.com/Amierza/pos-broissant/helpers"
	"github.com/Amierza/pos-broissant/repository"
	"github.com/go-playground/validator/v10"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error)
		LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
	}
	userService struct {
		userRepo   repository.UserRepository
		jwtService JWTService
	}
)

var (
	mu sync.Mutex
)

const (
	LOCAL_URL          = "http://localhost:8080"
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

func NewUserService(userRepo repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
		return dto.UserRegisterResponse{}, fmt.Errorf("validation errors: %v", errorMessages)
	}

	_, flag, err := s.userRepo.CheckEmailOrPhoneNumber(ctx, nil, req.Email, req.PhoneNumber)
	if err == nil || flag {
		return dto.UserRegisterResponse{}, dto.ErrEmailOrPhoneNumberAlreadyExists
	}

	user := entity.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserRegisterResponse{}, dto.ErrRegisterUser
	}

	return dto.UserRegisterResponse{
		ID:          userReg.ID.String(),
		FirstName:   userReg.FirstName,
		LastName:    userReg.LastName,
		Email:       userReg.Email,
		Password:    userReg.Password,
		PhoneNumber: userReg.PhoneNumber,
		Role:        userReg.Role,
	}, nil
}

func (s *userService) LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error())
		}
		return dto.UserLoginResponse{}, fmt.Errorf("validation errors: %v", errorMessages)
	}

	user, flag, err := s.userRepo.CheckEmailOrPhoneNumber(ctx, nil, req.Email, req.Password)
	if !flag || err != nil {
		return dto.UserLoginResponse{}, dto.ErrEmailNotFound
	}

	checkPass, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if !checkPass || err != nil {
		return dto.UserLoginResponse{}, dto.ErrPasswordDoesntMatch
	}

	accessToken, refreshToken, err := s.jwtService.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return dto.UserLoginResponse{}, dto.ErrGenerateToken
	}

	return dto.UserLoginResponse{
		ID:           user.ID.String(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		PhoneNumber:  user.PhoneNumber,
		Pin:          user.Pin,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPagination, err := c.userRepo.GetAllUserWithPaginationRepo(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.AllUserResponse
	for _, user := range dataWithPagination.Users {
		data := dto.AllUserResponse{
			ID:          user.ID.String(),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
			Pin:         user.Pin,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPagination.Page,
			PerPage: dataWithPagination.PerPage,
			Count:   dataWithPagination.Count,
			MaxPage: dataWithPagination.MaxPage,
		},
	}, nil
}
