export default {
  layouts: {
    sidebarMenu: {
      options: {
        search: "Search...",
        home: "Overview",
        projects: "Projects",
        tasks: "Tasks",
        reports: "Reports",
        charts: "Charts",
        settings: "Settings",
        manageUsers: "Users",
        manageRoles: "Roles",
        projectSettings: "Project",
        manageProjectTypes: "Type",
        manageProjectPriorities: "Priority",
        manageProjectStatuses: "Status",
        taskSettings: "Task",
        manageTaskPriorities: "Priority",
        manageTaskStatuses: "Status",
        disableNotifications: "Disable notifications",
        enableNotifications: "Enable notifications",
        switchToLightTheme: "Light theme",
        switchToDarkTheme: "Dark theme",
        profile: "Profile",
        signOut: "Logout",
      },
    },
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
        RemoteAPIAlertModal: {
          buttons: {
            backToHome: {
              label: "Back to Home",
            },
          },
        },
      },
      inputs: {
        ToggleInput: {
          buttons: {
            confirm: {
              toolTip: "Confirm",
            },
            cancel: {
              toolTip: "Cancel",
            },
          },
        },
        EstimatedTimeInput: {
          inputs: {
            estimatedDays: {
              label: "Estimated days",
              placeholder: "Enter estimated time (in days)",
            },
            estimatedHours: {
              label: "Estimated hours",
              placeholder: "Enter estimated time (in hours)",
            },
            estimatedMinutes: {
              label: "Estimated minutes",
              placeholder: "Enter estimated time (in minutes)",
            },
          },
        },
        TimeSpentInput: {
          inputs: {
            daysSpent: {
              label: "Days spent",
              placeholder: "Enter time spent (in days)",
            },
            hoursSpent: {
              label: "Hours spent",
              placeholder: "Enter time time (in hours)",
            },
            minutesSpent: {
              label: "Minutes spent",
              placeholder: "Enter time spent (in minutes)",
            },
          },
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
        ProjectPermissionSelect: {
          options: {
            updateProjectAllowed: "Update project allowed",
            updateProjectDenied: "Update project denied",
            deleteProjectAllowed: "Delete project allowed",
            deleteProjectDenied: "Delete project denied",
            viewProjectAllowed: "View project allowed",
            viewProjectDenied: "View project denied",
            addTaskAllowed: "Add task allowed",
            addTaskDenied: "Add task denied",
          },
          placeholder: "search by permission",
        },
        TaskPermissionSelect: {
          options: {
            updateTaskAllowed: "Update task allowed",
            updateTaskDenied: "Update task denied",
            deleteTaskAllowed: "Delete task allowed",
            deleteTaskDenied: "Delete task denied",
            viewTaskAllowed: "View task allowed",
            viewTaskDenied: "View task denied",
          },
          placeholder: "search by permission",
        },
        ToggleTagSelector: {
          buttons: {
            confirm: {
              toolTip: "Confirm",
            },
            cancel: {
              toolTip: "Cancel",
            },
          },
        },
      },
      pickers: {
        datepicker: {
          buttons: {
            now: {
              label: "now",
            },
            clear: {
              label: "clear",
            },
          },
        },
        ToggleDateTimePicker: {
          placeholder: "Select date and time",
        },
      },
      table: {
        header: {
          columns: {
            actions: "Actions",
          },
        },
        filters: {
          button: {
            clearFilters: {
              label: "Clear filters",
            },
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
        timer: {
          enable: {
            toolTip: "Enable timer",
          },
          disable: {
            toolTip: "Disable timer",
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
      Open: {
        label: "Open",
      },
      Download: {
        label: "Download",
      },
      Preview: {
        label: "Preview",
      },
      Settings: {
        label: "Settings",
      },
    },
    errorMessages: {
      Error: "Error",
      uncaugthException: "Uncaught exception",
      invalidAPIResponseCode: "Invalid API response code",
      unauthorizedOperation:
        "You are not authorized to perform this operation.",
    },
    warningMessages: {
      fieldIsRequired: "Field is required",
      fieldHasInvalidFormat: "Field has invalid format",
      fieldExceedsMaxLength: "Field exceeds max length",
      fieldIsBelowMinimumLength: "Field is below minimum length",
    },
    labels: {
      time: {
        hour: "hour | hours",
        minute: "minute | minutes",
        second: "second | seconds",
      },
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
            permissions: {
              label: "Permissions",
              superUser: {
                label: "Administrator",
              },
              normalUser: {
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
            updateProjectDenied: "Update project denied",
            deleteProjectAllowed: "Delete project allowed",
            deleteProjectDenied: "Delete project denied",
            viewProjectAllowed: "View project allowed",
            viewProjectDenied: "View project denied",
            addTaskAllowed: "Add task allowed",
            addTaskDenied: "Add task denied",
            updateTaskAllowed: "UpdateTaskAllowed",
            updateTaskDenied: "UpdateTaskDenied",
            deleteTaskAllowed: "Delete task allowed",
            deleteTaskDenied: "Delete task denied",
            viewTaskAllowed: "View task allowed",
            viewTaskDenied: "View task denied",
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
                updateProjectDenied: "Update project denied",
                deleteProjectAllowed: "Delete project allowed",
                deleteProjectDenied: "Delete project denied",
                viewProjectAllowed: "View project allowed",
                viewProjectDenied: "View project denied",
                addTaskAllowed: "Add task allowed",
                addTaskDenied: "Add task denied",
                updateTaskAllowed: "Update task allowed",
                updateTaskDenied: "Update task denied",
                deleteTaskAllowed: "Delete task allowed",
                deleteTaskDenied: "Delete task denied",
                viewTaskAllowed: "View task allowed",
                viewTaskDenied: "View task denied",
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
    note: {
      components: {
        NoteItem: {
          dialogs: {
            deleteConfirmation: {
              title: "Delete note",
              message: "You are about to delete the note from the system.",
            },
          },
        },
        ProjectNotesTab: {
          notifications: {
            projectNoteAdded: "Note has been added",
            projectNoteUpdated: "Note has been updated",
            projectNoteDeleted: "Note has been deleted",
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project note list",
            deleteError: "There was a problem while deleting the project note",
          },
        },
        // TODO: task notes tab
      },
    },
    projectPriority: {
      components: {
        ProjectPriorityForm: {
          headers: {
            addProjectPriority: "Add project priority",
            updateProjectPriority: "Update project priority",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter project priority name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            preview: {
              label: "Preview",
            },
            index: {
              label: "Index",
              placeholder: "Enter project priority index",
              errors: {
                alreadyExists: "Index",
              },
            },
          },
          errors: {
            loadError:
              "There was a problem while loading the project priority data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project priority",
            addError:
              "There was a problem while adding the project priority data",
            updateError:
              "There was a problem while updating the project priority data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
            indexAlreadyExists: "Index already exists",
          },
        },
        ManageProjectPrioritiesPage: {
          header: {
            title: "Manage project priorities",
          },
          pager: {
            totalItemsLabel: "Total project priorities: {total}",
          },
          notifications: {
            projectPriorityAdded: 'Project priority "{name}" has been added',
            projectPriorityUpdated:
              'Project priority "{name}" has been updated',
            projectPriorityDeleted:
              'Project priority "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project priority list",
            deleteError:
              "There was a problem while deleting the project priority",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project priority",
          },
        },
        ProjectPrioritiesTable: {
          header: {
            columns: {
              name: "Name",
              index: "Index",
            },
          },
          filters: {
            name: {
              placeholder: "search by name",
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project priority",
              message:
                'You are about to delete the project priority "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No project priorities found",
          },
        },
      },
    },
    projectStatus: {
      components: {
        ProjectStatusForm: {
          headers: {
            addProjectStatus: "Add project status",
            updateProjectStatus: "Update project status",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter project status name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            preview: {
              label: "Preview",
            },
            index: {
              label: "Index",
              placeholder: "Enter project status index",
              errors: {
                alreadyExists: "Index",
              },
            },
            flags: {
              label: "Flags",
            },
          },
          errors: {
            loadError:
              "There was a problem while loading the project status data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project status",
            addError:
              "There was a problem while adding the project status data",
            updateError:
              "There was a problem while updating the project status data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
            indexAlreadyExists: "Index already exists",
          },
        },
        ManageProjectStatusesPage: {
          header: {
            title: "Manage project statuses",
          },
          pager: {
            totalItemsLabel: "Total project statuses: {total}",
          },
          notifications: {
            projectStatusAdded: 'Project status "{name}" has been added',
            projectStatusUpdated: 'Project status "{name}" has been updated',
            projectStatusDeleted: 'Project status "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project status list",
            deleteError:
              "There was a problem while deleting the project status",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project status",
          },
        },
        ProjectStatusesTable: {
          header: {
            columns: {
              name: "Name",
              index: "Index",
              flags: "Flags",
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
                hasDefaultStatusOnCreation:
                  "set default status value on project creation",
                hasNotdefaultStatusOnCreation:
                  "do not set default status value on project creation",
                hasFillEmptyStartDate:
                  "set project start date on state change if unset",
                hasNotFillEmptyStartDate:
                  "do not set project start date on state change if unset",
                hasSetStartDate: "set project start date on state change",
                hasNotSetStartDate:
                  "do not set project start date on state change",
                hasFillEmptyFinishDate:
                  "set project finish date on state change if unset",
                hasNotFillEmptyFinishDate:
                  "do not set project finish date on state change if unset",
                hasSetFinishDate: "set project finish date on state change",
                hasNotSetFinishDate:
                  "do not set project finish date on state change",
                hasUnsetFinishDateOnLeave:
                  "unset project finish date when leaving this status",
                hasNotUnsetFinishDateOnLeave:
                  "do not unset project finish date when leaving this status",
              },
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project status",
              message:
                'You are about to delete the project status "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No project statuses found",
          },
        },
      },
    },
    projectType: {
      components: {
        ProjectTypeForm: {
          headers: {
            addProjectType: "Add project type",
            updateProjectType: "Update project type",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter project type name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            preview: {
              label: "Preview",
            },
          },
          errors: {
            loadError:
              "There was a problem while loading the project type data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project type",
            addError: "There was a problem while adding the project type data",
            updateError:
              "There was a problem while updating the project type data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
          },
        },
        ManageProjectTypesPage: {
          header: {
            title: "Manage project types",
          },
          pager: {
            totalItemsLabel: "Total project types: {total}",
          },
          notifications: {
            projectTypeAdded: 'Project type "{name}" has been added',
            projectTypeUpdated: 'Project type "{name}" has been updated',
            projectTypeDeleted: 'Project type "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project type list",
            deleteError: "There was a problem while deleting the project type",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project type",
          },
        },
        ProjectTypesTable: {
          header: {
            columns: {
              name: "Name",
            },
          },
          filters: {
            name: {
              placeholder: "search by name",
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project type",
              message:
                'You are about to delete the project type "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No project types found",
          },
        },
      },
    },
    projectPermission: {
      components: {
        ProjectPermissionForm: {
          headers: {
            addProjectPermission: "Add project permission",
          },
          inputs: {
            user: {
              label: "User",
              placeholder: "Select user",
              errors: {
                alreadyExists: "User already exists",
              },
            },
            role: {
              label: "Role",
              placeholder: "Select role",
              errors: {
                alreadyExists: "Role already exists",
              },
            },
          },
          errors: {
            addError:
              "There was a problem while adding the project permission data",
          },
          warnings: {
            userAlreadyExists: "User already exists",
            roleAlreadyExists: "Role already exists",
          },
        },
        projectPermissionsTable: {
          header: {
            columns: {
              user: "User",
              role: "Role",
              projectPermissions: "Project permissions",
              taskPermissions: "Task permissions",
            },
          },
          filters: {
            user: {
              placeholder: "search by user",
            },
            role: {
              placeholder: "search by role",
            },
          },
          body: {
            columns: {
              permissionsHints: {
                updateProjectAllowed: "Update project allowed",
                updateProjectDenied: "Update project denied",
                deleteProjectAllowed: "Delete project allowed",
                deleteProjectDenied: "Delete project denied",
                viewProjectAllowed: "View project allowed",
                viewProjectDenied: "View project denied",
                addTaskAllowed: "Add task allowed",
                addTaskDenied: "Add task denied",
                updateTaskAllowed: "Update task allowed",
                updateTaskDenied: "Update task denied",
                deleteTaskAllowed: "Delete task allowed",
                deleteTaskDenied: "Delete task denied",
                viewTaskAllowed: "View task allowed",
                viewTaskDenied: "View task denied",
              },
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project permission",
              message:
                'You are about to delete the project permission with user "{user}" and role "{role}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No project permissions found",
          },
        },
        projectPermissionsTab: {
          notifications: {
            projectPermissionAdded:
              'Project permission of user "{user}" with role "{role}" has been added',
            projectPermissionDeleted:
              'Project permission of user "{user}" with role "{role}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project permission list",
            deleteError:
              "There was a problem while deleting the project permission",
          },
        },
      },
    },
    projectAttachment: {
      components: {
        projectAttachmentsTable: {
          header: {
            columns: {
              name: "Name",
              size: "Size",
              contentType: "Content type",
              createdAt: "Created at",
              createdBy: "Created by",
            },
          },
          filters: {
            name: {
              placeholder: "search by name",
            },
            user: {
              placeholder: "search by user",
            },
            contentType: {
              placeholder: "search by content type",
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project attachment",
              message:
                "You are about to delete the project attachment {name} ({size}) from the system.",
            },
          },
          warnings: {
            noItemsFound: "No project attachments found",
          },
        },
        ProjectAttachmentsTab: {
          notifications: {
            projectAttachmentAdded:
              'Project attachment "{name}" has been added',
            projectAttachmentDeleted:
              'Project attachment "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project attachment list",
            deleteError:
              "There was a problem while deleting the project attachment",
          },
        },
      },
    },
    projectHistoryOperation: {
      components: {
        ProjectHistoryOperationsTable: {
          header: {
            columns: {
              user: "User",
              operationDate: "Operation date",
              operationType: "Operation type",
            },
          },
          filters: {
            user: {
              placeholder: "search by user",
            },
          },
          warnings: {
            noItemsFound: "No project history operations found",
          },
        },
      },
    },
    project: {
      components: {
        NewProjectForm: {
          headers: {
            addProject: "Add project",
          },
          inputs: {
            slug: {
              label: "Slug",
              placeholder: "Enter project slug",
            },
            summary: {
              label: "Summary",
              placeholder: "Enter project summary",
            },
            description: {
              label: "Description",
              placeholder: "Enter project description",
            },
          },
          selectors: {
            projectType: {
              label: "Project type",
              placeholder: "Select project type",
            },
            projectPriority: {
              label: "Project priority",
              placeholder: "Select project priority",
            },
            projectStatus: {
              label: "Project status",
              placeholder: "Select project status",
            },
            switches: {
              openProjectAfterCreate: {
                label: "Open project after create",
              },
            },
          },
          errors: {
            addError: "There was a problem while adding the project data",
          },
        },
        ManageProjectsPage: {
          header: {
            title: "Manage projects",
          },
          pager: {
            totalItemsLabel: "Total projects: {total}",
          },
          notifications: {
            projectAdded: 'Project "{summary}" has been added',
            projectUpdated: 'Project "{summary}" has been updated',
            projectDeleted: 'Project "{summary}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the project list",
            deleteError: "There was a problem while deleting the project",
            restoreError: "There was a problem while restoring the project",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project",
          },
        },
        ProjectsTable: {
          header: {
            columns: {
              slug: "Slug",
              type: "Type",
              priority: "Priority",
              status: "Status",
              summary: "Summary",
              createdAt: "Created at",
              createdBy: "Created by",
            },
            filters: {
              slug: {
                placeholder: "search by slug",
              },
              type: {
                placeholder: "search by type",
              },
              priority: {
                placeholder: "search by priority",
              },
              status: {
                placeholder: "search by status",
              },
              summary: {
                placeholder: "search by summary",
              },
              creator: {
                placeholder: "search by creator",
              },
            },
          },
          body: {
            columns: {},
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete project",
              message:
                'You are about to delete the project "{summary}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No projects found",
          },
        },
        ProjectPage: {
          tabs: {
            permissions: {
              label:
                "Permissions | Permissions ({count}) | Permissions ({count})",
            },
            attachments: {
              label:
                "Attachments | Attachments ({count}) | Attachments ({count})",
            },
            notes: {
              label: "Notes | Notes ({count}) | Notes ({count})",
            },
            history: {
              label: "History | History ({count}) | History ({count})",
            },
            tasks: {
              label: "Tasks | Tasks ({count}) | Tasks ({count})",
            },
          },
          errors: {
            loadError: "There was a problem while loading the project data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified project",
            addError: "There was a problem while adding the project data",
            updateError: "There was a problem while updating the project data",
          },
        },
        ProjectTasksTab: {
          notifications: {
            taskAdded: 'Task "{summary}" has been added',
          },
        },
      },
    },
    taskStatus: {
      components: {
        TaskStatusForm: {
          headers: {
            addTaskStatus: "Add task status",
            updateTaskStatus: "Update task status",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter task status name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            preview: {
              label: "Preview",
            },
            index: {
              label: "Index",
              placeholder: "Enter task status index",
              errors: {
                alreadyExists: "Index",
              },
            },
            flags: {
              label: "Flags",
            },
          },
          errors: {
            loadError: "There was a problem while loading the task status data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task status",
            addError: "There was a problem while adding the task status data",
            updateError:
              "There was a problem while updating the task status data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
            indexAlreadyExists: "Index already exists",
          },
        },
        ManageTaskStatusesPage: {
          header: {
            title: "Manage task statuses",
          },
          pager: {
            totalItemsLabel: "Total task statuses: {total}",
          },
          notifications: {
            taskStatusAdded: 'Task status "{name}" has been added',
            taskStatusUpdated: 'Task status "{name}" has been updated',
            taskStatusDeleted: 'Task status "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the task status list",
            deleteError: "There was a problem while deleting the task status",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task status",
          },
        },
        TaskStatusesTable: {
          header: {
            columns: {
              name: "Name",
              index: "Index",
              flags: "Flags",
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
                hasDefaultStatusOnCreation:
                  "set default status value on task creation",
                hasNotdefaultStatusOnCreation:
                  "do not set default status value on task creation",
                hasFillEmptyStartDate:
                  "set task start date on state change if unset",
                hasNotFillEmptyStartDate:
                  "do not set task start date on state change if unset",
                hasSetStartDate: "set task start date on state change",
                hasNotSetStartDate:
                  "do not set task start date on state change",
                hasFillEmptyFinishDate:
                  "set task finish date on state change if unset",
                hasNotFillEmptyFinishDate:
                  "do not set task finish date on state change if unset",
                hasSetFinishDate: "set task finish date on state change",
                hasNotSetFinishDate:
                  "do not set task finish date on state change",
                hasUnsetFinishDateOnLeave:
                  "unset task finish date when leaving this status",
                hasNotUnsetFinishDateOnLeave:
                  "do not unset task finish date when leaving this status",
              },
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete task status",
              message:
                'You are about to delete the task status "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No task statuses found",
          },
        },
      },
    },
    taskPriority: {
      components: {
        TaskPriorityForm: {
          headers: {
            addTaskPriority: "Add task priority",
            updateTaskPriority: "Update task priority",
          },
          inputs: {
            name: {
              label: "Name",
              placeholder: "Enter task priority name",
              errors: {
                alreadyExists: "Name already exists",
              },
            },
            preview: {
              label: "Preview",
            },
            index: {
              label: "Index",
              placeholder: "Enter task priority index",
              errors: {
                alreadyExists: "Index",
              },
            },
          },
          errors: {
            loadError:
              "There was a problem while loading the task priority data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task priority",
            addError: "There was a problem while adding the task priority data",
            updateError:
              "There was a problem while updating the task priority data",
          },
          warnings: {
            nameAlreadyExists: "Name already exists",
            indexAlreadyExists: "Index already exists",
          },
        },
        ManageTaskPrioritiesPage: {
          header: {
            title: "Manage task priorities",
          },
          pager: {
            totalItemsLabel: "Total task priorities: {total}",
          },
          notifications: {
            taskPriorityAdded: 'Task priority "{name}" has been added',
            taskPriorityUpdated: 'Task priority "{name}" has been updated',
            taskPriorityDeleted: 'Task priority "{name}" has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the task priority list",
            deleteError: "There was a problem while deleting the task priority",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task priority",
          },
        },
        TaskPrioritiesTable: {
          header: {
            columns: {
              name: "Name",
              index: "Index",
            },
          },
          filters: {
            name: {
              placeholder: "search by name",
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete task priority",
              message:
                'You are about to delete the task priority "{name}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No task priorities found",
          },
        },
      },
    },
    task: {
      components: {
        NewTaskForm: {
          headers: {
            addTask: "Add task",
          },
          inputs: {
            summary: {
              label: "Summary",
              placeholder: "Enter task summary",
            },
            description: {
              label: "Description",
              placeholder: "Enter task description",
            },
          },
          selectors: {
            taskPriority: {
              label: "Task priority",
              placeholder: "Select task priority",
            },
            taskStatus: {
              label: "Task status",
              placeholder: "Select task status",
            },
            switches: {
              openTaskAfterCreate: {
                label: "Open task after create",
              },
            },
          },
          collapse: {
            estimatedTimeCollapseItem: {
              label: "Estimated time",
            },
          },
          errors: {
            addError: "There was a problem while adding the task data",
          },
        },
        ManageTasksPage: {
          header: {
            title: "Manage tasks",
          },
          pager: {
            totalItemsLabel: "Total tasks: {total}",
          },
        },
        TasksTable: {
          header: {
            columns: {
              slug: "Slug",
              priority: "Priority",
              status: "Status",
              summary: "Summary",
              createdAt: "Created at",
              createdBy: "Created by",
            },
            filters: {
              slug: {
                placeholder: "search by slug",
              },
              priority: {
                placeholder: "search by priority",
              },
              status: {
                placeholder: "search by status",
              },
              summary: {
                placeholder: "search by summary",
              },
              creator: {
                placeholder: "search by creator",
              },
            },
          },
          body: {
            columns: {},
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete task",
              message:
                'You are about to delete the task "{summary}" from the system.',
            },
          },
          warnings: {
            noItemsFound: "No tasks found",
          },
        },
        TaskPage: {
          tabs: {
            attachments: {
              label:
                "Attachments | Attachments ({count}) | Attachments ({count})",
            },
            notes: {
              label: "Notes | Notes ({count}) | Notes ({count})",
            },
            timeTrackings: {
              label:
                "Time trackings | Time trackings ({count}) | Time trackings ({count})",
            },
            history: {
              label: "History | History ({count}) | History ({count})",
            },
          },
          errors: {
            loadError: "There was a problem while loading the task data",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task",
            addError: "There was a problem while adding the task data",
            updateError: "There was a problem while updating the task data",
          },
        },
        ProjectTasksTab: {
          pager: {
            totalItemsLabel: "Total tasks: {total}",
          },
        },
        TimeTrackingsTab: {
          notifications: {
            timeTrackingAdded: 'Time tracking "{summary}" has been added',
            timeTrackingDeleted: 'Time tracking "{summary}"  has been deleted',
          },
          errors: {
            refreshError:
              "There was a problem while refreshing the time tracking list",
            deleteError: "There was a problem while deleting the time tracking",
          },
        },
      },
    },
    timeTracking: {
      components: {
        TimeTrackingForm: {
          headers: {
            addTimeTracking: "Add time tracking",
          },
          inputs: {
            summary: {
              label: "Summary",
              placeholder: "type summary",
            },
            errors: {
              loadError:
                "There was a problem while loading the time tracking data",
              // TODO: deduplicate
              notFoundError: "We couldn’t find the specified time tracking",
              addError:
                "There was a problem while adding the time tracking data",
              updateError:
                "There was a problem while updating the time tracking data",
            },
          },
        },
        TimeTrackingsTable: {
          header: {
            columns: {
              summary: "Summary",
              totalTime: "Total time",
              createdAt: "Created at",
              createdBy: "Created by",
            },
          },
          filters: {
            summary: {
              placeholder: "search by summary",
            },
            user: {
              placeholder: "search by user",
            },
            contentType: {
              placeholder: "search by content type",
            },
          },
          dialogs: {
            deleteConfirmation: {
              title: "Delete time tracking",
              message:
                "You are about to delete the time tracking {summary} from the system.",
            },
          },
          warnings: {
            noItemsFound: "No time trackings found",
          },
        },
        taskTimeTrackingsTab: {
          errors: {
            refreshError:
              "There was a problem while refreshing the task time tracking list",
            deleteError:
              "There was a problem while deleting the task time tracking",
            // TODO: deduplicate
            notFoundError: "We couldn’t find the specified task time tracking",
          },
          notifications: {
            timeTrackingAdded: 'Time tracking "{summary}" has been added',
            timeTrackingUpdated: 'Time tracking "{summary}" has been updated',
            timeTrackingDeleted: 'Time tracking "{summary}" has been deleted',
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

  // other
  Error: "Error",
  Name: "Name",
  Preview: "Preview",
  Color: "Color",
};
