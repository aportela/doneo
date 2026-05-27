import type { RoleResponse as RoleDTO } from "../types/dto";

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

  constructor(data: RoleDTO) {
    this.id = data.id;
    this.name = data.name;
    this.permissions = data.permissions;
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
}

export const MAX_NAME_LENGTH = 32;
