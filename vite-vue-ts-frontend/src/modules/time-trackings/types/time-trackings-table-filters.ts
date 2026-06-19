import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface TimeTrackingsTableFilters {
  createdByUserId: string | null;
  createdAt: TimestampRange;
  summary: string;
  // TODO: totalSeconds
}
