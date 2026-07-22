import { h } from "vue";

import { NTag } from "naive-ui";

import { getNaiveUITagColorProperty } from "./color";

export const renderColoredTag = (name: string, color: string) =>
  h(
    NTag,
    {
      bordered: false,
      color: getNaiveUITagColorProperty(color ?? "#888888"),
    },
    {
      default: () => name,
    },
  );

export const renderLabel = (value: string) =>
  h("span", {}, { default: () => value });
