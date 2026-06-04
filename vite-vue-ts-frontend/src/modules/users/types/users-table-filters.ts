import type { UserPermissionFilter } from "./user-admin-permission-filter";
import type { TimestampRange } from "../../../shared/composables/timestamps";

export interface UsersTableFilters {
  permissions: UserPermissionFilter;
  name: string;
  email: string;
  createdAt: TimestampRange;
  updatedAt: TimestampRange;
  deletedAt: TimestampRange;
}
