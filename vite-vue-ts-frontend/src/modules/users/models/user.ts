import type {
  UserResponse as UserDTO,
  UserBaseResponse as UserBaseDTO,
} from "../types/dto";
import { IDate } from "../../../shared/types/idate";

interface UserPermissions {
  isSuperUser: boolean;
}

export class UserBase {
  id: string;
  name: string;

  constructor(data: UserBaseDTO) {
    this.id = data.id;
    this.name = data.name;
  }

  toDTO(): UserBaseDTO {
    return {
      id: this.id,
      name: this.name,
    };
  }
}

export class User {
  id: string;
  name: string;
  email: string;
  password: string | null;
  permissions: UserPermissions;
  createdAt: IDate;
  updatedAt: IDate | null;
  deletedAt: IDate | null;

  constructor(data: UserDTO) {
    this.id = data.id;
    this.name = data.name;
    this.email = data.email;
    this.password = null;
    this.permissions = {
      isSuperUser: data.permissions.isSuperUser,
    };
    this.createdAt = new IDate(data.createdAt);
    this.updatedAt = data.updatedAt !== null ? new IDate(data.updatedAt) : null;
    this.deletedAt = data.deletedAt !== null ? new IDate(data.deletedAt) : null;
  }

  toDTO(): UserDTO {
    return {
      id: this.id,
      name: this.name,
      email: this.email,
      permissions: {
        isSuperUser: this.permissions.isSuperUser,
      },
      createdAt: this.createdAt.msTimestamp,
      updatedAt: this.updatedAt?.msTimestamp ?? null,
      deletedAt: this.deletedAt?.msTimestamp ?? null,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
export const MAX_EMAIL_LENGTH = 255;

export const MIN_PASSWORD_LENGTH = 4;
