import { h, ref, computed, type Component } from "vue";
import { useI18n } from "vue-i18n";
import { RouterLink } from "vue-router";

import { NInput, type MenuOption } from "naive-ui";
import { renderIcon } from "../composables/naive-ui-icon";

import { useSessionStore } from "../../stores/session";

import {
  IconPresentation,
  IconUserCircle,
  IconBug,
  IconSitemap,
  IconFileAnalytics,
  IconUserCheck,
  IconSettings,
  IconUsers,
  IconChartHistogram,
  IconBookmark,
  IconFlagBolt,
  IconAdjustmentsBolt,
  IconLogout,
  IconId,
  IconSearch,
  IconBell,
  IconBellOff,
  IconMoon,
  IconSun,
  IconLayoutSidebarLeftExpand,
  IconLayoutNavbarExpand,
} from "@tabler/icons-vue";

interface AppMenuItem {
  key: string;
  label: string;
  icon?: Component;
  route?: string;
  show?: boolean;
  disabled?: boolean;
  children?: AppMenuItem[];
}

const menuOptionIconSize = 20;

const menuIcons = {
  home: renderIcon(IconPresentation)(menuOptionIconSize),
  projects: renderIcon(IconSitemap)(menuOptionIconSize),
  tasks: renderIcon(IconBug)(menuOptionIconSize),

  reports: renderIcon(IconFileAnalytics)(menuOptionIconSize),
  charts: renderIcon(IconChartHistogram)(menuOptionIconSize),

  users: renderIcon(IconUsers)(menuOptionIconSize),
  roles: renderIcon(IconUserCheck)(menuOptionIconSize),

  settings: renderIcon(IconSettings)(menuOptionIconSize),

  projectTypes: renderIcon(IconBookmark)(menuOptionIconSize),
  priorities: renderIcon(IconFlagBolt)(menuOptionIconSize),
  statuses: renderIcon(IconAdjustmentsBolt)(menuOptionIconSize),

  currentUser: renderIcon(IconUserCircle)(menuOptionIconSize),
  profile: renderIcon(IconId)(menuOptionIconSize),
  logout: renderIcon(IconLogout)(menuOptionIconSize),

  search: renderIcon(IconSearch)(menuOptionIconSize),
  notifications: renderIcon(IconBell)(menuOptionIconSize),
  notificationsOff: renderIcon(IconBellOff)(menuOptionIconSize),

  lightTheme: renderIcon(IconSun)(menuOptionIconSize),
  darkTheme: renderIcon(IconMoon)(menuOptionIconSize),

  sideNavigation: renderIcon(IconLayoutSidebarLeftExpand)(menuOptionIconSize),

  topNavigation: renderIcon(IconLayoutNavbarExpand)(menuOptionIconSize),
};

export { menuOptionIconSize };

