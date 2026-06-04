type TableHeaderColumnAlign = "left" | "right" | "center";

export interface TableHeaderColumn {
  label: string;
  field: string;
  visible: boolean;
  sortable?: boolean;
  align?: TableHeaderColumnAlign;
  isFiltered?: () => boolean;
}
