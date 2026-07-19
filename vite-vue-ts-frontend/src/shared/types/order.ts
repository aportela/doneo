export type SortDirection = "ASC" | "DESC";

export interface Order {
  field: string;
  direction: SortDirection;
}
