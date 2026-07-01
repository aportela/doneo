const DAY_SECONDS = 86400;
const HOUR_SECONDS = 3600;
const MINUTE_SECONDS = 60;

export { DAY_SECONDS, HOUR_SECONDS, MINUTE_SECONDS };

export const getSecondsDatetimeParts = (seconds: number) => {
  return {
    days: Math.floor(seconds / DAY_SECONDS),
    hours: Math.floor((seconds % DAY_SECONDS) / HOUR_SECONDS),
    minutes: Math.floor((seconds % HOUR_SECONDS) / MINUTE_SECONDS),
  };
};
