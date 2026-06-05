package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	IsTemp       bool      `json:"is_temp" db:"is_temp"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Task struct {
	ID              string    `json:"id" db:"id"`
	TextKey         string    `json:"text_key" db:"text_key"`
	DurationSeconds int       `json:"duration_seconds" db:"duration_seconds"`
	Category        string    `json:"category" db:"category"`
	IsPreset        bool      `json:"is_preset" db:"is_preset"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type UserTaskPool struct {
	UserID  string    `json:"user_id" db:"user_id"`
	TaskID  string    `json:"task_id" db:"task_id"`
	AddedAt time.Time `json:"added_at" db:"added_at"`
}

type TaskHistory struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	TaskID    string    `json:"task_id" db:"task_id"`
	StartedAt time.Time `json:"started_at" db:"started_at"`
	EndedAt   time.Time `json:"ended_at" db:"ended_at"`
	Outcome   string    `json:"outcome" db:"outcome"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserSettings struct {
	UserID          string    `json:"user_id" db:"user_id"`
	CustomTaskRatio int       `json:"custom_task_ratio" db:"custom_task_ratio"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type AdminUser struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type TaskCategory struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	LabelKey string `json:"label_key" db:"label_key"`
	IsActive bool   `json:"is_active" db:"is_active"`
}
