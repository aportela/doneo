export type SortDirection = "ASC" | "DESC";

export type Order = {
  field: string;
  direction: SortDirection;
};
