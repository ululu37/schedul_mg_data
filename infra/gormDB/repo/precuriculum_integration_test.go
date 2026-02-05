package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPreCuriculumRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Subject{}, &entities.PreCurriculum{}, &entities.SubjectInPreCurriculum{}))

	repo := &PreCuriculumRepo{DB: db}
	subRepo := &SubjectRepo{DB: db}

	// create subjects and capture IDs
	s1 := &entities.Subject{Name: "Math"}
	id1, err := subRepo.Create(s1)
	require.NoError(t, err)
	s2 := &entities.Subject{Name: "Physics"}
	id2, err := subRepo.Create(s2)
	require.NoError(t, err)

	var preID uint
	var sipID uint

	t.Run("Create", func(t *testing.T) {
		pre := &entities.PreCurriculum{
			Name: "Pre A",
			SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
				{SubjectID: id1, Credit: 3},
				{SubjectID: id2, Credit: 4},
			},
		}
		id, err := repo.Create(pre)
		require.NoError(t, err)
		require.NotZero(t, id)
		preID = id
		sipID = pre.SubjectInPreCurriculum[0].ID
	})

	t.Run("GetByID", func(t *testing.T) {
		p, err := repo.GetByID(preID)
		require.NoError(t, err)
		require.Equal(t, "Pre A", p.Name)
		require.Len(t, p.SubjectInPreCurriculum, 2)
	})

	t.Run("Listing", func(t *testing.T) {
		list, total, err := repo.Listing("Pre", 10, 1)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("AddSubject", func(t *testing.T) {
		// create another subject via SubjectRepo and use its returned ID
		s3 := &entities.Subject{Name: "Chem"}
		id3, err := subRepo.Create(s3)
		require.NoError(t, err)
		require.NoError(t, repo.AddSubject([]entities.SubjectInPreCurriculum{{PreCurriculumID: preID, SubjectID: id3, Credit: 2}}))
		p, err := repo.GetByID(preID)
		require.NoError(t, err)
		require.Len(t, p.SubjectInPreCurriculum, 3)
	})

	t.Run("RemoveSubject", func(t *testing.T) {
		require.NoError(t, repo.RemoveSubject([]uint{sipID}))
		p, err := repo.GetByID(preID)
		require.NoError(t, err)
		require.Len(t, p.SubjectInPreCurriculum, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		require.NoError(t, repo.Delete(preID))
		_, err := repo.GetByID(preID)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
