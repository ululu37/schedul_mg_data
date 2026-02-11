package usecase

import (
	"fmt"
	"scadulDataMono/domain/entities"
	"testing"
)

func TestCan(t *testing.T) {
	tests := []struct {
		name       string
		teachers   []entities.Teacher
		classrooms []entities.Classroom
		wantErr    bool
		errMsg     string
	}{
		// {
		// 	name: "Average hour too high",
		// 	teachers: []entities.Teacher{
		// 		{ID: 1},
		// 	},
		// 	classrooms: []entities.Classroom{
		// 		{
		// 			PreCurriculum: entities.PreCurriculum{
		// 				SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
		// 					{SubjectID: 1, Credit: 50}, // More than 40 hours per teacher
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: true,
		// 	errMsg:  "ไม่สามารถจัดตารางได้",
		// },
		{
			name: "Valid scheduling scenario",
			teachers: []entities.Teacher{
				{
					ID: 1,
					MySubject: []entities.TeacherMySubject{
						{SubjectID: 1, Preference: 4},
						{SubjectID: 2, Preference: 8},
						{SubjectID: 3, Preference: 8},
					},
				},
				{
					ID: 2,
					MySubject: []entities.TeacherMySubject{
						{SubjectID: 1, Preference: 4},
						{SubjectID: 2, Preference: 8},
					},
				},
				{
					ID: 3,
					MySubject: []entities.TeacherMySubject{
						{SubjectID: 1, Preference: 5},
						{SubjectID: 3, Preference: 8},
					},
				},
			},
			classrooms: []entities.Classroom{
				{
					ID: 1,
					PreCurriculum: entities.PreCurriculum{
						SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
							{SubjectID: 1, Credit: 3},
							{SubjectID: 2, Credit: 3},
						},
					},
				},
				{
					ID: 2,
					PreCurriculum: entities.PreCurriculum{
						SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
							{SubjectID: 1, Credit: 3},
							{SubjectID: 2, Credit: 2},
						},
					},
				},
				// {
				// 	ID: 3,
				// 	PreCurriculum: entities.PreCurriculum{
				// 		SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
				// 			{SubjectID: 1, Credit: 3},
				// 			{SubjectID: 2, Credit: 2},
				// 			{SubjectID: 3, Credit: 8},
				// 		},
				// 	},
				// },
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Can(tt.teachers, tt.classrooms)
			if err != nil {
				t.Errorf("Can() error = %v", err)
			}
			fmt.Println(got)
		})
	}
}
