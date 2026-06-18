import dayjs from "dayjs";

import { createStorageEntry } from "../composables/localStorage";

const localStorageDatetimeFormatMask = createStorageEntry<string | null>(
  "datetimeFormatMask",
  "YYYY-MM-DD HH:MM:SS",
);

// TODO: use store for real time changes without refresh page ???
const currentDatetimeMask =
  localStorageDatetimeFormatMask.get() ?? "YYYY-MM-DD HH:MM:SS";

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

  clear() {
    this.msTimestamp = null;
    this.date = null;
  }

  toLocaleString = () => {
    if (this.date === null && this.msTimestamp !== null) {
      this.date = new Date(this.msTimestamp);
    }
    return this.date?.toLocaleString();
  };

  toCustomMaskString = () => {
    if (this.date === null && this.msTimestamp !== null) {
      this.date = new Date(this.msTimestamp);
    }
    return dayjs(this.date).format(currentDatetimeMask);
  };
}
