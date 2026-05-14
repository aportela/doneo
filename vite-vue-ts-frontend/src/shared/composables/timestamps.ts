type WeekStart = 0 | 1 | 2 | 3 | 4 | 5 | 6;

type RangeType =
  | "today"
  | "yesterday"
  | "tomorrow"
  | "this_week"
  | "last_week"
  | "next_week"
  | "this_month"
  | "last_month"
  | "next_month"
  | "this_year"
  | "last_year"
  | "next_year";

interface RangeOptions {
  weekStartsOn?: WeekStart;
}

const startOfDay = (date: Date) => {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  return d;
};

const endOfDay = (date: Date) => {
  const d = new Date(date);
  d.setHours(23, 59, 59, 999);
  return d;
};

const getDayRange = (offsetDays: number) => {
  const base = new Date();
  base.setDate(base.getDate() + offsetDays);

  return {
    from: startOfDay(base).getTime(),
    to: endOfDay(base).getTime(),
  };
};

const getWeekRange = (offsetWeeks: number, weekStartsOn: WeekStart = 1) => {
  const now = new Date();

  const currentDay = now.getDay();
  let diffToStart = currentDay - weekStartsOn;

  if (diffToStart < 0) diffToStart += 7;

  const start = new Date(now);
  start.setDate(now.getDate() - diffToStart + offsetWeeks * 7);

  const end = new Date(start);
  end.setDate(start.getDate() + 6);

  return {
    from: startOfDay(start).getTime(),
    to: endOfDay(end).getTime(),
  };
};

const getMonthRange = (offsetMonths: number) => {
  const now = new Date();

  const start = new Date(now.getFullYear(), now.getMonth() + offsetMonths, 1);
  const end = new Date(now.getFullYear(), now.getMonth() + offsetMonths + 1, 0);

  return {
    from: startOfDay(start).getTime(),
    to: endOfDay(end).getTime(),
  };
};

const getYearRange = (offsetYears: number) => {
  const now = new Date();

  const start = new Date(now.getFullYear() + offsetYears, 0, 1);
  const end = new Date(now.getFullYear() + offsetYears, 11, 31);

  return {
    from: startOfDay(start).getTime(),
    to: endOfDay(end).getTime(),
  };
};

export const getRange = (type: RangeType, options: RangeOptions = {}) => {
  const weekStartsOn = options.weekStartsOn ?? 1;

  switch (type) {
    case "yesterday":
      return getDayRange(-1);

    case "today":
      return getDayRange(0);

    case "tomorrow":
      return getDayRange(1);

    case "this_week":
      return getWeekRange(0, weekStartsOn);

    case "last_week":
      return getWeekRange(-1, weekStartsOn);

    case "next_week":
      return getWeekRange(1, weekStartsOn);

    case "this_month":
      return getMonthRange(0);

    case "last_month":
      return getMonthRange(-1);

    case "next_month":
      return getMonthRange(1);

    case "this_year":
      return getYearRange(0);

    case "last_year":
      return getYearRange(-1);

    case "next_year":
      return getYearRange(1);

    default:
      return getDayRange(0);
  }
};
