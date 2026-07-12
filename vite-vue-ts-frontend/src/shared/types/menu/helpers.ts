import type {
  AppMenuAction,
  AppMenuCustom,
  AppMenuDivider,
  AppMenuGroup,
  AppMenuRoute,
  AppMenuItem,
} from "./types";

import type { MenuIconName } from "./icons";

type MenuBaseOptions = {
  key: string;
  label: string;
  icon?: MenuIconName;
  show?: boolean;
  disabled?: boolean;
};

interface MenuRouteOptions extends MenuBaseOptions {
  route: string;
  children?: AppMenuItem[];
}

export function menuRoute(options: MenuRouteOptions): AppMenuRoute {
  return {
    type: "route",
    ...options,
  };
}

interface MenuGroupOptions extends MenuBaseOptions {
  children: AppMenuItem[];
}

export function menuGroup(options: MenuGroupOptions): AppMenuGroup {
  return {
    type: "group",
    ...options,
  };
}

interface MenuActionOptions extends MenuBaseOptions {
  action: string;
}

export function menuAction(options: MenuActionOptions): AppMenuAction {
  return {
    type: "action",
    ...options,
  };
}

interface MenuDividerOptions {
  key: string;
  show?: boolean;
}

export function menuDivider(options: MenuDividerOptions): AppMenuDivider {
  return {
    type: "divider",
    ...options,
  };
}

interface MenuCustomOptions {
  key: string;
  render: AppMenuCustom["render"];
  show?: boolean;
}

export function menuCustom(options: MenuCustomOptions): AppMenuCustom {
  return {
    type: "custom",
    ...options,
  };
}
