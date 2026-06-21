import type { TimeTrackingResponse as TimeTrackingDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class TimeTracking {
  id: string | null;
  createdBy: UserBase;
  createdAt: IDate | null;
  summary: string;
  spentTime: number;

  constructor(data?: TimeTrackingDTO) {
    this.id = data?.id ?? null;
    this.createdBy = new UserBase(data?.createdBy);
    this.createdAt = data?.createdAt ? new IDate(data.createdAt) : null;
    this.summary = data?.summary ?? "";
    this.spentTime = data?.spentTime ?? 0;
  }

  toDTO(): TimeTrackingDTO {
    return {
      id: this.id ?? "",
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      summary: this.summary ?? "",
      spentTime: this.spentTime ?? 0,
    };
  }

  geti18nTimeParts() {
    const days = Math.floor(this.spentTime / 86400);
    const hours = Math.floor((this.spentTime % 86400) / 3600);
    const minutes = Math.floor((this.spentTime % 3600) / 60);

    return [
      { key: "shared.labels.time.day", count: days },
      { key: "shared.labels.time.hour", count: hours },
      { key: "shared.labels.time.minute", count: minutes },
    ].filter(({ count }) => count > 0);
  }
}
