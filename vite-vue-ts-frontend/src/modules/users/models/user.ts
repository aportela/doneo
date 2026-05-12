import type { UserResponse as UserDTO } from "../types/dto";

class IDate {
  msTimestamp: number;
  date: Date;

  constructor(msTimestamp: number) {
    this.msTimestamp = msTimestamp;
    this.date = new Date(msTimestamp);
  }

  toLocaleString = () => {
    return this.date.toLocaleString();
  };
}

export class User {
  id: string;
  name: string;
  email: string;
  password: string | null;
  isSuperUser: boolean;
  createdAt: IDate;
  updatedAt: IDate | null;
  deletedAt: IDate | null;
  avatarUrl: string | null;

  constructor(data: UserDTO) {
    this.id = data.id;
    this.name = data.name;
    this.email = data.email;
    this.password = null;
    this.isSuperUser = data.isSuperUser;
    this.createdAt = new IDate(data.createdAt);
    this.updatedAt = data.updatedAt !== null ? new IDate(data.updatedAt) : null;
    this.deletedAt = data.deletedAt !== null ? new IDate(data.deletedAt) : null;
    this.avatarUrl = data.avatarUrl;
  }

  toDTO(): UserDTO {
    return {
      id: this.id,
      name: this.name,
      email: this.email,
      isSuperUser: this.isSuperUser,
      createdAt: this.createdAt.msTimestamp,
      updatedAt: this.updatedAt?.msTimestamp ?? null,
      deletedAt: this.deletedAt?.msTimestamp ?? null,
      avatarUrl: this.avatarUrl ?? "",
    };
  }
}

export const maxNameLength = 32;
export const maxEmailLength = 255;

export const minPasswordLength = 4;
