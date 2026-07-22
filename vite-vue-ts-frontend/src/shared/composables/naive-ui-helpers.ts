import { h, type Component, type VNodeChild } from "vue";

import { NIcon, NTag } from "naive-ui";

import { getNaiveUITagColorProperty } from "./color";
import type { IconProps } from "@tabler/icons-vue";

export const renderColoredTag = (
  name: string,
  color: string,
  bordered: boolean,
  className?: string,
): VNodeChild =>
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

export const renderLabel = (value: string, className?: string): VNodeChild =>
  className ? h("span", { class: className }, { default: () => value }) : value;

export const renderIcon =
  (icon: Component, props: IconProps = {}) =>
  (): VNodeChild =>
    h(NIcon, props, {
      default: () => h(icon),
    });
