import type { ProjectPermissionSelectValue } from "../../../shared/types/project-permission-select-value";
import type { TaskPermissionSelectValue } from "../../../shared/types/task-permission-select-value";

export interface RolesTableFilters {
  name: string;
  projectPermission: ProjectPermissionSelectValue | null;
  taskPermission: TaskPermissionSelectValue | null;
}
