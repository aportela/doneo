import { h } from "vue";

import { NTag } from "naive-ui";

import { getNaiveUITagColorProperty } from "./color";

export const renderColoredTag = (
  name: string,
  color: string,
  bordered: boolean,
  className?: string,
) =>
  h(
    NTag,
    {
      bordered: bordered,
      color: getNaiveUITagColorProperty(color ?? "#888888"),
      class: className,
    },
    {
      default: () => name,
    },
  );

export const renderLabel = (value: string, className?: string) =>
  h("span", { class: className }, { default: () => value });
