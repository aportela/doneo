import type { PageResponse as PageDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class Page {
  id: string;
  createdBy: UserBase;
  createdAt: IDate;
  updatedAt: IDate | null;
  title: string;
  body: string;

  constructor(data?: PageDTO) {
    this.id = data?.id ?? "";
    this.createdBy = new UserBase(data?.createdBy);
    this.createdAt = new IDate(data?.createdAt ?? null);
    this.updatedAt = data?.updatedAt ? new IDate(data.updatedAt) : null;
    this.title = data?.title ?? "";
    this.body = data?.body ?? "";
  }

  toDTO(): PageDTO {
    return {
      id: this.id,
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      updatedAt: this.updatedAt?.msTimestamp ?? null,
      title: this.title,
      body: this.body,
    };
  }
}
