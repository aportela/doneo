import type { RoleResponse as RoleDTO } from "../types/dto";

type Permissions = {
  allowCreate: boolean;
  allowUpdate: boolean;
  allowDelete: boolean;
  allowView: boolean;
  allowList: boolean;
  allowExecute: boolean;
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
        allowCreate: this.permissions.allowCreate,
        allowUpdate: this.permissions.allowUpdate,
        allowDelete: this.permissions.allowDelete,
        allowView: this.permissions.allowView,
        allowList: this.permissions.allowList,
        allowExecute: this.permissions.allowExecute,
      },
    };
  }
}

export const maxNameLength = 32;
