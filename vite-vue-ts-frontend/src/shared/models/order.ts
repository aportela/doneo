import type { Order, SortDirection } from "../types/order";

export class ClientOrder implements Order {
  field: string;
  direction: SortDirection;

  constructor(field: string, direction: SortDirection) {
    this.field = field;
    this.direction = direction;
  }

  toggle = (field: string) => {
    if (field !== this.field) {
      this.field = field;
      this.direction = "ASC";
    } else {
      this.direction = this.direction === "ASC" ? "DESC" : "ASC";
    }
  };

  toDTO = (): Order => {
    return {
      field: this.field,
      direction: this.direction,
    };
  };
}
