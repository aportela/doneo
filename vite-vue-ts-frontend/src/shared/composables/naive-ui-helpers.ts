import { h, type Component, type VNodeChild } from "vue";
import { useI18n } from "vue-i18n";

import { NIcon, NTag, NTooltip } from "naive-ui";

import { hexToRgba, oklchToHex } from "./color";
import type { IconProps } from "@tabler/icons-vue";
import type { Role } from "../../modules/roles/models/role";
import {
  DONEO_ICON_ACTION_ADD,
  DONEO_ICON_ACTION_DELETE,
  DONEO_ICON_ACTION_EDIT,
  DONEO_ICON_ACTION_SHOW,
} from "../types/icons";

type TFunction = ReturnType<typeof useI18n>["t"];

export const getNaiveUITagColorProperty = (base: string) => {
  return {
    color: hexToRgba(base, 0.2),
    textColor: hexToRgba(base, 1),
    borderColor: hexToRgba(base, 0.5),
  };
};

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

export const renderLabel = (
  value: string | number,
  className?: string,
): VNodeChild =>
  className ? h("span", { class: className }, { default: () => value }) : value;

export const renderIcon =
  (icon: Component, props: IconProps = {}) =>
  (): VNodeChild =>
    h(NIcon, props, {
      default: () => h(icon),
    });

type PermissionIcon = {
  allowed: boolean;
  icon: Component;
  allowedKey: string;
  deniedKey: string;
};

// return rendered permissions array
const renderRolePermissionIcons = (
  permissions: PermissionIcon[],
  t: TFunction,
): VNodeChild =>
  h(
    "div",
    { class: "doneo-flex doneo-gap-2" },
    permissions.map((permission) =>
      h(
        NTooltip,
        { trigger: "hover" },
        {
          trigger: () =>
            h(NIcon, {
              size: 20,
              component: permission.icon,
              class: [
                "doneo-cursor-help",
                { "doneo-disabled-icon": !permission.allowed },
              ],
            }),
          default: () =>
            t(
              permission.allowed ? permission.allowedKey : permission.deniedKey,
            ),
        },
      ),
    ),
  );

// get project permissions array from role row
export const renderProjectPermissionIcons = (
  row: Role,
  t: TFunction,
): VNodeChild => {
  return renderRolePermissionIcons(
    [
      {
        allowed: row.permissions.allowUpdateProject,
        icon: DONEO_ICON_ACTION_EDIT,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectDenied",
      },
      {
        allowed: row.permissions.allowDeleteProject,
        icon: DONEO_ICON_ACTION_DELETE,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectDenied",
      },
      {
        allowed: row.permissions.allowViewProject,
        icon: DONEO_ICON_ACTION_SHOW,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectDenied",
      },
      {
        allowed: row.permissions.allowAddTask,
        icon: DONEO_ICON_ACTION_ADD,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskDenied",
      },
    ],
    t,
  );
};

// get task permissions array from role row
export const renderTaskPermissionIcons = (
  row: Role,
  t: TFunction,
): VNodeChild => {
  return renderRolePermissionIcons(
    [
      {
        allowed: row.permissions.allowUpdateTask,
        icon: DONEO_ICON_ACTION_EDIT,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskDenied",
      },
      {
        allowed: row.permissions.allowDeleteTask,
        icon: DONEO_ICON_ACTION_DELETE,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskDenied",
      },
      {
        allowed: row.permissions.allowViewTask,
        icon: DONEO_ICON_ACTION_SHOW,
        allowedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed",
        deniedKey:
          "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskDenied",
      },
    ],
    t,
  );
};

const COLOR_PICKER_SWATCHES_COUNT = 32;
const COLOR_PICKER_SWATCHES_START_HUE = 15;

export const ColorPickerSwatches = Array.from(
  { length: COLOR_PICKER_SWATCHES_COUNT },
  (_, i) =>
    oklchToHex(
      0.72, // luminosity
      0.14, // color intensity
      (COLOR_PICKER_SWATCHES_START_HUE +
        (i * 360) / COLOR_PICKER_SWATCHES_COUNT) %
        360,
    ),
);
