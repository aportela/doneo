import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectsTableFilters {
  typeId: string | null;
  priorityId: string | null;
  summary: string;
  createdAt: TimestampRange;
  createdByUserId: string | null;
}
