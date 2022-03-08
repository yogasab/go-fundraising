package repository

import (
	"go-fundraising/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	Save(user entity.User) (entity.User, error)
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
