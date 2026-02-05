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
	aiAgent := aiAgent.NewAiAgent(
		"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",
		"https://openrouter.ai/api/v1/chat/completions",
	)

	// Init usecases
	subjectMg := &usecase.SubjectMg{SubjectRepo: subjectRepo}
	preCurriculum := &usecase.PreCurriculum{PreRepo: preCurriculumRepo, SubjectMg: subjectMg}
	teacherMg := usecase.NewTeacherMg(teacherRepo, authRepo)
	teacherEverlute := usecase.NewTeacherEverlute(teacherMg, subjectMg, aiAgent)

	// Start echo server
	server := echohttp.NewEchoServer(&cfg, subjectMg, preCurriculum, teacherMg, teacherEverlute)
	server.Start()
}
