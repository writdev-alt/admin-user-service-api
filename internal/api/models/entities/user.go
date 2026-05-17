package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID             uuid.UUID  `gorm:"column:uuid;type:char(36);not null;uniqueIndex:idx_uuid_unique" json:"uuid"`
	Username         string     `gorm:"column:username;type:varchar(255);not null" json:"username"`
	Phone            *string    `gorm:"column:phone;type:varchar(255)" json:"phone"`
	Country          *string    `gorm:"column:country;type:varchar(255)" json:"country"`
	Email            string     `gorm:"column:email;type:varchar(255);not null" json:"email"`
	EmailVerifiedAt  *time.Time `gorm:"column:email_verified_at;type:timestamp" json:"emailVerifiedAt,omitempty"`
	TwoFactorEnabled bool       `gorm:"column:two_factor_enabled;type:tinyint(1);not null;default:0" json:"twoFactorEnabled"`
	Status           bool       `gorm:"column:status;type:tinyint(1);not null;default:1" json:"status"`
	Password         string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	CreatedBy        uint64     `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy        uint64     `gorm:"column:updated_by" json:"updatedBy"`
	DeletedBy        uint64     `gorm:"column:deleted_by" json:"deletedBy"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp" json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt,omitempty"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt,omitempty"`
}

func (User) TableName() string { return "users" }

func (m *User) BeforeCreate(*gorm.DB) error {
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

func (m *User) BeforeUpdate(*gorm.DB) error {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}
