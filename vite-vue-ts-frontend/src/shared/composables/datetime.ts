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

export const formatDuration = (spentTime: number): string => {
  // TODO: i18N
  const days = Math.floor(spentTime / 86400);
  const hours = Math.floor((spentTime % 86400) / 3600);
  const minutes = Math.floor((spentTime % 3600) / 60);
  const seconds = spentTime % 60;
  const parts: string[] = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (minutes > 0) parts.push(`${minutes}m`);
  if (seconds > 0 || parts.length === 0) parts.push(`${seconds}s`);
  return parts.join(" ");
};
