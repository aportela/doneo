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
      selectors: {
        dateFilter: {
          options: {
            anyDate: "Any date",
            customDate: "Custom date",
            yesterday: "Yesterday",
            today: "Today",
            tomorrow: "Tomorrow",
            lastWeek: "Last week",
            thisWeek: "This week",
            nextWeek: "Next week",
            lastMonth: "Last month",
            thisMonth: "This month",
            nextMonth: "Next month",
            lastYear: "Last year",
            thisYear: "This year",
            nextYear: "Next year",
          },
          placeholder: "select date",
        },
      },
      table: {
        header: {
          columns: {
            actions: "Actions",
          },
        },
      },
      pager: {
        selector: {
          options: {
            allResults: "All results",
            nnnResultsPage: "{number} results/page",
          },
        },
        labels: {
          totalResults: "Total results:",
          currentPageOfTotal: "Page {currentPage} of {totalPages}",
        },
      },
      buttons: {
        github: {
          toolTip: "Open project page",
        },
        notifications: {
          enable: {
            toolTip: "Enable notifications",
          },
          disable: {
            toolTip: "Disable notifications",
          },
        },
        colorScheme: {
          darkMode: {
            toolTip: "Switch to dark mode",
          },
          lightMode: {
            toolTip: "Switch to light mode",
          },
        },
        navigationMode: {
          topNavigation: {
            toolTip: "Switch to top navigation",
          },
          sideNavigation: {
            toolTip: "Switch to side navigation",
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
    role: {
      components: {
        RoleForm: {
          headers: {
            addRole: "Add role",
            updateRole: "Update role",
            rolePermissions: "Role permissions",
            projectPermissions: "Project permissions",
            taskPermissions: "Task permissions",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter role name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
          },
          permissionSwitches: {
            updateProjectAllowed: "Update project allowed",
            deleteProjectAllowed: "Delete project allowed",
            viewProjectAllowed: "View project allowed",
            addTaskAllowed: "Add task allowed",
            updateTaskAllowed: "UpdateTaskAllowed",
            deleteTaskAllowed: "Delete task allowed",
            viewTaskAllowed: "View task allowed",
          },
          errors: {
            loadError: "There was a problem while loading the role data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified role",
            addError: "There was a problem while adding the role data",
            updateError: "There was a problem while updating the role data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
          },
        },
        ManageRolesPage: {
          header: {
            title: "Manage roles",
          },
          pager: {
            totalItemsLabel: "Total roles: {total}",
          },
          notifications: {
            roleAdded: 'Role "{name}" has been added',
            roleUpdated: 'Role "{name}" has been updated',
            roleDeleted: 'Role "{name}" has been deleted',
          },
          errors: {
            refreshError: "There was a problem while refreshing the role list",
            deleteError: "There was a problem while deleting the role",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified role",
          },
        },
        RolesTable: {
          header: {
            columns: {
              name: "Name",
              projectPermissions: "Project permissions",
              taskPermissions: "Task permissions",
            },
          },
          filters: {
            name: {
              placeholder: "search by name",
            },
          },
          body: {
            columns: {
              permissionsHints: {
                updateProjectAllowed: "Update project allowed",
                deleteProjectAllowed: "Delete project allowed",
                viewProjectAllowed: "View project allowed",
                addTaskAllowed: "Add task allowed",
                updateTaskAllowed: "Update task allowed",
                deleteTaskAllowed: "Delete task allowed",
                viewTaskAllowed: "View task allowed",
              },
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete role",
              message:
                'You are about to delete the role "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No roles found",
          },
        },
      },
    },
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
};
