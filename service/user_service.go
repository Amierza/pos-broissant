package service

import (
	"context"
	"regexp"
	"sync"

	"github.com/Amierza/pos-broissant/dto"
	"github.com/Amierza/pos-broissant/entity"
	"github.com/Amierza/pos-broissant/repository"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error)
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

	if !isValidEmail(req.Email) {
		return dto.UserRegisterResponse{}, dto.ErrInvalidEmail
	}

	_, flag, err := s.userRepo.CheckEmailOrPhoneNumber(ctx, nil, req.Email, req.PhoneNumber)
	if err == nil || flag {
		return dto.UserRegisterResponse{}, dto.ErrEmailOrPhoneNumberAlreadyExists
	}

	if len(req.Password) < 8 {
		return dto.UserRegisterResponse{}, dto.ErrPasswordLessThanEight
	}

	if len(req.Pin) != 6 {
		return dto.UserRegisterResponse{}, dto.ErrLengthPinMustBeSix
	}

	user := entity.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Pin:         req.Pin,
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
		Pin:         userReg.Pin,
	}, nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
