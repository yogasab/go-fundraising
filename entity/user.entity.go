package entity

import "time"

type User struct {
	ID             int64     `gorm:"primary_key:auto_increment"`
	Name           string    `gorm:"type:varchar(100)"`
	Email          string    `gorm:"type:varchar(100);unique"`
	Password       string    `gorm:"type:varchar(100)"`
	Occupation     string    `gorm:"type:varchar(100)"`
	AvatarFileName string    `gorm:"type:varchar(100)"`
	Token          string    `gorm:"type:varchar(100)"`
	Role           string    `gorm:"type:varchar(100)"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
