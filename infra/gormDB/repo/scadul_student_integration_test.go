package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestScadulStudentRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Subject{}, &entities.Teacher{}, &entities.Classroom{}, &entities.ScadulStudent{}, &entities.SubjectInScadulStudent{}))

	repo := &ScadulStudentRepo{DB: db}

	// create teacher and subject and classroom
	teacher := &entities.Teacher{Name: "Mr T"}
	require.NoError(t, db.Create(teacher).Error)
	subject := &entities.Subject{Name: "Math"}
	require.NoError(t, db.Create(subject).Error)
	classroom := &entities.Classroom{TeacherID: teacher.ID, Year: 1}
	require.NoError(t, db.Create(classroom).Error)

	// create scadul student
	s := &entities.ScadulStudent{ClassroomID: classroom.ID, UseIn: "2026/T1"}
	id, err := repo.Create(s)
	require.NoError(t, err)
	require.NotZero(t, id)

	// add subject entries
	require.NoError(t, repo.AddSubjects([]entities.SubjectInScadulStudent{{ScadulStudentID: id, TeacherID: teacher.ID, SubjectID: subject.ID, Order: 1}}))

	// get and verify
	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "2026/T1", got.UseIn)
	require.GreaterOrEqual(t, len(got.SubjectInScadulStudent), 1)
	require.Equal(t, teacher.ID, got.SubjectInScadulStudent[0].TeacherID)

	// listing by useIn
	list, total, err := repo.Listing("2026/T1", 1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int64(1))
	require.GreaterOrEqual(t, len(list), 1)

	// update
	up, err := repo.Update(id, &entities.ScadulStudent{ClassroomID: classroom.ID, UseIn: "2026/T2"})
	require.NoError(t, err)
	require.Equal(t, "2026/T2", up.UseIn)

	// remove subject
	require.NoError(t, repo.RemoveSubjects([]uint{got.SubjectInScadulStudent[0].ID}))
	g2, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Len(t, g2.SubjectInScadulStudent, 0)

	// delete
	require.NoError(t, repo.Delete(id))
	err = db.First(&entities.ScadulStudent{}, id).Error
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
