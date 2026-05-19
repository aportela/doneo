export default {
  "Sign in": "Iniciar sesión",

  // common actions/operations
  Actions: "Acciones",
  Add: "Añadir",
  Update: "Modificar",
  Save: "Guardar",
  Delete: "Eliminar",
  Restore: "Restaurar",
  Cancel: "Cancelar",
  Refresh: "Recargar",
  "Do you want to continue ?": "¿ Desea continuar ?",

  // errors
  "Uncaught exception": "Excepción no controlada",
  "Invalid API response code": "Código de respuesta de la API incorrecto",

  // text filter
  searchByNameDefaultPlaceholder: "buscar por nombre",
  searchByEmailDefaultPlaceholder: "buscar por email",

  // date filter
  "select date": "seleccione una fecha",
  "Any date": "Cualquier fecha",
  "Custom date": "Fecha específica",
  Yesterday: "Ayer",
  Today: "Hoy",
  Tomorrow: "Mañana",
  "Last week": "La semana pasada",
  "This week": "Esta semana",
  "Next week": "La próxima semana",
  "Last month": "El mes pasado",
  "This month": "Este mes",
  "Next month": "El próximo mes",
  "Last year": "El año pasado",
  "This year": "Este año",
  "Next year": "El próximo año",

  // user admin permission filter
  userAdminPermissionSelectorValueAll: "Todos",
  userAdminPermissionSelectorValueOnlyAdministrators: "Sólo administradores",
  userAdminPermissionSelectorValueOnlyUsers: "Sólo usuarios",

  // users
  "Manage users": "Gestión de usuarios",
  TotalUsersPagerLabel: "Total usuarios: {total}",
  UserAdminPermissionTableHeader: "Permisos",
  UsernameTableHeader: "Nombre",
  EmailTableHeader: "Email",
  CreatedAtTableHeader: "Creado en",
  UpdatedAtTableHeader: "Actualizado en",
  DeletedAtTableHeader: "Eliminado en",
  AdminTypeValue: "Administrador",
  UserTypeValue: "Usuario",
  "Add user": "Añadir usuario",
  "Update user": "Modificar usuario",
  "Delete user": "Eliminar usuario",
  "Restore user": "Restaurar usuario",
  "No users found": "No se han encontrado usuarios",
  deleteUserConfirmation:
    'Está a punto de eliminar el usuario "{name}" del sistema.',
  restoreUserConfirmation:
    'Está a punto de restaurar el usuario "{name}" del sistema.',
  "There was a problem while refreshing the user list":
    "Se ha producido un error durante la búsqueda de usuarios",
  "We couldn’t find the specified user": "El usuario especificado no existe",
  "There was a problem while deleting the user":
    "Se ha producido un error al eliminar el usuario",
  "There was a problem while restoring the user":
    "Se ha producido un error al restaurar el usuario",
  "There was a problem while adding the user data":
    "Se ha producido un error al añadir el usuario",
  "There was a problem while updating the user data":
    "Se ha producido un error al modificar el usuario",
  "There was a problem while loading the user data":
    "Se ha producido un error al cargar el usuario",

  userAddedNotification: 'El usuario "{name}" se ha añadido',
  userUpdatedNotification: 'El usuario "{name}" se ha modificado',
  userDeletedNotification: 'El usuario "{name}" se ha eliminado',
  userRestoredNotification: 'El usuario "{name}" se ha restaurado',
  userFormNameLabel: "Nombre",
  userFormEmailLabel: "Email",
  userFormPasswordLabel: "Contraseña",
  userFormPermissionsLabel: "Permisos",
  userFormNameFieldPlaceholder: "Introduzca el nombre",
  userFormEmailFieldPlaceholder: "Introduzca la dirección de correo",
  userFormPasswordFieldPlaceholder: "Introduzca la contraseña",
  userFormNameFieldEmptyError: "El nombre está vacío",
  userFormNameFieldTooLargeError: "El nombre es muy largo",
  userFormNameFieldAlreadyExists: "El nombre ya existe",
  userFormEmailFieldEmptyError: "El email está vacío",
  userFormEmailFieldInvalidError: "El email no tiene un formato válido",
  userFormEmailFieldTooLargeError: "El email es muy largo",
  userFormEmailFieldAlreadyExists: "El email ya existe",
  userFormPasswordFieldEmptyError: "La contraseña está vacía",
  userFormPasswordFieldTooShortError: "La contraseña es muy corta",
  userFormChangePasswordButtonLabel: "Modificar contraseña",
  superUserPermission: "Super usuario",
  normalUserPermission: "Usuario normal",

  // role permission typefilter
  rolePermissionTypeSelectorValueAll: "Todos",
  rolePermissionTypeSelectorValueAllowed: "Sólo permitidos",
  rolePermissionTypeSelectorValueNotAllowed: "Sólo no permitidos",

  // user roles
  "Manage roles": "Gestión de roles",
  TotalRolesPagerLabel: "Total roles: {total}",
  RoleNameTableHeader: "Nombre",
  RoleProjectPermissionsTableHeader: "Permisos de proyecto",
  RoleTasksPermissionsTableHeader: "Permisos de tarea",
  "Update project allowed": "Permitir modificar proyecto",
  "Delete project allowed": "Permitir eliminar proyecto",
  "View project allowed": "Permitir ver proyecto",
  "Add task allowed": "Permitir añadir tarea",
  "Update task allowed": "Permitir modificar tarea",
  "Delete task allowed": "Permitir eliminar tarea",
  "View task allowed": "Permitir ver tarea",
  "Add role": "Añadir rol",
  "Update role": "Modificar rol",
  "Delete role": "Eliminar rol",
  "No roles found": "No se han encontrado roles",
  "Role permissions": "Permisos del rol",
  "Project permissions": "Permisos del proyecto",
  "Task permissions": "Permisos de la tarea",
  deleteRoleConfirmation:
    'Está a punto de eliminar el rol "{name}" del sistema.',
  "There was a problem while refreshing the role list":
    "Se ha producido un error durante la búsqueda de roles",
  "We couldn’t find the specified role": "El rol especificado no existe",
  "There was a problem while deleting the role":
    "Se ha producido un error al eliminar el rol",
  "There was a problem while adding the role data":
    "Se ha producido un error al añadir el rol",
  "There was a problem while updating the role data":
    "Se ha producido un error al modificar el rol",
  "There was a problem while loading the role data":
    "Se ha producido un error al cargar el rol",
  roleAddedNotification: 'El rol "{name}" se ha añadido',
  roleUpdatedNotification: 'El rol "{name}" se ha modificado',
  roleDeletedNotification: 'El rol "{name}" se ha eliminado',
  roleFormNameLabel: "Nombre",
  roleFormNameFieldPlaceholder: "Introduzca el nombre",
  roleFormNameFieldEmptyError: "El nombre está vacío",
  roleFormNameFieldTooLargeError: "El nombre es muy largo",
  roleFormNameFieldAlreadyExists: "El nombre ya existe",

  // project types
  "Manage project types": "Gestión de tipos de proyectos",
  TotalProjectTypesPagerLabel: "Total tipos de proyecto: {total}",
  ProjectTypeNameTableHeader: "Nombre",
  "Add project type": "Añadir tipo de proyecto",
  "Update project type": "Actualizar tipo de proyecto",
  "Delete project type": "Eliminar tipo de proyecto",
  "No project types found": "No se han encontrado tipos de proyectos",
  deleteProjectTypeConfirmation:
    'Está a punto de eliminar el tipo de proyecto "{name}" del sistema.',
  "There was a problem loading the project type data":
    "Se produjo un problema al cargar los datos del tipo de proyecto",
  "There was a problem adding the project type data":
    "Se produjo un problema al añadir los datos del tipo de proyecto",
  "There was a problem updating the project type data":
    "Se produjo un problema al modificar los datos del tipo de proyecto",
  "There was a problem deleting the project type data":
    "Se produjo un problema al eliminar los datos del tipo de proyecto",
  "There was a problem while adding the project type data":
    "Se produjo un problema al añadir los datos del tipo de proyecto",
  "There was a problem while updating the project type data":
    "Se produjo un problema al modificar los datos del tipo de proyecto",
  "There was a problem while loading the project type data":
    "Se produjo un problema al cargar los datos del tipo de proyecto",
  "We couldn’t find the specified project type":
    "El tipo de proyecto especificado no existe",
  projectTypeFormNameFieldEmptyError: "El nombre está vacío",
  projectTypeFormNameFieldTooLargeError: "El nombre es muy largo",
  projectTypeFormNameFieldAlreadyExists: "El nombre ya existe",
  projectTypeFormNameFieldPlaceholder:
    "Introduzca el nombre del tipo de proyecto",
  projectTypeAddedNotification: 'El tipo de proyecto "{name}" se ha añadido',
  projectTypeUpdatedNotification:
    'El tipo de proyecto "{name}" se ha modificado',
  projectTypeDeletedNotification:
    'El tipo de proyecto "{name}" se ha eliminado',

  // project statuses
  "Manage project statuses": "Gestión de estados de proyectos",
  TotalProjectStatusesPagerLabel: "Total estados de proyecto: {total}",
  ProjectStatusNameTableHeader: "Nombre",
  ProjectStatusIndexTableHeader: "Índice",
  "No project statuses found": "No se han encontrado estados de proyectos",
  "Add project status": "Añadir estado de proyecto",
  "Update project status": "Actualizar estado de proyecto",
  "Delete project status": "Eliminar estado de proyecto",
  deleteProjectStatusConfirmation:
    'Está a punto de eliminar el estado de proyecto "{name}" del sistema.',
  "There was a problem loading the project status data":
    "Se produjo un problema al cargar los datos del estado de proyecto",
  "There was a problem adding the project status data":
    "Se produjo un problema al añadir los datos del estado de proyecto",
  "There was a problem updating the project status data":
    "Se produjo un problema al modificar los datos del estado de proyecto",
  "There was a problem deleting the project status data":
    "Se produjo un problema al eliminar los datos del estado de proyecto",
  "There was a problem while loading the project status data":
    "Se produjo un problema al cargar los datos del estado de proyecto",
  "We couldn’t find the specified project status":
    "El estado de proyecto especificado no existe",
  projectStatusFormNameFieldEmptyError: "El nombre está vacío",
  projectStatusFormNameFieldTooLargeError: "El nombre es muy largo",
  projectStatusFormNameFieldAlreadyExists: "El nombre ya existe",
  projectStatusFormNameFieldPlaceholder:
    "Introduzca el nombre del estado de proyecto",
  projectStatusAddedNotification:
    'El estado de proyecto "{name}" se ha añadido',
  projectStatusUpdatedNotification:
    'El estado de proyecto "{name}" se ha modificado',
  projectStatusDeletedNotification:
    'El estado de proyecto "{name}" se ha eliminado',
  projectStatusIndexMovedNotification:
    'El índice del estado de proyecto "{name}" se ha movido',

  // auth
  "Login to your account": "Inicie sesión con su cuenta",
  "Session lost... re-auth required":
    "La sesión se ha perdido... se requiere re-autentificación",
  loginFormEmailLabel: "Email",
  loginFormPasswordLabel: "Contraseña",
  loginFormEmailFieldPlaceholder: "Introduzca su correo electrónico",
  loginFormPasswordFieldPlaceholder: "Introduzca su contraseña",
  loginFormEmailFieldEmptyError: "El email está vacío",
  loginFormEmailFieldInvalidError: "El email no tiene un formato válido",
  loginFormEmailFieldTooLargeError: "El email es muy largo",
  loginFormPasswordFieldEmptyError: "La contraseña está vacía",
  loginFormPasswordFieldTooShortError: "La contraseña es muy corta",
  loginFormEmailFieldAPINotFound: "El email no existe",
  loginFormPasswordFieldIncorrect: "La contraseña es incorrecta",
  "Sign in error": "Error al iniciar la sesión",

  // other
  Error: "Error",
  Name: "Nombre",
  Preview: "Vista previa",
  Color: "Color",

  // tooltips
  "show password": "mostrar contraseña",
  "hide password": "ocultar contraseña",

  // notifications
  "Enable notifications": "Activar notificaciones",
  "Disable notifications": "Desactivar notificaciones",

  // hints
  "Open project page": "Abrir página del proyecto",
  "Switch to dark mode": "Cambiar a modo oscuro",
  "Switch to light mode": "Cambiar a modo claro",
  "Switch to top navigation": "Cambiar a navegación superior",
  "Switch to side navigation": "Cambiar a navegación lateral",

  // pager
  labelPageOfTotalPages: "Página {currentPage} de {totalPages}",
  "Total results:": "Total resultados:",
  "All results": "Todos los resultados",
  labelSelectorResultsPage: "{number} resultados/página",
};
