import type { AttachmentResponse as AttachmentDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";
import { isAudio, isImage, isPDF } from "../../../shared/composables/fileUtils";

export class Attachment {
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
    return `/api/wc/attachments/project/${projectId}/attachment/${this.id}/download`;
  };

  getPreviewURL = (projectId: string): string => {
    return `/api/wc/attachments/project/${projectId}/attachment/${this.id}/inline`;
  };

  getBgDownloadURL = (projectId: string): string => {
    return `/wc/attachments/project/${projectId}/attachment/${this.id}/download`;
  };

  allowImagePreview = (): boolean => {
    return isImage(this.name);
  };

  allowAudioPreview = (): boolean => {
    // TODO: with cache
    /*
    return (
      isAudio(this.name) &&
      document.createElement("audio").canPlayType(this.contentType) !== ""
    );
    */
    return isAudio(this.name);
  };

  allowPDFPreview = (): boolean => {
    return isPDF(this.name);
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
