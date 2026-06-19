import type { TimeTrackingResponse as TimeTrackingDTO } from "../types/dto";
import { UserBase } from "../../users/models/user";
import { IDate } from "../../../shared/types/idate";

export class TimeTracking {
  id: string | null;
  createdBy: UserBase;
  createdAt: IDate | null;
  summary: string;
  totalSeconds: number;

  constructor(data?: TimeTrackingDTO) {
    this.id = data?.id ?? null;
    this.createdBy = new UserBase(data?.createdBy);
    this.createdAt = data?.createdAt ? new IDate(data.createdAt) : null;
    this.summary = data?.summary ?? "";
    this.totalSeconds = data?.totalSeconds ?? 0;
  }

  toDTO(): TimeTrackingDTO {
    return {
      id: this.id ?? "",
      createdBy: this.createdBy.toDTO(),
      createdAt: this.createdAt?.msTimestamp ?? 0,
      summary: this.summary ?? "",
      totalSeconds: this.totalSeconds ?? 0,
    };
  }

  geti18nTimeParts() {
    const hours = Math.floor(this.totalSeconds / 3600);
    const minutes = Math.floor((this.totalSeconds % 3600) / 60);
    const seconds = this.totalSeconds % 60;

    return [
      { key: "shared.labels.time.hour", count: hours },
      { key: "shared.labels.time.minute", count: minutes },
      { key: "shared.labels.time.second", count: seconds },
    ].filter(({ count }) => count > 0);
  }
}
