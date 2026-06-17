import type { AttachmentResponse as AttachmentDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";
import { isImage } from "../../../shared/composables/fileUtils";

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
    return `/api/wc/attachments/project/${projectId}/attachment/${this.id}`;
  };

  getAxiosDownloadURL = (projectId: string): string => {
    return `/wc/attachments/project/${projectId}/attachment/${this.id}`;
  };

  getPreviewURL = (projectId: string): string => {
    return `/api/wc/attachments/project/${projectId}/attachment/${this.id}`;
  };

  allowImagePreview = (): boolean => {
    return isImage(this.name) ?? false;
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
