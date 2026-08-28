package api

import (
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/config"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/models/permission"
	_ "github.com/abozorov/school_online/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Option struct {
	Conf       *config.Config
	Middleware *middleware.Middleware
	Handler    *handlers.Handler
}

func NewRouter(opt *Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), opt.Middleware.Logging(), opt.Middleware.CORS())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve Frontend static assets directly from Go
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.StaticFile("/", "./frontend/index.html")

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if (len(path) >= 4 && path[:4] == "/api") || (len(path) >= 8 && path[:8] == "/swagger") {
			c.JSON(404, gin.H{"error": "Endpoint not found"})
			return
		}
		c.File("./frontend/index.html")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// auth
	authApi := router.Group("/api/auth")
	authApi.POST(
		"login",
		opt.Handler.Login,
	)

	// user
	userApi := router.Group("/api/user")
	userApi.Use(opt.Middleware.Auth())
	userApi.GET(
		"/:id",
		opt.Middleware.RBAC(permission.UserView),
		opt.Handler.GetById,
	)

	userApi.GET(
		"/list",
		opt.Middleware.RBAC(permission.UserList),
		opt.Handler.List,
	)

	userApi.POST(
		"",
		opt.Middleware.RBAC(permission.UserCreate),
		opt.Handler.Create,
	)

	userApi.PATCH(
		"",
		opt.Middleware.RBAC(permission.UserUpdate),
		opt.Handler.UpdateById,
	)

	userApi.DELETE(
		"/:id",
		opt.Middleware.RBAC(permission.UserDelete),
		opt.Handler.DeleteById,
	)

	userApi.GET(
		"/subject/:id",
		opt.Middleware.RBAC(permission.SubjectCreate),
		opt.Handler.GetSubjectById,
	)

	userApi.GET(
		"/subject/list",
		opt.Middleware.RBAC(permission.SubjectCreate),
		opt.Handler.ListSubjects,
	)

	userApi.POST(
		"/subject",
		opt.Middleware.RBAC(permission.SubjectCreate),
		opt.Handler.CreateSubject,
	)

	userApi.PATCH(
		"/subject",
		opt.Middleware.RBAC(permission.SubjectCreate),
		opt.Handler.UpdateSubject,
	)

	// classroom
	classroomApi := router.Group("/api/classroom")
	classroomApi.Use(opt.Middleware.Auth())
	classroomApi.GET(
		"/:id",
		opt.Middleware.RBAC(permission.ClassroomView),
		opt.Handler.GetClassroomById,
	)

	classroomApi.GET(
		"/list",
		opt.Middleware.RBAC(permission.ClassroomList),
		opt.Handler.ListClassrooms,
	)

	classroomApi.POST(
		"",
		opt.Middleware.RBAC(permission.ClassroomCreate),
		opt.Handler.CreateClassroom,
	)

	classroomApi.PATCH(
		"",
		opt.Middleware.RBAC(permission.ClassroomUpdate),
		opt.Handler.UpdateClassroomById,
	)

	classroomApi.DELETE(
		"/:id",
		opt.Middleware.RBAC(permission.ClassroomDelete),
		opt.Handler.DeleteClassroomById,
	)

	// schedule
	scheduleApi := router.Group("/api/schedule")
	scheduleApi.Use(opt.Middleware.Auth())
	scheduleApi.GET(
		"/classroom/:id",
		opt.Middleware.RBAC(permission.ScheduleClassroomView),
		opt.Handler.GetScheduleByClassroomId,
	)

	scheduleApi.GET(
		"/teacher/:id",
		opt.Middleware.RBAC(permission.ScheduleTeacherView),
		opt.Handler.GetScheduleByTeacherId,
	)

	scheduleApi.POST(
		"",
		opt.Middleware.RBAC(permission.ScheduleCreate),
		opt.Handler.CreateSchedule,
	)

	scheduleApi.PATCH(
		"",
		opt.Middleware.RBAC(permission.ScheduleUpdate),
		opt.Handler.UpdateScheduleById,
	)

	scheduleApi.DELETE(
		"/:id",
		opt.Middleware.RBAC(permission.ScheduleDelete),
		opt.Handler.DeleteScheduleById,
	)

	// journal
	journalApi := router.Group("/api/journal")
	journalApi.Use(opt.Middleware.Auth())
	journalApi.POST(
		"/student/:id",
		opt.Middleware.RBAC(permission.JournalStudentViewMe),
		opt.Handler.GetJournalByStudentId,
	)

	journalApi.POST(
		"/classroom/:id",
		opt.Middleware.RBAC(permission.JournalClassroomView),
		opt.Handler.GetJournalByClassroomId,
	)

	journalApi.PATCH(
		"",
		opt.Middleware.RBAC(permission.JournalGradeUpdate),
		opt.Handler.UpdateJournal,
	)
	return router
}