// TODO: i18n
export function useMenu() {
  interface MenuRouteOptions {
    key: string;
    route: string;
    label: string;
    icon?: Component;
    disabled?: boolean;
    show?: boolean;
  }

  function menuRoute({
    key,
    route,
    label,
    icon,
    disabled,
    show,
  }: MenuRouteOptions): MenuOption {
    return {
      key,
      disabled,
      show,
      label: () =>
        h(
          RouterLink,
          {
            to: { name: route },
          },
          {
            default: () => t(label),
          },
        ),
      icon: icon ? renderIcon(icon)(menuOptionIconSize) : undefined,
    };
  }

  function menuGroup(
    key: string,
    label: string,
    children: MenuOption[],
    options?: {
      icon?: Component;
      show?: boolean;
      disabled?: boolean;
    },
  ): MenuOption {
    return {
      key,
      label,
      children,
      show: options?.show,
      disabled: options?.disabled,
      icon: options?.icon
        ? renderIcon(options.icon)(menuOptionIconSize)
        : undefined,
    };
  }

  function menuDivider(key: string, show = true): MenuOption {
    return {
      key,
      type: "divider",
      show,
      props: {
        style: {
          marginLeft: "32px",
        },
      },
    };
  }

  const { t } = useI18n();

  const menuT = (key: string) => t(`layouts.sidebarMenu.options.${key}`);

  const sessionStore = useSessionStore();

  const lightTheme = ref<boolean>(true);
  const darkTheme = ref<boolean>(false);
  const notificationsDisabled = ref<boolean>(false);
  const notificationsEnabled = ref<boolean>(true);
  const topNavigation = ref<boolean>(true);
  const sideNavigation = ref<boolean>(false);

  const menuOptions = computed<MenuOption[]>(() => {
    return [
      {
        label: () =>
          h(NInput, {
            placeholder: t("search..."),
            clearable: true,
          }),
        key: "search",
        show: false,
        icon: menuIcons.search,
      },
      menuRoute({
        key: "home",
        route: "home",
        label: menuT("home"),
        icon: IconPresentation,
      }),
      menuGroup(
        "workspace",
        "Workspace",
        [
          menuRoute({
            key: "workspaceProjects",
            route: "manageProjects",
            label: menuT("projects"),
            icon: IconSitemap,
          }),
          menuRoute({
            key: "workspaceTasks",
            route: "manageTasks",
            label: menuT("tasks"),
            icon: IconBug,
          }),
        ],
        { show: true },
      ),
      {
        label: t("layouts.sidebarMenu.options.reports"),
        key: "reports",
        disabled: true,
        show: false,
        icon: menuIcons.reports,
      },
      {
        label: t("layouts.sidebarMenu.options.charts"),
        key: "charts",
        disabled: true,
        show: false,
        icon: menuIcons.charts,
      },
      menuDivider("divider-2", false),
      {
        label: "Projects",
        show: true,
        children: [
          {
            label: "Project 1",
            show: true,
            children: [
              {
                label: "Tasks",
                show: true,
              },
              {
                label: "Pages",
                show: true,
              },
            ],
          },
        ],
      },
      {
        label: t("layouts.sidebarMenu.options.settings"),
        key: "settings",
        show: sessionStore.sessionUserIsAdmin,
        children: [
          {
            label: () =>
              h(
                RouterLink,
                {
                  to: {
                    name: "manageUsers",
                    params: {},
                  },
                },
                { default: () => t("layouts.sidebarMenu.options.manageUsers") },
              ),
            key: "manageUsers",
            icon: menuIcons.users,
          },
          {
            label: () =>
              h(
                RouterLink,
                {
                  to: {
                    name: "manageRoles",
                    params: {},
                  },
                },
                { default: () => t("layouts.sidebarMenu.options.manageRoles") },
              ),
            key: "roles",
            icon: menuIcons.roles,
          },
          {
            label: t("layouts.sidebarMenu.options.projectSettings"),
            key: "projectSettings",
            icon: menuIcons.settings,
            children: [
              {
                label: () =>
                  h(
                    RouterLink,
                    {
                      to: {
                        name: "manageProjectTypes",
                        params: {},
                      },
                    },
                    {
                      default: () =>
                        t("layouts.sidebarMenu.options.manageProjectTypes"),
                    },
                  ),
                key: "manageProjectTypes",
                icon: menuIcons.projectTypes,
              },
              {
                label: () =>
                  h(
                    RouterLink,
                    {
                      to: {
                        name: "manageProjectPriorities",
                        params: {},
                      },
                    },
                    {
                      default: () =>
                        t(
                          "layouts.sidebarMenu.options.manageProjectPriorities",
                        ),
                    },
                  ),
                key: "manageProjectPriorities",
                icon: menuIcons.priorities,
              },
              {
                label: () =>
                  h(
                    RouterLink,
                    {
                      to: {
                        name: "manageProjectStatuses",
                        params: {},
                      },
                    },
                    {
                      default: () =>
                        t("layouts.sidebarMenu.options.manageProjectStatuses"),
                    },
                  ),
                key: "manageProjectStatuses",
                icon: menuIcons.statuses,
              },
            ],
          },
          {
            label: t("layouts.sidebarMenu.options.taskSettings"),
            key: "taskSettings",
            icon: menuIcons.settings,
            children: [
              {
                label: () =>
                  h(
                    RouterLink,
                    {
                      to: {
                        name: "manageTaskPriorities",
                        params: {},
                      },
                    },
                    {
                      default: () =>
                        t("layouts.sidebarMenu.options.manageTaskPriorities"),
                    },
                  ),
                key: "manageTaskPriorities",
                icon: menuIcons.priorities,
              },
              {
                label: () =>
                  h(
                    RouterLink,
                    {
                      to: {
                        name: "manageTaskStatuses",
                        params: {},
                      },
                    },
                    {
                      default: () =>
                        t("layouts.sidebarMenu.options.manageTaskStatuses"),
                    },
                  ),
                key: "manageTaskStatuses",
                icon: menuIcons.statuses,
              },
            ],
          },
        ],
      },
      menuDivider("divider-3", true),
      {
        label: sessionStore.sessionUserName,
        key: "myuser",
        icon: menuIcons.currentUser,
        children: [
          {
            label: "Side navigation",
            key: "sideNavigation",
            show: topNavigation.value,
            icon: menuIcons.sideNavigation,
          },
          {
            label: "Top navigation",
            key: "topNavigation",
            show: sideNavigation.value,
            icon: menuIcons.topNavigation,
          },
          {
            label: t("layouts.sidebarMenu.options.disableNotifications"),
            key: "disableNotifications",
            show: notificationsDisabled.value,
            icon: menuIcons.notificationsOff,
          },
          {
            label: t("layouts.sidebarMenu.options.enableNotifications"),
            key: "enableNotifications",
            show: notificationsEnabled.value,
            icon: menuIcons.notifications,
          },
          {
            label: t("layouts.sidebarMenu.options.switchToLightTheme"),
            key: "switchToLightTheme",
            show: darkTheme.value,
            icon: menuIcons.lightTheme,
          },
          {
            label: t("layouts.sidebarMenu.options.switchToDarkTheme"),
            key: "switchToDarkTheme",
            show: lightTheme.value,
            icon: menuIcons.darkTheme,
          },
          {
            label: () =>
              h(
                RouterLink,
                {
                  to: {
                    name: "profile",
                    params: {},
                  },
                },
                { default: () => t("layouts.sidebarMenu.options.profile") },
              ),
            key: "profile",
            icon: menuIcons.profile,
          },
          {
            label: t("layouts.sidebarMenu.options.signOut"),
            key: "signout",
            icon: menuIcons.logout,
          },
        ],
      },
    ] as MenuOption[];
  });

  return {
    menuOptions,
    lightTheme,
    darkTheme,
    notificationsDisabled,
    notificationsEnabled,
    topNavigation,
    sideNavigation,
  };
}
