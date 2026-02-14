package echohttp

import (
	"net/http"
	"scadulDataMono/config"

	router "scadulDataMono/infra/echo_http/router"
	"scadulDataMono/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type echoServer struct {
	config               *config.Config
	subjectUsecase       *usecase.SubjectMg
	preCurriculumUsecase *usecase.PreCurriculum
	teacherMgUsecase     *usecase.TeacherMg
	teacherEverlute      *usecase.TeacherEverlute
	termUsecase          *usecase.Term
	studentMg            *usecase.StudentMg
	classroomMg          *usecase.ClassroomMg
	curriculumMg         *usecase.CurriculumMg
	scadulStudentMg      *usecase.ScadulStudentMg
	scadulTeacherMg      *usecase.ScadulTeacherMg
	auth                 *usecase.Auth
	scheduleUsecase      *usecase.ScheduleUsecase
	importPrecuriculum   *usecase.ImportPrecuriculum
	app                  *echo.Echo
}

func NewEchoServer(config *config.Config,
	subjectUsecase *usecase.SubjectMg,
	preCurriculumUsecase *usecase.PreCurriculum,
	teacherMgUsecase *usecase.TeacherMg,
	teacherEverlute *usecase.TeacherEverlute,
	termUsecase *usecase.Term,
	studentMg *usecase.StudentMg,
	classroomMg *usecase.ClassroomMg,
	curriculumMg *usecase.CurriculumMg,
	scadulStudentMg *usecase.ScadulStudentMg,
	scadulTeacherMg *usecase.ScadulTeacherMg,
	auth *usecase.Auth,
	scheduleUsecase *usecase.ScheduleUsecase,
	importPrecuriculum *usecase.ImportPrecuriculum,
) *echoServer {
	return &echoServer{
		config:               config,
		subjectUsecase:       subjectUsecase,
		preCurriculumUsecase: preCurriculumUsecase,
		teacherMgUsecase:     teacherMgUsecase,
		teacherEverlute:      teacherEverlute,
		termUsecase:          termUsecase,
		studentMg:            studentMg,
		classroomMg:          classroomMg,
		curriculumMg:         curriculumMg,
		scadulStudentMg:      scadulStudentMg,
		scadulTeacherMg:      scadulTeacherMg,
		auth:                 auth,
		scheduleUsecase:      scheduleUsecase,
		importPrecuriculum:   importPrecuriculum,
	}
}

func (s *echoServer) Start() {
	s.app = echo.New()

	// CORS middleware
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{s.config.Server.CorsOrigin},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Register PreCurriculum routes
	router.RegisterPreCurriculumRoutes(s.app, s.preCurriculumUsecase, s.importPrecuriculum)
	router.RegisterTeacherRoutes(s.app, s.teacherMgUsecase, s.teacherEverlute)
	router.RegisterTermRoutes(s.app, s.termUsecase)
	router.RegisterStudentRoutes(s.app, s.studentMg)
	//router.RegisterSubjectRoutes(s.app, s.subjectUsecase)
	router.RegisterClassroomRoutes(s.app, s.classroomMg)
	router.RegisterCurriculumRoutes(s.app, s.curriculumMg)
	router.RegisterScadulStudentRoutes(s.app, s.scadulStudentMg)
	router.RegisterScadulTeacherRoutes(s.app, s.scadulTeacherMg)
	router.RegisterAuthRoutes(s.app, s.auth)
	router.RegisterScheduleRoutes(s.app, s.scheduleUsecase)

	// TODO: add other routes and middleware
	s.app.Logger.Fatal(s.app.Start(s.config.Server.Port))
}
