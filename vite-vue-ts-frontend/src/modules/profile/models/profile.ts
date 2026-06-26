import type { ProfileResponse as ProfileDTO } from "../types/dto";
import { IDate } from "../../../shared/types/idate";

export class Profile {
  id: string | null;
  name: string | null;
  email: string | null;
  createdAt: IDate | null;
  updatedAt: IDate | null;

  constructor(data?: ProfileDTO) {
    this.id = data?.id ?? null;
    this.name = data?.name ?? null;
    this.email = data?.email ?? null;
    this.createdAt = data?.createdAt ? new IDate(data.createdAt) : null;
    this.updatedAt = data?.updatedAt ? new IDate(data.updatedAt) : null;
  }

  toDTO(): ProfileDTO {
    return {
      id: this.id ?? "",
      name: this.name ?? "",
      email: this.email ?? "",
      createdAt: this.createdAt?.msTimestamp ?? 0,
      updatedAt: this.updatedAt?.msTimestamp ?? null,
    };
  }
}
