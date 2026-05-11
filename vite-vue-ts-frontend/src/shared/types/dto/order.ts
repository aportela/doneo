export type SortType = "ASC" | "DESC";

export type Order = {
  field: string;
  sort: SortType;
};
