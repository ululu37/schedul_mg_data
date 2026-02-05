package repo

import (
	"testing"

	"scadulDataMono/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCuriculumRepo_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&entities.Subject{}, &entities.Term{}, &entities.PreCurriculum{}, &entities.SubjectInPreCurriculum{}, &entities.Curriculum{}, &entities.SubjectInCurriculum{}))

	repo := &CurriculumRepo{DB: db}

	// create subjects via SubjectRepo and a pre-curiculum
	subjectRepo := &SubjectRepo{DB: db}
	s1 := &entities.Subject{Name: "Math"}
	s2 := &entities.Subject{Name: "Physics"}

	_, err = subjectRepo.Create(s1)
	require.NoError(t, err)
	_, err = subjectRepo.Create(s2)
	require.NoError(t, err)

	// create a PreCuriculum with those subjects
	preRepo := &PreCuriculumRepo{DB: db}
	pre := &entities.PreCurriculum{
		Name: "Pre A",
		SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{
			{SubjectID: s1.ID, Credit: 3},
			{SubjectID: s2.ID, Credit: 4},
		},
	}
	preID, err := preRepo.Create(pre)
	preInst, _ := preRepo.GetByID(preID)
	require.NoError(t, err)

	// create terms
	term1 := &entities.Term{Name: "T1"}
	term2 := &entities.Term{Name: "T2"}
	require.NoError(t, db.Create(term1).Error)
	require.NoError(t, db.Create(term2).Error)

	var curID uint
	var sic1ID uint

	t.Run("Create with subjects", func(t *testing.T) {
		cur := &entities.Curriculum{
			Name: "Cur A",
			SubjectInCurriculum: []entities.SubjectInCurriculum{
				{SubjectInPreCurriculumID: preInst.SubjectInPreCurriculum[0].ID, TermID: term1.ID},
				{SubjectInPreCurriculumID: preInst.SubjectInPreCurriculum[0].ID, TermID: term2.ID},
			},
		}
		id, err := repo.Create(cur)
		require.NoError(t, err)
		require.NotZero(t, id)
		curID = id
		// reload from DB to capture assigned child IDs (GORM may fill nested IDs only after load)
		created, err := repo.GetByID(id, "")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(created.SubjectInCurriculum), 1)
		sic1ID = created.SubjectInCurriculum[0].ID
	})

	t.Run("GetByID - all terms", func(t *testing.T) {
		c, err := repo.GetByID(curID, "")
		require.NoError(t, err)
		require.Equal(t, "Cur A", c.Name)
		require.Len(t, c.SubjectInCurriculum, 2)
	})

	t.Run("GetByID - filter term", func(t *testing.T) {
		c, err := repo.GetByID(curID, "T1")
		require.NoError(t, err)
		require.Len(t, c.SubjectInCurriculum, 1)
		require.Equal(t, "T1", c.SubjectInCurriculum[0].Term.Name)
	})

	t.Run("Listing and count", func(t *testing.T) {
		// add another curiculum
		_, err := repo.Create(&entities.Curriculum{Name: "Other"})
		require.NoError(t, err)

		list, total, err := repo.Listing("Cur", 10, 1)
		require.NoError(t, err)
		require.GreaterOrEqual(t, total, int64(1))
		require.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("Update name", func(t *testing.T) {
		up, err := repo.Update(curID, &entities.Curriculum{Name: "Cur A Updated"})
		require.NoError(t, err)
		require.Equal(t, "Cur A Updated", up.Name)
	})

	t.Run("AddSubject", func(t *testing.T) {
		// create another subject via SubjectRepo and a pre-curiculum for it
		s3 := &entities.Subject{Name: "Chem"}
		subjectRepo := &SubjectRepo{DB: db}
		id3, err := subjectRepo.Create(s3)
		require.NoError(t, err)
		pre := &entities.PreCurriculum{Name: "Pre Chem", SubjectInPreCurriculum: []entities.SubjectInPreCurriculum{{SubjectID: id3, Credit: 2}}}
		preRepo := &PreCuriculumRepo{DB: db}
		preID2, err := preRepo.Create(pre)
		require.NoError(t, err)
		// create term T3
		term3 := &entities.Term{Name: "T3"}
		require.NoError(t, db.Create(term3).Error)
		// add to curiculum using PreCuriculumID and TermID
		err = repo.AddSubject([]entities.SubjectInCurriculum{{CurriculumID: curID, SubjectInPreCurriculumID: preID2, TermID: term3.ID}})
		require.NoError(t, err)
		c, err := repo.GetByID(curID, "")
		require.NoError(t, err)
		require.Len(t, c.SubjectInCurriculum, 3)
	})

	t.Run("UpdateTerm", func(t *testing.T) {
		// create term TX
		tx := &entities.Term{Name: "TX"}
		require.NoError(t, db.Create(tx).Error)
		err := repo.UpdateTerm([]SubjectTermUpdate{{ID: sic1ID, TermID: tx.ID}})
		require.NoError(t, err)
		c, err := repo.GetByID(curID, "TX")
		require.NoError(t, err)
		require.Len(t, c.SubjectInCurriculum, 1)
		require.Equal(t, "TX", c.SubjectInCurriculum[0].Term.Name)
	})

	t.Run("RemoveSubject", func(t *testing.T) {
		require.NoError(t, repo.RemoveSubject([]uint{sic1ID}))
		c, err := repo.GetByID(curID, "")
		require.NoError(t, err)
		// len decreases by 1
		require.Len(t, c.SubjectInCurriculum, 2)
	})

	t.Run("Delete curiculum", func(t *testing.T) {
		require.NoError(t, repo.Delete(curID))
		_, err := repo.GetByID(curID, "")
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
