import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectHistoryOperationsTableFilters {
  userId: string | null;
  createdAt: TimestampRange;
}
