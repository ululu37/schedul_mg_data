package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestTeacherRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Auth{}, &entities.Subject{}, &entities.Teacher{}, &entities.TeacherMySubject{}))

	repo := &TeacherRepo{DB: db}

	var id uint

	t.Run("Create", func(t *testing.T) {
		id, err = repo.Create(&entities.Teacher{Name: "Mr T", Resume: "res", AuthID: 1})
		require.NoError(t, err)
		require.NotZero(t, id)
	})

	t.Run("Listing - by name", func(t *testing.T) {
		list, total, err := repo.Listing("Mr T", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Listing - by auth_id", func(t *testing.T) {
		list, total, err := repo.Listing("1", 1, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Update", func(t *testing.T) {
		up, err := repo.Update(id, &entities.Teacher{Name: "Mr T2", Resume: "r2", AuthID: 2})
		require.NoError(t, err)
		require.Equal(t, "Mr T2", up.Name)
		require.Equal(t, uint(2), up.AuthID)
	})

	t.Run("AddMySubject and GetByIDs", func(t *testing.T) {
		// create subjects
		s1 := &entities.Subject{Name: "Math"}
		require.NoError(t, db.Create(s1).Error)
		s2 := &entities.Subject{Name: "Physics"}
		require.NoError(t, db.Create(s2).Error)

		require.NoError(t, repo.AddMySubject(id, []entities.TeacherMySubject{{SubjectID: s1.ID, Preference: 1}, {SubjectID: s2.ID, Preference: 2}}))

		list, err := repo.GetByIDs([]uint{id})
		require.NoError(t, err)
		require.Len(t, list, 1)
		require.GreaterOrEqual(t, len(list[0].MySubject), 1)
	})

	t.Run("UpdatePreference and RemoveMySubject", func(t *testing.T) {
		list, err := repo.GetByIDs([]uint{id})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(list), 1)
		m := list[0].MySubject[0]

		// update preference
		require.NoError(t, repo.UpdatePreference([]PreferenceUpdate{{ID: m.ID, Preference: 9}}))
		// verify
		list2, err := repo.GetByIDs([]uint{id})
		require.NoError(t, err)
		found := false
		for _, mm := range list2[0].MySubject {
			if mm.ID == m.ID {
				require.Equal(t, 9, mm.Preference)
				found = true
			}
		}
		require.True(t, found)

		// remove
		require.NoError(t, repo.RemoveMySubject([]uint{m.ID}))
		// verify removed
		list3, err := repo.GetByIDs([]uint{id})
		require.NoError(t, err)
		for _, mm := range list3[0].MySubject {
			require.NotEqual(t, m.ID, mm.ID)
		}
	})

	t.Run("Delete teacher", func(t *testing.T) {
		require.NoError(t, repo.Delete(id))
		_, _, err := repo.Listing("Mr T2", 1, 10)
		require.NoError(t, err)
	})
}
