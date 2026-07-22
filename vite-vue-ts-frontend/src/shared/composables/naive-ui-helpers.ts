import { h, type Component } from "vue";

import { NIcon, NTag } from "naive-ui";

import { getNaiveUITagColorProperty } from "./color";
import type { IconProps } from "@tabler/icons-vue";

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

export const renderIcon =
  (icon: Component, props: IconProps = {}) =>
  () =>
    h(NIcon, props, {
      default: () => h(icon),
    });
