import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface ProjectsTableFilters {
  slug: string;
  typeId: string | null;
  priorityId: string | null;
  statusId: string | null;
  summary: string;
  createdAt: TimestampRange;
  createdByUserId: string | null;
}
