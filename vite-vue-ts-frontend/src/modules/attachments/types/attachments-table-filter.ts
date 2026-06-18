import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface AttachmentsTableFilters {
  name: string;
  createdByUserId: string | null;
  createdAt: TimestampRange;
  contentType: string | null;
}
