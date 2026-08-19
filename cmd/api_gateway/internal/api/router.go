package api

import (
	_ "github.com/abozorov/school_online/docs"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/config"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/models/permission"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

type Option struct {
	Conf       *config.Config
	Middleware *middleware.Middleware
	Handler    *handlers.Handler
}

func NewRouter(opt *Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), opt.Middleware.Logging())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	userApi.POST(
		"/subject",
		opt.Middleware.RBAC(permission.SubjectCreate),
		opt.Handler.CreateSubject,
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
	journalApi.GET(
		"/student/:id",
		opt.Middleware.RBAC(permission.JournalStudentViewMe),
		opt.Handler.GetJournalByStudentId,
	)
	
	journalApi.GET(
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
