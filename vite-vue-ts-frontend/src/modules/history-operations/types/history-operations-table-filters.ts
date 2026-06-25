import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface HistoryOperationsTableFilters {
  userId: string | null;
  createdAt: TimestampRange;
  operationType: number | null;
}
