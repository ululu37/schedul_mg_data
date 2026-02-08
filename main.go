package main

import (
	"scadulDataMono/config"
	aiAgent "scadulDataMono/infra/Agent"
	echohttp "scadulDataMono/infra/echo_http"
	"scadulDataMono/infra/gormDB"
	"scadulDataMono/infra/gormDB/repo"
	"scadulDataMono/usecase"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	db := gormDB.NewPostgresDatabase(&cfg.Database).Connect()

	// Init repos
	subjectRepo := repo.NewSubjectRepo(db)
	preCurriculumRepo := &repo.PreCuriculumRepo{DB: db}
	teacherRepo := repo.NewTeacherRepo(db)
	authRepo := repo.NewAuthRepo(db)
	termRepo := repo.NewTermRepo(db)
	studentRepo := &repo.StudentRepo{DB: db}
	classroomRepo := &repo.ClassroomRepo{DB: db}
	curriculumRepo := &repo.CurriculumRepo{DB: db}
	scadulStudentRepo := &repo.ScadulStudentRepo{DB: db}
	scadulTeacherRepo := &repo.ScadulTeacherRepo{DB: db}
	aiAgent := aiAgent.NewAiAgent(
		cfg.AiAgent.ApiKey,
		cfg.AiAgent.Url,
	)

	// Init usecases
	subjectMg := &usecase.SubjectMg{SubjectRepo: subjectRepo}
	preCurriculum := &usecase.PreCurriculum{PreRepo: preCurriculumRepo, SubjectMg: subjectMg}
	teacherMg := usecase.NewTeacherMg(teacherRepo, authRepo)
	teacherEverlute := usecase.NewTeacherEverlute(teacherMg, subjectMg, aiAgent)
	termUsecase := usecase.NewTermUsecase(termRepo)
	studentMg := usecase.NewStudentMg(studentRepo, authRepo)
	classroomMg := usecase.NewClassroomMg(classroomRepo)
	curriculumMg := usecase.NewCurriculumMg(curriculumRepo)
	scadulStudentMg := usecase.NewScadulStudentMg(scadulStudentRepo)
	scadulTeacherMg := usecase.NewScadulTeacherMg(scadulTeacherRepo)
	auth := usecase.NewAuth(authRepo, studentMg, teacherMg)

	// Start echo server
	server := echohttp.NewEchoServer(&cfg, subjectMg, preCurriculum, teacherMg, teacherEverlute, termUsecase, studentMg, classroomMg, curriculumMg, scadulStudentMg, scadulTeacherMg, auth)
	server.Start()
}
