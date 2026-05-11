export const hexToRgba = (hex: string, alphaOverride?: number) => {
  if (!hex) return `rgba(0,0,0,1)`;
  let h = hex.replace("#", "");
  let r,
    g,
    b,
    a = 1;
  if (h.length === 8) {
    // RRGGBBAA
    r = parseInt(h.slice(0, 2), 16);
    g = parseInt(h.slice(2, 4), 16);
    b = parseInt(h.slice(4, 6), 16);
    a = parseInt(h.slice(6, 8), 16) / 255;
  } else {
    // RRGGBB
    r = parseInt(h.slice(0, 2), 16);
    g = parseInt(h.slice(2, 4), 16);
    b = parseInt(h.slice(4, 6), 16);
  }

  const alpha = alphaOverride ?? a;

  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
};

export const tagColor = (base: string) => {
  return {
    color: hexToRgba(base, 0.2),
    textColor: hexToRgba(base, 1),
    borderColor: hexToRgba(base, 0.5),
  };
};

export const getNaiveUITagColorProperty = (base: string) => {
  return {
    color: hexToRgba(base, 0.2),
    textColor: hexToRgba(base, 1),
    borderColor: hexToRgba(base, 0.5),
  };
};

const clamp = (v: number): number => {
  if (v < 0) return 0;
  if (v > 255) return 255;
  return Math.round(v);
};

export const hslToRgb = (
  h: number,
  s: number,
  l: number,
): [number, number, number] => {
  let r: number, g: number, b: number;

  const c = (1 - Math.abs(2 * l - 1)) * s;
  const x = c * (1 - Math.abs(((h / 60) % 2) - 1));
  const m = l - c / 2;

  if (h < 60) {
    [r, g, b] = [c, x, 0];
  } else if (h < 120) {
    [r, g, b] = [x, c, 0];
  } else if (h < 180) {
    [r, g, b] = [0, c, x];
  } else if (h < 240) {
    [r, g, b] = [0, x, c];
  } else if (h < 300) {
    [r, g, b] = [x, 0, c];
  } else {
    [r, g, b] = [c, 0, x];
  }

  return [clamp((r + m) * 255), clamp((g + m) * 255), clamp((b + m) * 255)];
};

export const generateRandomSoftHexColor = (): string => {
  const h = Math.floor(Math.random() * 360);
  const s = 0.45 + Math.random() * 0.25; // 45% - 70%
  const l = 0.55 + Math.random() * 0.2; // 55% - 75%

  const [r, g, b] = hslToRgb(h, s, l);

  const toHex = (n: number) => n.toString(16).padStart(2, "0").toUpperCase();

  return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
};
