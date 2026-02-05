package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestAuthRepo_Integration_CRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Auth{}))

	repo := &AuthRepo{DB: db}

	var id uint
	var up *entities.Auth
	updated := &entities.Auth{Username: "updated", Password: "newpass", HumanType: "t", Role: 3}

	t.Run("Create", func(t *testing.T) {
		newAuth := &entities.Auth{Username: "testuser", Password: "pass", HumanType: "s", Role: 2}
		var err error
		id, err = repo.Create(newAuth)
		require.NoError(t, err)
		require.NotZero(t, id)
	})

	t.Run("GetByUsername - found", func(t *testing.T) {
		got, err := repo.GetByUsername("testuser")
		require.NoError(t, err)
		require.Equal(t, "testuser", got.Username)
		require.Equal(t, "pass", got.Password)
		require.Equal(t, "s", got.HumanType)
		require.Equal(t, 2, got.Role)
	})

	t.Run("Update", func(t *testing.T) {
		var err error
		up, err = repo.Update(id, updated)
		require.NoError(t, err)
		require.Equal(t, "updated", up.Username)
		require.Equal(t, "newpass", up.Password)
	})

	t.Run("GetByUsername - old not found", func(t *testing.T) {
		_, err := repo.GetByUsername("testuser")
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("GetByUsername - updated found", func(t *testing.T) {
		got2, err := repo.GetByUsername("updated")
		require.NoError(t, err)
		require.Equal(t, up.ID, got2.ID)
	})

	t.Run("HumanType helpers", func(t *testing.T) {
		// up was set to HumanType "t"
		require.True(t, up.IsTeacher())
		require.False(t, up.IsStudent())
		require.Equal(t, "teacher", up.HumanTypeDisplay())

		// create other auth
		other := &entities.Auth{Username: "x", Password: "x", HumanType: ""}
		require.True(t, other.IsOther())
		require.Equal(t, "other", other.HumanTypeDisplay())
	})

	t.Run("GetByID - found", func(t *testing.T) {
		got, err := repo.GetByID(id)
		require.NoError(t, err)
		require.Equal(t, id, got.ID)
	})

	t.Run("DeleteByID - success", func(t *testing.T) {
		require.NoError(t, repo.DeleteByID(id))
	})

	t.Run("GetByID - not found after delete", func(t *testing.T) {
		_, err := repo.GetByID(id)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("DeleteByID - not found", func(t *testing.T) {
		require.ErrorIs(t, repo.DeleteByID(9999), gorm.ErrRecordNotFound)
	})

	t.Run("Update - not found", func(t *testing.T) {
		_, err := repo.Update(9999, updated)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
