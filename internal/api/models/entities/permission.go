package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID  `gorm:"column:uuid;type:char(36);not null;uniqueIndex:idx_uuid_unique" json:"uuid"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	GuardName   string     `gorm:"column:guard_name;type:varchar(255);not null;default:web" json:"guardName"`
	Description *string    `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
	CreatedBy   *uint64    `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy   *uint64    `gorm:"column:updated_by" json:"updatedBy"`
	DeletedBy   *uint64    `gorm:"column:deleted_by" json:"deletedBy"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt,omitempty"`
}

func (Permission) TableName() string { return "permissions" }

func (m *Permission) BeforeCreate(*gorm.DB) error {
	if m.UUID == (uuid.UUID{}) {
		m.UUID = uuid.New()
	}
	now := time.Now()
	if m.CreatedAt == nil {
		m.CreatedAt = &now
	}
	if m.UpdatedAt == nil {
		m.UpdatedAt = &now
	}
	return nil
}

func (m *Permission) BeforeUpdate(*gorm.DB) error {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}
