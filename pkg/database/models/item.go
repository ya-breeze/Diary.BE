package models

import (
	"github.com/google/uuid"
)

type Item struct {
	// gorm.Model

	Date  string `gorm:"index"`
	Title string
	Text  string
	IDs   IDList `gorm:"type:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

// func (u Item) FromDB() goserver.Item {
// 	return goserver.Item{
// 		Email:     u.Login,
// 		StartDate: u.StartDate,
// 	}
// }
