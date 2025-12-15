package model

import "time"

type UserRecord struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	UID       string    `gorm:"column:uid" json:"uid"`
	Name      string    `gorm:"column:name" json:"name"`
	Concat    string    `gorm:"column:concat" json:"concat"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
