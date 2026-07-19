import type { SortOrder } from "../common";

export class SortClass {
  field: string;
  order: SortOrder;

  constructor(field: string, order: SortOrder) {
    this.field = field;
    this.order = order;
  }

  toggleSort = (field: string) => {
    if (field !== this.field) {
      this.field = field;
      this.order = "ASC";
    } else {
      this.order = this.order === "ASC" ? "DESC" : "ASC";
    }
  };
}
