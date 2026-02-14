package repo

import (
	"scadulDataMono/domain/entities"

	"gorm.io/gorm"
)

type SchedulePromptRepo struct {
	DB *gorm.DB
}

func NewSchedulePromptRepo(db *gorm.DB) *SchedulePromptRepo {
	return &SchedulePromptRepo{DB: db}
}

func (r *SchedulePromptRepo) Create(prompt *entities.SchedulePrompt) (*entities.SchedulePrompt, error) {
	if err := r.DB.Create(prompt).Error; err != nil {
		return nil, err
	}
	return prompt, nil
}

func (r *SchedulePromptRepo) GetAll() ([]entities.SchedulePrompt, error) {
	var prompts []entities.SchedulePrompt
	if err := r.DB.Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}

func (r *SchedulePromptRepo) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entities.SchedulePrompt{}).Error
}
