/**
 * -----------------------------------------------------------------------------
 * @generated
 *
 * This file was generated with AI assistance and may have been modified
 * afterwards by human contributors.
 *
 * Please review changes before production use.
 * Manual edits may be overwritten if the file is regenerated.
 * -----------------------------------------------------------------------------
 */
import type { GlobalThemeOverrides } from "naive-ui";

type CatppuccinPalette = {
  background: string;
  surface: string;
  surface0: string;
  overlay: string;

  text: string;
  subtext: string;

  rosewater: string;
  flamingo: string;
  pink: string;
  mauve: string;
  red: string;
  maroon: string;
  peach: string;
  yellow: string;
  green: string;
  teal: string;
  sky: string;
  sapphire: string;
  blue: string;
  lavender: string;
};

const mocha: CatppuccinPalette = {
  background: "#1E1E2E",
  surface: "#313244",
  surface0: "#45475A",
  overlay: "#6C7086",

  text: "#CDD6F4",
  subtext: "#BAC2DE",

  rosewater: "#F5E0DC",
  flamingo: "#F2CDCD",
  pink: "#F5C2E7",
  mauve: "#CBA6F7",
  red: "#F38BA8",
  maroon: "#EBA0AC",
  peach: "#FAB387",
  yellow: "#F9E2AF",
  green: "#A6E3A1",
  teal: "#94E2D5",
  sky: "#89DCEB",
  sapphire: "#74C7EC",
  blue: "#89B4FA",
  lavender: "#B4BEFE",
};

const macchiato: CatppuccinPalette = {
  background: "#24273A",
  surface: "#363A4F",
  surface0: "#494D64",
  overlay: "#6E738D",

  text: "#CAD3F5",
  subtext: "#B8C0E0",

  rosewater: "#F4DBD6",
  flamingo: "#F0C6C6",
  pink: "#F5BDE6",
  mauve: "#C6A0F6",
  red: "#ED8796",
  maroon: "#EE99A0",
  peach: "#F5A97F",
  yellow: "#EED49F",
  green: "#A6DA95",
  teal: "#8BD5CA",
  sky: "#91D7E3",
  sapphire: "#7DC4E4",
  blue: "#8AADF4",
  lavender: "#B7BDF8",
};

const frappe: CatppuccinPalette = {
  background: "#303446",
  surface: "#414559",
  surface0: "#51576D",
  overlay: "#737994",

  text: "#C6D0F5",
  subtext: "#B5BFE2",

  rosewater: "#F2D5CF",
  flamingo: "#EEBEBE",
  pink: "#F4B8E4",
  mauve: "#CA9EE6",
  red: "#E78284",
  maroon: "#EA999C",
  peach: "#EF9F76",
  yellow: "#E5C890",
  green: "#A6D189",
  teal: "#81C8BE",
  sky: "#99D1DB",
  sapphire: "#85C1DC",
  blue: "#8CAAEE",
  lavender: "#BABBF1",
};

const latte: CatppuccinPalette = {
  background: "#EFF1F5",
  surface: "#E6E9EF",
  surface0: "#CCD0DA",
  overlay: "#9CA0B0",

  text: "#4C4F69",
  subtext: "#5C5F77",

  rosewater: "#DC8A78",
  flamingo: "#DD7878",
  pink: "#EA76CB",
  mauve: "#8839EF",
  red: "#D20F39",
  maroon: "#E64553",
  peach: "#FE640B",
  yellow: "#DF8E1D",
  green: "#40A02B",
  teal: "#179299",
  sky: "#04A5E5",
  sapphire: "#209FB5",
  blue: "#1E66F5",
  lavender: "#7287FD",
};

const createCatppuccinTheme = (c: CatppuccinPalette): GlobalThemeOverrides => {
  return {
    common: {
      primaryColor: c.mauve,
      primaryColorHover: c.lavender,
      primaryColorPressed: c.pink,
      primaryColorSuppl: c.mauve,

      infoColor: c.blue,
      successColor: c.green,
      warningColor: c.yellow,
      errorColor: c.red,

      bodyColor: c.background,
      cardColor: c.background,
      modalColor: c.background,

      borderColor: c.surface0,

      textColorBase: c.text,
      textColor1: c.text,
      textColor2: c.subtext,
      textColor3: c.overlay,

      inputColor: c.surface,
      placeholderColor: c.overlay,

      hoverColor: `${c.mauve}22`,
    },

    Layout: {
      color: c.background,
      siderColor: c.surface,
      headerColor: c.surface,
      footerColor: c.surface,
    },

    Card: {
      color: c.background,
      borderColor: c.surface0,
    },

    Button: {
      colorPrimary: c.mauve,
      colorHoverPrimary: c.lavender,
      colorPressedPrimary: c.pink,
      textColorPrimary: c.background,
    },

    Input: {
      color: c.surface,
      border: `1px solid ${c.surface0}`,
      borderHover: `1px solid ${c.blue}`,
      borderFocus: `1px solid ${c.mauve}`,
    },

    Menu: {
      color: c.background,
      itemColorHover: c.surface,
      itemColorActive: c.surface0,
      itemTextColor: c.text,
      itemTextColorActive: c.mauve,
    },

    DataTable: {
      thColor: c.surface,
      tdColor: c.background,
      borderColor: c.surface0,
    },

    Tabs: {
      barColor: c.mauve,
      tabTextColorActive: c.mauve,
    },

    Switch: {
      railColor: c.surface0,
      railColorActive: c.mauve,
    },
  };
};

export const catppuccinMocha = createCatppuccinTheme(mocha);
export const catppuccinMacchiato = createCatppuccinTheme(macchiato);
export const catppuccinFrappe = createCatppuccinTheme(frappe);
export const catppuccinLatte = createCatppuccinTheme(latte);
