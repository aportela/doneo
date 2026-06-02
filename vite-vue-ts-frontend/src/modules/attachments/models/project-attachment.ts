import type { AttachmentResponse as AttachmentDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class ProjectAttachment {
  id: string | null;
  createdBy: UserBase;
  createdAt: IDate;
  name: string;
  contentType: string;
  size: number;

  constructor(data?: AttachmentDTO) {
    this.id = data?.id ?? null;
    this.createdBy = new UserBase(data?.createdBy);
    this.createdAt = new IDate(data?.createdAt ?? null);
    this.name = data?.name ?? "";
    this.contentType = data?.contentType ?? "";
    this.size = data?.size ?? 0;
  }

  getDownloadURL = (projectId: string): string => {
    return `/api/projects/${projectId}/attachments/${this.id}`;
  };

  getPreviewURL = (projectId: string): string => {
    return `/api/projects/${projectId}/attachments/${this.id}`;
  };

  toDTO(): AttachmentDTO {
    return {
      id: this.id ?? "",
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      name: this.name,
      contentType: this.contentType,
      size: this.size,
    };
  }
}
