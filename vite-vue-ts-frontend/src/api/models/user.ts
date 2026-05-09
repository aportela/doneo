import type { User as UserDTO } from "../types/dto/user";

class IDate {
  msTimestamp: number;
  date: Date;

  constructor(msTimestamp: number) {
    this.msTimestamp = msTimestamp;
    this.date = new Date(msTimestamp);
  }
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
  avatarURL: string | null;

  constructor(data: UserDTO) {
    this.id = data.id;
    this.name = data.name;
    this.email = data.email;
    this.password = null;
    this.isSuperUser = data.isSuperUser;
    this.createdAt = new IDate(data.createdAt);
    this.updatedAt = data.updatedAt !== null ? new IDate(data.updatedAt) : null;
    this.deletedAt = data.deletedAt !== null ? new IDate(data.deletedAt) : null;
    this.avatarURL = data.avatarURL;
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
      avatarURL: this.avatarURL ?? "",
    };
  }
}
