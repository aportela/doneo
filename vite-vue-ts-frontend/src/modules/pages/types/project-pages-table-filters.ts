import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectPagesTableFilters {
  title: string | null;
  userId: string | null;
  createdAt: TimestampRange;
  updatedAt: TimestampRange;
}
