package app

import (
	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/config"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers/attachmenthandler"
	"github.com/aportela/doneo/internal/handlers/historyoperationhandler"
	"github.com/aportela/doneo/internal/handlers/identityhandler"
	"github.com/aportela/doneo/internal/handlers/notehandler"
	"github.com/aportela/doneo/internal/handlers/profilehandler"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
	"github.com/aportela/doneo/internal/handlers/projectpermissionhandler"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/taskhandler"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/tasktimetrackinghandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/handlers/usertimerhandler"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/repositories/tagrepository"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/taskrepository"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/repositories/tasktimetrackingrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/repositories/usertimerrepository"
	"github.com/aportela/doneo/internal/services/attachmentservice"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/services/identityservice"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/aportela/doneo/internal/services/profileservice"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/services/taskservice"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/services/tasktimetrackingservice"
	"github.com/aportela/doneo/internal/services/userservice"
	"github.com/aportela/doneo/internal/services/usertimerservice"
)

type App struct {
	DB    database.Database
	Cfg   config.Configuration
	Cache cache.PermissionCache

	AttachmentHandler        attachmenthandler.AttachmentHandler
	IdentityHandler          identityhandler.IdentityHandler
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
	TaskTimeTrackingHandler  tasktimetrackinghandler.TaskTimeTrackingHandler
	UserTimerHandler         usertimerhandler.UserTimerHandler
	UserHandler              userhandler.UserHandler
	ProfileHandler           profilehandler.ProfileHandler
}

func NewApp(
	db database.Database,
	cfg config.Configuration,
	cache cache.PermissionCache,

) *App {

	attachmentRepository := attachmentrepository.NewRepository()
	historyOperationRepository := historyoperationrepository.NewRepository()
	noteRepository := noterepository.NewRepository()
	projectPermissionRepository := projectpermissionrepository.NewRepository()
	projectPriorityRepository := projectpriorityrepository.NewRepository()
	projectRepository := projectrepository.NewRepository()
	projectStatusRepository := projectstatusrepository.NewRepository()
	taskRepository := taskrepository.NewRepository()
	tagRepository := tagrepository.NewRepository()
	projectTypeRepository := projecttyperepository.NewRepository()
	roleRepository := rolerepository.NewRepository()
	taskPriorityRepository := taskpriorityrepository.NewRepository()
	//taskRelationRepository := taskrelationrepository.NewRepository()
	taskStatusRepository := taskstatusrepository.NewRepository()
	taskTimeTrackingRepository := tasktimetrackingrepository.NewRepository()
	userTimerRepository := usertimerrepository.NewRepository()
	userRepository := userrepository.NewRepository()

	authorizationService := authorizationservice.NewService(db, cache, userRepository, projectPermissionRepository)
	historyOperationService := historyoperationservice.NewService(authorizationService, historyOperationRepository)
	attachmentService := attachmentservice.NewService(db, cfg.Storage.AttachmentsPath, authorizationService, historyOperationService, attachmentRepository)
	identityService := identityservice.NewService(db, userRepository)
	noteService := noteservice.NewService(db, authorizationService, historyOperationService, noteRepository)
	projectPermissionService := projectpermissionservice.NewService(db, cache, authorizationService, historyOperationService, projectPermissionRepository)
	projectPriorityService := projectpriorityservice.NewService(db, authorizationService, projectPriorityRepository)
	projectService := projectservice.NewService(db, cache, authorizationService, historyOperationService, projectRepository)
	projectStatusService := projectstatusservice.NewService(db, authorizationService, projectStatusRepository)
	taskService := taskservice.NewService(db, authorizationService, historyOperationService, taskRepository, tagRepository)
	projectTypeService := projecttypeservice.NewService(db, authorizationService, projectTypeRepository)
	roleService := roleservice.NewService(db, authorizationService, roleRepository)
	taskPriorityService := taskpriorityservice.NewService(db, authorizationService, taskPriorityRepository)
	taskStatusService := taskstatusservice.NewService(db, authorizationService, taskStatusRepository)
	taskTimeTrackingService := tasktimetrackingservice.NewService(db, authorizationService, historyOperationService, taskTimeTrackingRepository)
	userTimerService := usertimerservice.NewService(db, userTimerRepository)
	userService := userservice.NewService(db, cache, authorizationService, userRepository)
	profileService := profileservice.NewService(db, cache, authorizationService, userRepository)

	attachmentHandler := attachmenthandler.NewHandler(attachmentService, cfg.Storage.MaxUploadFilesize)
	identityHandler := identityhandler.NewHandler(identityService, cfg.Auth.SecretKey, cfg.Auth.AccessTokenExpirationHours, cfg.Auth.RefreshTokenExpirationDays)
	historyOperationHandler := historyoperationhandler.NewHandler(db, historyOperationService)
	noteHandler := notehandler.NewHandler(noteService)
	projectPermissionHandler := projectpermissionhandler.NewHandler(projectPermissionService)
	projectPriorityHandler := projectpriorityhandler.NewHandler(projectPriorityService)
	projectHandler := projecthandler.NewHandler(projectService)
	projectStatusHandler := projectstatushandler.NewHandler(projectStatusService)
	taskHandler := taskhandler.NewHandler(taskService)
	projectTypeHandler := projecttypehandler.NewHandler(projectTypeService)
	roleHandler := rolehandler.NewHandler(roleService)
	taskPriorityHandler := taskpriorityhandler.NewHandler(taskPriorityService)
	taskStatusHandler := taskstatushandler.NewHandler(taskStatusService)
	taskTimeTrackingHandler := tasktimetrackinghandler.NewHandler(taskTimeTrackingService)
	userTimerHandler := usertimerhandler.NewHandler(userTimerService)
	userHandler := userhandler.NewHandler(userService)
	profileHandler := profilehandler.NewHandler(profileService)

	return &App{
		DB:    db,
		Cfg:   cfg,
		Cache: cache,

		AttachmentHandler:        attachmentHandler,
		IdentityHandler:          identityHandler,
		HistoryOperationHandler:  historyOperationHandler,
		NoteHandler:              noteHandler,
		ProjectPermissionHandler: projectPermissionHandler,
		ProjectPriorityHandler:   projectPriorityHandler,
		ProjectHandler:           projectHandler,
		ProjectStatusHandler:     projectStatusHandler,
		TaskHandler:              taskHandler,
		ProjectTypeHandler:       projectTypeHandler,
		RoleHandler:              roleHandler,
		TaskPriorityHandler:      taskPriorityHandler,
		TaskStatusHandler:        taskStatusHandler,
		TaskTimeTrackingHandler:  taskTimeTrackingHandler,
		UserTimerHandler:         userTimerHandler,
		UserHandler:              userHandler,
		ProfileHandler:           profileHandler,
	}
}
