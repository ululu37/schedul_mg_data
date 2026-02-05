package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSubjectRepo_CRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Subject{}))

	repo := &SubjectRepo{DB: db}

	var id uint
	t.Run("Create", func(t *testing.T) {
		s := &entities.Subject{Name: "Math"}
		var err error
		id, err = repo.Create(s)
		require.NoError(t, err)
		require.NotZero(t, id)
	})

	t.Run("Update", func(t *testing.T) {
		up, err := repo.Update(id, &entities.Subject{Name: "MathV2"})
		require.NoError(t, err)
		require.Equal(t, "MathV2", up.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		require.NoError(t, repo.Delete(id))
		// delete again -> not found
		require.ErrorIs(t, repo.Delete(id), gorm.ErrRecordNotFound)
	})
}
