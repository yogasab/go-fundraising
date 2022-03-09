package repository

import (
	"go-fundraising/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	Save(user entity.User) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(user entity.User) (entity.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	// Which table we query on
	var user entity.User
	err := r.db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}
