import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectTasksTableFilters {
  slug: string | null;
  priorityId: string | null;
  statusId: string | null;
  summary: string;
  createdAt: TimestampRange;
  createdByUserId: string | null;
}
