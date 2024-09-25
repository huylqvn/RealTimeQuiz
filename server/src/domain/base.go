package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status int8

const (
	StatusInactive Status = 0
	StatusActive   Status = 1
)

type BaseModel struct {
	ID        uint64     `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = &now
	m.UpdatedAt = &now
	return
}

type BaseUUIDModel struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m *BaseUUIDModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return
}

func generateUUID() string {
	return uuid.New().String()
}

func (m *BaseUUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.ID = generateUUID()
	m.CreatedAt = &now
	m.UpdatedAt = &now
	return
}

type SuccessResponse struct {
	Message string `json:"message"`
}

const (
	ColUpdatedAt = "updated_at"
	ColCreatedAt = "created_at"
	ColID        = "id"
)
