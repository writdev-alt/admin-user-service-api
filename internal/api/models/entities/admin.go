package entities

import "github.com/google/uuid"

type Admin struct {
	ID     uint64    `gorm:"column:id;primaryKey" json:"id"`
	UUID   uuid.UUID `gorm:"column:uuid;type:char(36);not null" json:"uuid"`
	Status bool      `gorm:"column:status;type:tinyint(1);not null" json:"status"`
}

func (Admin) TableName() string { return "admins" }
