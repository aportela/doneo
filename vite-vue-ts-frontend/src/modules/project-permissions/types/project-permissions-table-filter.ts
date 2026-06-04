import type { ProjectPermissionSelectValue } from "../../../shared/types/project-permission-select-value";
import type { TaskPermissionSelectValue } from "../../../shared/types/task-permission-select-value";

export interface ProjectPermissionsTableFilters {
  userId: string | null;
  roleId: string | null;
  projectPermission: ProjectPermissionSelectValue | null;
  taskPermission: TaskPermissionSelectValue | null;
}
