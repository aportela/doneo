import type { VNode } from "vue";
import type { MenuIconName } from "./icons";

export type AppMenuItem =
  | AppMenuRoute
  | AppMenuAction
  | AppMenuGroup
  | AppMenuDivider
  | AppMenuCustom;

export interface AppMenuBase {
  key: string;
  label: string;

  icon?: MenuIconName;

  show?: boolean;
  disabled?: boolean;
}

export interface AppMenuRoute extends AppMenuBase {
  type: "route";
  route: string;
}

export interface AppMenuAction extends AppMenuBase {
  type: "action";
  action: string;
}

export interface AppMenuGroup extends AppMenuBase {
  type: "group";
  children: AppMenuItem[];
}

export interface AppMenuDivider {
  type: "divider";
  key: string;

  show?: boolean;
}

export interface AppMenuCustom {
  type: "custom";
  key: string;

  show?: boolean;

  render: () => VNode;
}
