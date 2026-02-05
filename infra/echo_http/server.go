package echohttp

import (
	"scadulDataMono/config"
	router "scadulDataMono/infra/echo_http/router"
	"scadulDataMono/usecase"

	"github.com/labstack/echo/v4"
)

type echoServer struct {
	config               *config.Config
	subjectUsecase       *usecase.SubjectMg
	preCurriculumUsecase *usecase.PreCurriculum
	teacherMgUsecase     *usecase.TeacherMg
	teacherEverlute      *usecase.TeacherEverlute
	app                  *echo.Echo
}

func NewEchoServer(config *config.Config,
	subjectUsecase *usecase.SubjectMg,
	preCurriculumUsecase *usecase.PreCurriculum,
	teacherMgUsecase *usecase.TeacherMg,
	teacherEverlute *usecase.TeacherEverlute,
) *echoServer {
	return &echoServer{
		config:               config,
		subjectUsecase:       subjectUsecase,
		preCurriculumUsecase: preCurriculumUsecase,
		teacherMgUsecase:     teacherMgUsecase,
		teacherEverlute:      teacherEverlute,
	}
}

func (s *echoServer) Start() {
	s.app = echo.New()

	// Register PreCurriculum routes
	router.RegisterPreCurriculumRoutes(s.app, s.preCurriculumUsecase)
	router.RegisterTeacherRoutes(s.app, s.teacherMgUsecase, s.teacherEverlute)

	// TODO: add other routes and middleware
	s.app.Logger.Fatal(s.app.Start(s.config.Server.Port))
}
