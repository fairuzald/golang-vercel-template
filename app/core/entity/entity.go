package entity

import "time"

type Entity interface {
	GetID() string
	SetID(id string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

type BaseEntity struct {
	ID        string    `json:"id" firestore:"id"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}

func (e *BaseEntity) GetID() string {
	return e.ID
}

func (e *BaseEntity) SetID(id string) {
	e.ID = id
}

func (e *BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e *BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

func (e *BaseEntity) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}
