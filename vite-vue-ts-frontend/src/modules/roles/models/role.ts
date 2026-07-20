import {
  type AddRequest,
  type RoleResponse as RoleDTO,
  type UpdateRequest,
  getDefaultPermissions,
} from "../types/dto";

type Permissions = {
  allowUpdateProject: boolean;
  allowDeleteProject: boolean;
  allowViewProject: boolean;
  allowAddTask: boolean;
  allowUpdateTask: boolean;
  allowDeleteTask: boolean;
  allowViewTask: boolean;
};

export class Role {
  id: string;
  name: string;
  permissions: Permissions;

  constructor(data?: RoleDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.permissions = data?.permissions ?? getDefaultPermissions();
  }

  toDTO(): RoleDTO {
    return {
      id: this.id,
      name: this.name,
      permissions: {
        allowUpdateProject: this.permissions.allowUpdateProject,
        allowDeleteProject: this.permissions.allowDeleteProject,
        allowViewProject: this.permissions.allowViewProject,
        allowAddTask: this.permissions.allowAddTask,
        allowUpdateTask: this.permissions.allowUpdateTask,
        allowDeleteTask: this.permissions.allowDeleteTask,
        allowViewTask: this.permissions.allowViewTask,
      },
    };
  }

  toAddRoleRequestPayload(): AddRequest {
    return {
      name: this.name,
      permissions: this.permissions,
    };
  }

  toUpdateRoleRequestPayload(): UpdateRequest {
    return {
      id: this.id,
      name: this.name,
      permissions: this.permissions,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
