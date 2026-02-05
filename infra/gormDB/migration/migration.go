package main

import (
	"scadulDataMono/config"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/gormDB"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.Auth{},
		&entities.Teacher{},
		&entities.Student{},

		&entities.Subject{},
		&entities.Term{},
		&entities.Curriculum{},
		&entities.SubjectInCurriculum{},

		&entities.TeacherMySubject{},

		&entities.Classroom{},
		&entities.ScadulStudent{},
		&entities.ScadulTeacher{},

		&entities.SubjectInScadulStudent{},
		&entities.SubjectInScadulTeacher{},

		
		&entities.PreCurriculum{},
		&entities.SubjectInPreCurriculum{},
	)
}
func main() {
	conf := config.LoadConfig()
	db := gormDB.NewPostgresDatabase(&conf.Database)
	Migrate(db.Connect())
}
