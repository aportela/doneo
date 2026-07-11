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

  constructor(data?: UserBaseDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
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
  password: string;
  permissions: UserPermissions;
  createdAt: IDate;
  updatedAt: IDate | null;
  deletedAt: IDate | null;

  constructor(data?: UserDTO) {
    this.id = data?.id ?? "";
    this.name = data?.name ?? "";
    this.email = data?.email ?? "";
    this.password = "";
    this.permissions = {
      isSuperUser: data?.permissions.isSuperUser ?? false,
    };
    this.createdAt = data?.createdAt
      ? new IDate(data.createdAt)
      : new IDate(new Date().getTime());
    this.updatedAt = data?.updatedAt ? new IDate(data.updatedAt) : null;
    this.deletedAt = data?.deletedAt ? new IDate(data.deletedAt) : null;
  }

  toDTO(): UserDTO {
    return {
      id: this.id,
      name: this.name,
      email: this.email,
      permissions: {
        isSuperUser: this.permissions.isSuperUser ?? false,
      },
      createdAt: this.createdAt.msTimestamp ?? new Date().getTime(),
      updatedAt: this.updatedAt?.msTimestamp ?? null,
      deletedAt: this.deletedAt?.msTimestamp ?? null,
    };
  }
}

export const MAX_NAME_LENGTH = 32;
export const MAX_EMAIL_LENGTH = 255;
export const MIN_PASSWORD_LENGTH = 4;
