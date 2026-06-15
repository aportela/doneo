package app

import (
	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/config"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers/attachmenthandler"
	"github.com/aportela/doneo/internal/handlers/authhandler"
	"github.com/aportela/doneo/internal/handlers/historyoperationhandler"
	"github.com/aportela/doneo/internal/handlers/notehandler"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
	"github.com/aportela/doneo/internal/handlers/projectpermissionhandler"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/taskhandler"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/timerhandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/taskrepository"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/repositories/timerrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/attachmentservice"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/authservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/services/taskservice"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/services/timerservice"
	"github.com/aportela/doneo/internal/services/userservice"
)

type App struct {
	DB    database.Database
	Cfg   config.Configuration
	Cache cache.PermissionCache

	AttachmentHandler        attachmenthandler.AttachmentHandler
	AuthHandler              authhandler.AuthHandler
	HistoryOperationHandler  historyoperationhandler.HistoryOperationHandler
	NoteHandler              notehandler.NoteHandler
	ProjectPermissionHandler projectpermissionhandler.ProjectPermissionHandler
	ProjectPriorityHandler   projectpriorityhandler.ProjectPriorityHandler
	ProjectHandler           projecthandler.ProjectHandler
	ProjectStatusHandler     projectstatushandler.ProjectStatusHandler
	TaskHandler              taskhandler.TaskHandler
	ProjectTypeHandler       projecttypehandler.ProjectTypeHandler
	RoleHandler              rolehandler.RoleHandler
	TaskPriorityHandler      taskpriorityhandler.TaskPriorityHandler
	TaskStatusHandler        taskstatushandler.TaskStatusHandler
	TimerHandler             timerhandler.TimerHandler
	UserHandler              userhandler.UserHandler
}

func NewApp(
	db database.Database,
	cfg config.Configuration,
	cache cache.PermissionCache,

) *App {

	attachmentRepository := attachmentrepository.NewRepository(db)
	historyOperationRepository := historyoperationrepository.NewRepository(db)
	noteRepository := noterepository.NewRepository(db)
	projectPermissionRepository := projectpermissionrepository.NewRepository(db)
	projectPriorityRepository := projectpriorityrepository.NewRepository(db)
	projectRepository := projectrepository.NewRepository(db)
	projectStatusRepository := projectstatusrepository.NewRepository(db)
	projectTaskRepository := taskrepository.NewRepository(db)
	projectTypeRepository := projecttyperepository.NewRepository(db)
	roleRepository := rolerepository.NewRepository(db)
	//tagRepository := tagrepository.NewRepository(db)
	taskPriorityRepository := taskpriorityrepository.NewRepository(db)
	//taskRelationRepository := taskrelationrepository.NewRepository(db)
	taskStatusRepository := taskstatusrepository.NewRepository(db)
	timerRepository := timerrepository.NewRepository(db)
	userRepository := userrepository.NewRepository(db)

	historyOperationService := historyoperationservice.NewService(db, historyOperationRepository)
	attachmentService := attachmentservice.NewService(db, historyOperationService, attachmentRepository)
	authorizationService := authorizationservice.NewService(db, cache)
	authService := authservice.NewService(db, userRepository)
	noteService := noteservice.NewService(db, historyOperationService, noteRepository)
	projectPermissionService := projectpermissionservice.NewService(db, cache, historyOperationService, projectPermissionRepository)
	projectPriorityService := projectpriorityservice.NewService(db, projectPriorityRepository)
	projectService := projectservice.NewService(db, authorizationService, historyOperationService, projectRepository)
	projectStatusService := projectstatusservice.NewService(db, projectStatusRepository)
	taskService := taskservice.NewService(db, projectTaskRepository)
	projectTypeService := projecttypeservice.NewService(db, projectTypeRepository)
	roleService := roleservice.NewService(db, roleRepository)
	taskPriorityService := taskpriorityservice.NewService(db, taskPriorityRepository)
	taskStatusService := taskstatusservice.NewService(db, taskStatusRepository)
	timerService := timerservice.NewService(db, timerRepository)
	userService := userservice.NewService(db, userRepository)

	attachmentHandler := attachmenthandler.NewHandler(attachmentService, cfg.Storage.AttachmentsPath)
	authHandler := authhandler.NewHandler(authService, cfg.Auth.SecretKey, cfg.Auth.AccessTokenExpirationHours, cfg.Auth.RefreshTokenExpirationDays)
	historyOperationHandler := historyoperationhandler.NewHandler(historyOperationService)
	noteHandler := notehandler.NewHandler(noteService)
	projectPermissionHandler := projectpermissionhandler.NewHandler(projectPermissionService)
	projectPriorityHandler := projectpriorityhandler.NewHandler(projectPriorityService)
	projectHandler := projecthandler.NewHandler(projectService)
	projectStatusHandler := projectstatushandler.NewHandler(projectStatusService)
	taskHandler := taskhandler.NewHandler(taskService)
	ProjectTypeHandler := projecttypehandler.NewHandler(projectTypeService)
	roleHandler := rolehandler.NewHandler(roleService)
	TaskPriorityHandler := taskpriorityhandler.NewHandler(taskPriorityService)
	TaskStatusHandler := taskstatushandler.NewHandler(taskStatusService)
	TimerHandler := timerhandler.NewHandler(timerService)
	userHandler := userhandler.NewHandler(userService)

	return &App{
		DB:    db,
		Cfg:   cfg,
		Cache: cache,

		AttachmentHandler:        *attachmentHandler,
		AuthHandler:              *authHandler,
		HistoryOperationHandler:  *historyOperationHandler,
		NoteHandler:              *noteHandler,
		ProjectPermissionHandler: *projectPermissionHandler,
		ProjectPriorityHandler:   *projectPriorityHandler,
		ProjectHandler:           *projectHandler,
		ProjectStatusHandler:     *projectStatusHandler,
		TaskHandler:              *taskHandler,
		ProjectTypeHandler:       *ProjectTypeHandler,
		RoleHandler:              *roleHandler,
		TaskPriorityHandler:      *TaskPriorityHandler,
		TaskStatusHandler:        *TaskStatusHandler,
		TimerHandler:             *TimerHandler,
		UserHandler:              *userHandler,
	}
}
