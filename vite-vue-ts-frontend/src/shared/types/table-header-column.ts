type TableHeaderColumnAlign = "left" | "right" | "center";

export interface TableHeaderColumn {
  label: string;
  field: string;
  sortable?: boolean;
  align?: TableHeaderColumnAlign;
  visible?: boolean;
  isFiltered?: () => boolean;
}
