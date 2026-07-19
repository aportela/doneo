import type {
  UserResponse as UserDTO,
  UserBaseResponse as UserBaseDTO,
  UserPermissions,
  AddRequest,
  UpdateRequest,
} from "../types/dto";
import { IDate } from "../../../shared/types/idate";

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

export class User extends UserBase {
  email: string;
  password: string;
  permissions: UserPermissions;
  createdAt: IDate;
  updatedAt: IDate | null;
  deletedAt: IDate | null;

  constructor(data?: UserDTO) {
    super(data);

    this.email = data?.email ?? "";
    this.password = "";
    this.permissions = {
      isSuperUser: data?.permissions?.isSuperUser ?? false,
    };
    this.createdAt = data?.createdAt
      ? new IDate(data.createdAt)
      : new IDate(Date.now());
    this.updatedAt = data?.updatedAt ? new IDate(data.updatedAt) : null;
    this.deletedAt = data?.deletedAt ? new IDate(data.deletedAt) : null;
  }

  override toDTO(): UserDTO {
    return {
      ...super.toDTO(),
      email: this.email,
      permissions: {
        isSuperUser: this.permissions.isSuperUser,
      },
      createdAt: this.createdAt.msTimestamp ?? Date.now(),
      updatedAt: this.updatedAt?.msTimestamp ?? null,
      deletedAt: this.deletedAt?.msTimestamp ?? null,
    };
  }

  toAddUserRequestPayload(): AddRequest {
    return {
      name: this.name,
      email: this.email,
      password: this.password,
      permissions: {
        isSuperUser: this.permissions.isSuperUser,
      },
    };
  }

  toUpdateUserRequestPayload(): UpdateRequest {
    return {
      id: this.id,
      name: this.name,
      email: this.email,
      password: this.password ? this.password : undefined,
      permissions: {
        isSuperUser: this.permissions.isSuperUser,
      },
    };
  }
}

export const MAX_NAME_LENGTH = 32;
export const MAX_EMAIL_LENGTH = 255;
export const MIN_PASSWORD_LENGTH = 4;
