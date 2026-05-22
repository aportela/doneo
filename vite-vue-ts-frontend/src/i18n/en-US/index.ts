export default {
  layouts: {
    errorNotFound: {
      title: "404 Not Found",
      description: "We can't find what you're looking for.",
      buttons: {
        goHome: "Go home",
        goBack: "Go back",
      },
    },
  },
  shared: {
    components: {
      dialogs: {
        confirmation: {
          continueMessage: "Do you want to continue ?",
        },
      },
      table: {
        header: {
          columns: {
            actions: "Actions",
          },
        },
      },
    },
    buttons: {
      Save: {
        label: "Save",
      },
      Cancel: {
        label: "Cancel",
      },
      Delete: {
        label: "Delete",
      },
      Restore: {
        label: "Restore",
      },
      Refresh: {
        label: "Refresh",
      },
      Add: {
        label: "Add",
      },
      Update: {
        label: "Update",
      },
    },
    errorMessages: {
      Error: "Error",
      uncaugthException: "Uncaught exception",
      invalidAPIResponseCode: "Invalid API response code",
    },
    warningMessages: {
      fieldIsRequired: "Field is required",
      fieldHasInvalidFormat: "Field has invalid format",
      fieldExceedsMaxLength: "Field exceeds max length",
      fieldIsBelowMinimumLength: "Field is below minimum length",
    },
  },
  modules: {
    // auth module
    auth: {
      components: {
        LoginForm: {
          inputs: {
            email: {
              label: "Email",
              placeholder: "Enter your email address",
            },
            password: {
              label: "Password",
              placeholder: "Enter your password",
              hidePasswordTooltipIcon: "hide password",
              showPasswordTooltipIcon: "show passwword",
            },
          },
          buttons: {
            signIn: {
              label: "Sign in",
            },
          },
          warnings: {
            noAccountFoundForThisEmail: "No account found for this email",
            incorrectPassword: "Incorrect password",
          },
          errors: {
            signInError: "Sign in error",
          },
        },
        LoginPage: {
          headerMessage: "Login to your account",
          headerSlogan: "Turn chaos into progress.",
        },
        reauthModal: {
          title: "Session lost... re-auth required",
        },
      },
    },
    user: {
      components: {
        UserForm: {
          headers: {
            addUser: "Add user",
            updateUser: "Update user",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter user name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            email: {
              label: "Email",
              placeholder: "Enter user email address",
              errors: {
                alreadyExists: "Email already exists",
              },
            },
            password: {
              label: "Password",
              placeholder: "Enter user password",
              hidePasswordTooltipIcon: "hide password",
              showPasswordTooltipIcon: "show passwword",
            },
          },
          buttons: {
            changePassword: {
              label: "Change password",
            },
            signIn: {
              label: "Sign in",
            },
          },
          radios: {
            Permissions: {
              SuperUser: {
                label: "Administrator",
              },
              NormalUser: {
                label: "Normal user",
              },
            },
          },
          errors: {
            loadError: "There was a problem while loading the user data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified user",
            addError: "There was a problem while adding the user data",
            updateError: "There was a problem while updating the user data",
          },
        },
        ManageUsersPage: {
          header: {
            title: "Manage users",
          },
          pager: {
            totalItemsLabel: "Total users: {total}",
          },
          notifications: {
            userAdded: 'User "{name}" has been added',
            userUpdated: 'User "{name}" has been updated',
            userDeleted: 'User "{name}" has been deleted',
            userRestored: 'User "{name}" has been restored',
          },
          errors: {
            refreshError: "There was a problem while refreshing the user list",
            deleteError: "There was a problem while deleting the user",
            restoreError: "There was a problem while restoring the user",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified user",
          },
        },
        UsersTable: {
          header: {
            columns: {
              permissions: "Permissions",
              name: "Name",
              email: "Email",
              createdAt: "Created at",
              updatedAt: "Updated at",
              deletedAt: "Deleted at",
            },
            filters: {
              name: {
                placeholder: "search by name",
              },
              email: {
                placeholder: "search by email",
              },
            },
          },
          body: {
            columns: {
              permissions: {
                administrator: "Administrator",
                user: "User",
              },
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete user",
              message:
                'You are about to delete the user "{name}" from the system.',
            },
            undeleteConfirmation: {
              title: "Restore user",
              message:
                'You are about to restore the user "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No users found",
          },
        },
        UserPermissionsFilterSelector: {
          options: {
            any: "Any",
            onlyAdministrators: "Only administrators",
            onlyUsers: "Only users",
          },
        },
      },
    },
    role: {},
  },

  // common actions/operations
  Actions: "Actions",
  Add: "Add",
  Update: "Update",
  Save: "Save",
  Delete: "Delete",
  Restore: "Restore",
  Cancel: "Cancel",
  Refresh: "Refresh",
  "Move up": "Move up",
  "Move down": "Move down",
  "Do you want to continue ?": "Do you want to continue ?",

  // errors
  "Uncaught exception": "Uncaught exception",
  "Invalid API response code": "Invalid API response code",

  // text filter
  searchByNameDefaultPlaceholder: "search by name",
  searchByEmailDefaultPlaceholder: "search by email",

  // date filter
  "select date": "select date",
  "Any date": "Any date",
  "Custom date": "Custom date",
  Yesterday: "Yesterday",
  Today: "Today",
  Tomorrow: "Tomorrow",
  "Last week": "Last week",
  "This week": "This week",
  "Next week": "Next week",
  "Last month": "Last month",
  "This month": "This month",
  "Next month": "Next month",
  "Last year": "Last year",
  "This year": "This year",
  "Next year": "Next year",

  // user admin permission filter
  userAdminPermissionSelectorValueAll: "All",
  userAdminPermissionSelectorValueOnlyAdministrators: "Only administrators",
  userAdminPermissionSelectorValueOnlyUsers: "Only users",

  // role permission typefilter
  rolePermissionTypeSelectorValueAll: "All",
  rolePermissionTypeSelectorValueAllowed: "Only allowed",
  rolePermissionTypeSelectorValueNotAllowed: "Only not allowed",

  // user roles
  "Manage roles": "Manage roles",
  TotalRolesPagerLabel: "Total roles: {total}",
  RoleNameTableHeader: "Name",
  RoleProjectPermissionsTableHeader: "Project permissions",
  RoleTasksPermissionsTableHeader: "Task permissions",
  "Update project allowed": "Update project allowed",
  "Delete project allowed": "Delete project allowed",
  "View project allowed": "View project allowed",
  "Add task allowed": "Add task allowed",
  "Update task allowed": "Update task allowed",
  "Delete task allowed": "Delete task allowed",
  "View task allowed": "View task allowed",
  "Add role": "Add role",
  "Update role": "Update role",
  "Delete role": "Delete role",
  "No roles found": "No roles found",
  "Role permissions": "Role permissions",
  "Project permissions": "Project permissions",
  "Task permissions": "Task permissions",
  deleteRoleConfirmation:
    'You are about to delete the role "{name}" from the system.',
  "There was a problem while refreshing the role list":
    "There was a problem while refreshing the role list",
  "We couldn’t find the specified role": "We couldn’t find the specified role",
  "There was a problem while deleting the role":
    "There was a problem while deleting the role",
  "There was a problem while adding the role data":
    "There was a problem while adding the role data",
  "There was a problem while updating the role data":
    "There was a problem while updating the role data",
  "There was a problem while loading the role data":
    "There was a problem while loading the role data",
  roleAddedNotification: 'Role "{name}" has been added',
  roleUpdatedNotification: 'Role "{name}" has been updated',
  roleDeletedNotification: 'Role "{name}" has been deleted',
  roleFormNameLabel: "Name",
  roleFormNameFieldPlaceholder: "Enter role name",
  roleFormNameFieldEmptyError: "Name is empty",
  roleFormNameFieldTooLargeError: "Name is too large",
  roleFormNameFieldAlreadyExists: "Name already exists",

  // project types
  "Manage project types": "Manage project types",
  TotalProjectTypesPagerLabel: "Total project types: {total}",
  ProjectTypeNameTableHeader: "Name",
  "Add project type": "Add project type",
  "Update project type": "Update project type",
  "Delete project type": "Delete project type",
  "No project types found": "No project types found",
  deleteProjectTypeConfirmation:
    'You are about to delete the project type "{name}" from the system.',
  "There was a problem loading the project type data":
    "There was a problem loading the project type data",
  "There was a problem adding the project type data":
    "There was a problem adding the project type data",
  "There was a problem updating the project type data":
    "There was a problem updating the project type data",
  "There was a problem deleting the project type data":
    "There was a problem deleting the project type data",
  "There was a problem while adding the project type data":
    "There was a problem while adding the project type data",
  "There was a problem while updating the project type data":
    "There was a problem while updating the project type data",
  "There was a problem while loading the project type data":
    "There was a problem while loading the project type data",
  "We couldn’t find the specified project type":
    "We couldn’t find the specified project type",
  projectTypeFormNameFieldEmptyError: "Name is empty",
  projectTypeFormNameFieldTooLargeError: "Name is too large",
  projectTypeFormNameFieldAlreadyExists: "Name already exists",
  projectTypeFormNameFieldPlaceholder: "Enter project type name",
  projectTypeAddedNotification: 'Project type "{name}" has been added',
  projectTypeUpdatedNotification: 'Project type "{name}" has been updated',
  projectTypeDeletedNotification: 'Project type "{name}" has been deleted',

  // project statuses
  "Manage project statuses": "Manage project statuses",
  TotalProjectStatusesPagerLabel: "Total project statuses: {total}",
  ProjectStatusNameTableHeader: "Name",
  "Add project status": "Add project status",
  "Update project status": "Update project status",
  "Delete project status": "Delete project status",
  "No project statuses found": "No project statuses found",
  deleteProjectStatusConfirmation:
    'You are about to delete the project status "{name}" from the system.',
  "There was a problem loading the project status data":
    "There was a problem loading the project status data",
  "There was a problem adding the project status data":
    "There was a problem adding the project status data",
  "There was a problem updating the project status data":
    "There was a problem updating the project status data",
  "There was a problem deleting the project status data":
    "There was a problem deleting the project status data",
  "There was a problem while loading the project status data":
    "There was a problem while loading the project status data",
  "We couldn’t find the specified project status":
    "We couldn’t find the specified project status",
  projectStatusFormNameFieldEmptyError: "Name is empty",
  projectStatusFormNameFieldTooLargeError: "Name is too large",
  projectStatusFormNameFieldAlreadyExists: "Name already exists",
  projectStatusFormNameFieldPlaceholder: "Enter project status name",
  projectStatusAddedNotification: 'Project status "{name}" has been added',
  projectStatusUpdatedNotification: 'Project status "{name}" has been updated',
  projectStatusDeletedNotification: 'Project status "{name}" has been deleted',

  // project priorities
  "Manage project priorities": "Manage project priorities",
  TotalProjectPrioritiesPagerLabel: "Total project priorities: {total}",
  ProjectPriorityNameTableHeader: "Name",
  "Add project priority": "Add project priority",
  "Update project priority": "Update project priority",
  "Delete project priority": "Delete project priority",
  "No project priorities found": "No project priorities found",
  deleteProjectPriorityConfirmation:
    'You are about to delete the project priority "{name}" from the system.',
  "There was a problem loading the project priority data":
    "There was a problem loading the project priority data",
  "There was a problem adding the project priority data":
    "There was a problem adding the project priority data",
  "There was a problem updating the project priority data":
    "There was a problem updating the project priority data",
  "There was a problem deleting the project priority data":
    "There was a problem deleting the project priority data",
  "There was a problem while loading the project priority data":
    "There was a problem while loading the project priority data",
  "We couldn’t find the specified project priority":
    "We couldn’t find the specified project priority",
  projectPriorityFormNameFieldEmptyError: "Name is empty",
  projectPriorityFormNameFieldTooLargeError: "Name is too large",
  projectPriorityFormNameFieldAlreadyExists: "Name already exists",
  projectPriorityFormNameFieldPlaceholder: "Enter project priority name",
  projectPriorityAddedNotification: 'Project priority "{name}" has been added',
  projectPriorityUpdatedNotification:
    'Project priority "{name}" has been updated',
  projectPriorityDeletedNotification:
    'Project priority "{name}" has been deleted',

  // other
  Error: "Error",
  Name: "Name",
  Preview: "Preview",
  Color: "Color",

  // tooltips
  "show password": "show password",
  "hide password": "hide password",

  // notifications
  "Enable notifications": "Enable notifications",
  "Disable notifications": "Disable notifications",

  // hints
  "Open project page": "Open project page",
  "Switch to dark mode": "Switch to dark mode",
  "Switch to light mode": "Switch to light mode",
  "Switch to top navigation": "Switch to top navigation",
  "Switch to side navigation": "Switch to side navigation",

  // pager
  labelPageOfTotalPages: "Page {currentPage} of {totalPages}",
  "Total results:": "Total results:",
  "All results": "All results",
  labelSelectorResultsPage: "{number} results/page",
};
