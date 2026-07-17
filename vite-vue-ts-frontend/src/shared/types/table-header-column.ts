import type { VNodeChild } from "vue";

export interface TableHeaderColumn<T> {
  label: string;
  field: string;
  visible: boolean;
  sortable?: boolean;
  align?: "left" | "right" | "center";
  isFiltered?: () => boolean; // TODO: remove func ?
  render: (value: T) => VNodeChild;
}
