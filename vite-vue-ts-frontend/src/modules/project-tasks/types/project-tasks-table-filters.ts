import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface TasksTableFilters {
  slug: string | null;
  priorityId: string | null;
  statusId: string | null;
  summary: string;
  createdAt: TimestampRange;
  createdByUserId: string | null;
}
