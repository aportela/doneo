import type { VNodeChild } from "vue";

type TableHeaderColumnAlign = "left" | "right" | "center";

export interface TableHeaderColumn<T> {
  label: string;
  field: string;
  visible: boolean;
  sortable?: boolean;
  align?: TableHeaderColumnAlign;
  isFiltered?: () => boolean;
  render: (value: T) => VNodeChild;
}
