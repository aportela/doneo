import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectAttachmentsTableFilters {
  name: string;
  createdByUserId: string | null;
  createdAt: TimestampRange;
  contentType: string | null;
}
