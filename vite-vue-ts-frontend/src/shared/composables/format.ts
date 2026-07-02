// TODO: i18n for locale param
export const formatBytes = (bytes: number, locale = "es-ES") => {
  if (bytes === 0) return "0 B";

  const units = ["B", "KB", "MB", "GB", "TB"];
  const exponent = Math.floor(Math.log(bytes) / Math.log(1024));

  const value = bytes / 1024 ** exponent;

  return `${new Intl.NumberFormat(locale, {
    maximumFractionDigits: 2,
  }).format(value)} ${units[exponent]}`;
};
