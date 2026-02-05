package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestStudentRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Auth{}, &entities.Teacher{}, &entities.Classroom{}, &entities.Subject{}, &entities.PreCurriculum{}, &entities.Curriculum{}, &entities.SubjectInCurriculum{}, &entities.Student{}))

	subRepo := &SubjectRepo{DB: db}
	subRepo.Create(&entities.Subject{Name: "Math"})

	// create teacher and classroom
	teacher := &entities.Teacher{Name: "Mr T"}
	require.NoError(t, db.Create(teacher).Error)
	classroom := &entities.Classroom{TeacherID: teacher.ID, Year: 1}
	require.NoError(t, db.Create(classroom).Error)

	// create pre-curiculum and cuuriculum
	pre := &entities.PreCurriculum{Name: "Pre A", SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{{SubjectID: 1, Credit: 3}}}
	require.NoError(t, db.Create(pre).Error)
	term := &entities.Term{Name: "T1"}
	require.NoError(t, db.Create(term).Error)
	cur := &entities.Curriculum{Name: "Cur A", SubjectInCurriculum: []entities.SubjectInCurriculum{{SubjectInPreCurriculumID: pre.SubjectInPreCurriculum[0].ID, TermID: term.ID}}}
	require.NoError(t, db.Create(cur).Error)

	repo := &StudentRepo{DB: db}

	// create students
	s1 := &entities.Student{Name: "Alice", CurriculumID: cur.ID, ClassroomID: classroom.ID}
	s2 := &entities.Student{Name: "Bob", CurriculumID: cur.ID, ClassroomID: classroom.ID}

	id1, err := repo.Create(s1)
	require.NoError(t, err)
	id2, err := repo.Create(s2)
	require.NoError(t, err)
	require.NotZero(t, id1)
	require.NotZero(t, id2)

	t.Run("GetByClassRoom", func(t *testing.T) {
		list, err := repo.GetByClassRoom(classroom.ID)
		require.NoError(t, err)
		require.Len(t, list, 2)
	})

	t.Run("Listing - by name", func(t *testing.T) {
		list, total, err := repo.Listing("Alice", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Listing - by curriculum name", func(t *testing.T) {
		list, total, err := repo.Listing("Cur A", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Listing - by precuriculum name", func(t *testing.T) {
		list, total, err := repo.Listing("Pre A", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Listing - by teacher name", func(t *testing.T) {
		list, total, err := repo.Listing("Mr T", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})
}
