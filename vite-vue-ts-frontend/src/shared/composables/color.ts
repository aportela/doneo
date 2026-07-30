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

export const oklchToHex = (l: number, c: number, h: number): string => {
  const hRad = (h * Math.PI) / 180;

  const a = c * Math.cos(hRad);
  const b = c * Math.sin(hRad);

  // OKLCH -> OKLab
  const l_ = l + 0.3963377774 * a + 0.2158037573 * b;
  const m_ = l - 0.1055613458 * a - 0.0638541728 * b;
  const s_ = l - 0.0894841775 * a - 1.291485548 * b;

  const l3 = l_ ** 3;
  const m3 = m_ ** 3;
  const s3 = s_ ** 3;

  // OKLab -> XYZ
  const x = 1.2268798734 * l3 - 0.5578149965 * m3 + 0.2813910502 * s3;
  const y = -0.0405757626 * l3 + 1.1122868294 * m3 - 0.0717110667 * s3;
  const z = -0.0763729497 * l3 - 0.4214933239 * m3 + 1.5869240244 * s3;

  // XYZ -> RGB
  let r = 3.2404542 * x - 1.5371385 * y - 0.4985314 * z;
  let g = -0.969266 * x + 1.8760108 * y + 0.041556 * z;
  let b2 = 0.0556434 * x - 0.2040259 * y + 1.0572252 * z;

  const gamma = (v: number) =>
    v <= 0.0031308 ? 12.92 * v : 1.055 * Math.pow(v, 1 / 2.4) - 0.055;

  r = Math.min(1, Math.max(0, gamma(r)));
  g = Math.min(1, Math.max(0, gamma(g)));
  b2 = Math.min(1, Math.max(0, gamma(b2)));

  return (
    "#" +
    [r, g, b2]
      .map((v) =>
        Math.round(v * 255)
          .toString(16)
          .padStart(2, "0"),
      )
      .join("")
  );
};
