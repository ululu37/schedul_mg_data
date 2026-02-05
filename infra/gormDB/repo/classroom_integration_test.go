package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestClassroomRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Teacher{}, &entities.Classroom{}, &entities.Student{}))

	repo := &ClassroomRepo{DB: db}

	// create teachers
	teacher1 := &entities.Teacher{Name: "Mr T"}
	require.NoError(t, db.Create(teacher1).Error)
	teacher2 := &entities.Teacher{Name: "Ms S"}
	require.NoError(t, db.Create(teacher2).Error)

	// Create classroom
	c := &entities.Classroom{TeacherID: teacher1.ID, Year: 1}
	id, err := repo.Create(c)
	require.NoError(t, err)
	require.NotZero(t, id)

	// reload and check
	got := &entities.Classroom{}
	require.NoError(t, db.Preload("Teacher").First(got, id).Error)
	require.Equal(t, teacher1.ID, got.TeacherID)

	// Listing by teacher name
	list, total, err := repo.Listing("Mr T", 1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int64(1))
	require.GreaterOrEqual(t, len(list), 1)
	require.Equal(t, "Mr T", list[0].Teacher.Name)

	// Listing by year (numeric search)
	list, total, err = repo.Listing("1", 1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, int64(1))

	// Update classroom
	up, err := repo.Update(id, &entities.Classroom{TeacherID: teacher2.ID, Year: 3})
	require.NoError(t, err)
	require.Equal(t, teacher2.ID, up.TeacherID)
	require.Equal(t, 3, up.Year)

	// Delete classroom
	require.NoError(t, repo.Delete(id))
	err = db.First(&entities.Classroom{}, id).Error
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
