export class IDate {
  // TODO: readonly msTimestamp ?, remove nulls
  msTimestamp: number | null;
  date: Date | null;

  constructor(msTimestamp: number | null) {
    this.msTimestamp = msTimestamp;
    this.date = msTimestamp !== null ? new Date(msTimestamp) : null;
  }

  hasValue() {
    return this.date !== null && this.msTimestamp !== null;
  }

  toLocaleString = () => {
    if (this.date === null && this.msTimestamp !== null) {
      this.date = new Date(this.msTimestamp);
    }
    return this.date?.toLocaleString();
  };
}
