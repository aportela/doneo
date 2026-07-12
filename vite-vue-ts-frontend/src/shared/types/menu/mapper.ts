import { h, type Component } from "vue";
import { RouterLink } from "vue-router";

import type { MenuMixedOption } from "naive-ui/es/menu/src/interface";

import type { AppMenuItem } from "./types";
import { menuIcons, MENU_ICON_SIZE } from "./icons";

const renderMenuIcon = (iconName?: keyof typeof menuIcons) => {
  if (!iconName) {
    return undefined;
  }

  const Icon = menuIcons[iconName];

  return () =>
    h(Icon as Component, {
      size: MENU_ICON_SIZE,
    });
};

const mapChildren = (
  children?: AppMenuItem[],
): MenuMixedOption[] | undefined => {
  return children
    ?.map(mapMenuItem)
    .filter((item): item is MenuMixedOption => item !== null);
};

const mapMenuItem = (item: AppMenuItem): MenuMixedOption | null => {
  if (item.show === false) {
    return null;
  }

  switch (item.type) {
    case "divider":
      return {
        key: item.key,
        type: "divider",
      };

    case "custom":
      return {
        key: item.key,
        render: item.render,
      };

    case "route":
      return {
        key: item.key,
        label: () =>
          h(
            RouterLink,
            { to: { name: item.route } },
            {
              default: () => item.label,
            },
          ),
        icon: renderMenuIcon(item.icon),
        disabled: item.disabled,
      };

    case "group":
      return {
        key: item.key,
        label: item.label,
        icon: renderMenuIcon(item.icon),
        disabled: item.disabled,
        children: mapChildren(item.children),
      };

    case "action":
      return {
        key: item.key,
        label: item.label,
        icon: renderMenuIcon(item.icon),
        disabled: item.disabled,
      };
  }
};

export const buildMenuOptions = (items: AppMenuItem[]): MenuMixedOption[] => {
  return items
    .map(mapMenuItem)
    .filter((item): item is MenuMixedOption => item !== null);
};
