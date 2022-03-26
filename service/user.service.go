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
	CheckEmailAvailability(request dto.CheckEmailRequest) (bool, error)
	SaveAvatar(id int, fileLocation string) (entity.User, error)
	GetUserByID(userID int) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
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

func (s *userService) CheckEmailAvailability(request dto.CheckEmailRequest) (bool, error) {
	email := request.Email

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	// Check if user is available through default value of ID
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) SaveAvatar(id int, fileLocation string) (entity.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return user, err
	}
	user.AvatarFileName = fileLocation
	updatedUser, err := s.userRepository.Update(user)
	if err != nil {
		return user, err
	}
	return updatedUser, nil
}

func (s *userService) GetUserByID(userID int) (entity.User, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user with the correspond ID")
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]entity.User, error) {
	users, err := s.userRepository.FindAll()
	if err != nil {
		return users, err
	}
	return users, nil
}
