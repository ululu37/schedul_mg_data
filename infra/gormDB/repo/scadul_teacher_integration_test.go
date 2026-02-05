package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestScadulTeacherRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Subject{}, &entities.Teacher{}, &entities.ScadulTeacher{}, &entities.SubjectInScadulTeacher{}))

	repo := &ScadulTeacherRepo{DB: db}

	// create teacher and subject
	teacher := &entities.Teacher{Name: "Mr T"}
	require.NoError(t, db.Create(teacher).Error)
	subject := &entities.Subject{Name: "Math"}
	require.NoError(t, db.Create(subject).Error)

	// create scadul teacher
	s := &entities.ScadulTeacher{TeacherID: teacher.ID, UseIn: "2026/T1"}
	id, err := repo.Create(s)
	require.NoError(t, err)
	require.NotZero(t, id)

	// add subject entries
	require.NoError(t, repo.AddSubjects([]entities.SubjectInScadulTeacher{{ScadulTeacherID: id, TeacherID: teacher.ID, SubjectID: subject.ID, Order: 1}}))

	// get and verify
	got, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "2026/T1", got.UseIn)
	require.GreaterOrEqual(t, len(got.SubjectInScadulTeacher), 1)
	require.Equal(t, teacher.ID, got.SubjectInScadulTeacher[0].TeacherID)

	// listing by useIn
	list, total, err := repo.Listing("2026/T1", 1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int64(1))
	require.GreaterOrEqual(t, len(list), 1)

	// update
	up, err := repo.Update(id, &entities.ScadulTeacher{TeacherID: teacher.ID, UseIn: "2026/T2"})
	require.NoError(t, err)
	require.Equal(t, "2026/T2", up.UseIn)

	// remove subject
	require.NoError(t, repo.RemoveSubjects([]uint{got.SubjectInScadulTeacher[0].ID}))
	g2, err := repo.GetByID(id)
	require.NoError(t, err)
	require.Len(t, g2.SubjectInScadulTeacher, 0)

	// delete
	require.NoError(t, repo.Delete(id))
	err = db.First(&entities.ScadulTeacher{}, id).Error
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
