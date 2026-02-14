package entities

import "gorm.io/gorm"

type SchedulePrompt struct {
	gorm.Model
	Prompt string `json:"prompt"`
}
