package entity

import (
	"github.com/Amierza/pos-broissant/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Pin         string    `gorm:"default:null" json:"pin"`
	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	u.ID = uuid.New()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Pin, err = helpers.HashPin(u.Pin)
	if err != nil {
		return err
	}
	return nil
}
