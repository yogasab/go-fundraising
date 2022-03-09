package service

import (
	"errors"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService interface {
	RegisterUser(request dto.RegisterRequest) (entity.User, error)
	LoginUser(request dto.LoginRequest) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) RegisterUser(request dto.RegisterRequest) (entity.User, error) {
	user := entity.User{}
	user.Name = request.Name
	user.Email = request.Email
	user.Occupation = request.Occupation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
		return user, nil
	}
	user.Password = string(hashedPassword)
	user.Role = "user"

	newUser, err := s.userRepository.Save(user)
	if err != nil {
		log.Println(err.Error())
		return user, nil
	}
	return newUser, nil
}

func (s *userService) LoginUser(request dto.LoginRequest) (entity.User, error) {
	email := request.Email
	password := request.Password

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return user, errors.New("User with correspond email is not registered")
	}
	if user.ID == 0 {
		return user, errors.New("User with correspond email is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("Invalid credentials check your email or password")
	}

	return user, nil
}
