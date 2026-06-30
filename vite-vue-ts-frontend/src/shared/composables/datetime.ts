export const geti18nTimeParts = (seconds: number) => {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  return [
    { key: "shared.labels.time.day", count: days },
    { key: "shared.labels.time.hour", count: hours },
    { key: "shared.labels.time.minute", count: minutes },
  ].filter(({ count }) => count > 0);
};

export const defaultDateTimeMask = "YYYY-MM-DD HH:MM:ss";
