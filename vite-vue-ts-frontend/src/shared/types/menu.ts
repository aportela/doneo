import { h, ref, computed } from "vue";
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

const menuOptionIconSize = 20;

export { menuOptionIconSize };

// TODO: i18n
export function useMenu() {
  const { t } = useI18n();

  const sessionStore = useSessionStore();

  const lightTheme = ref<boolean>(true);
  const darkTheme = ref<boolean>(false);
  const notificationsDisabled = ref<boolean>(false);
  const notificationsEnabled = ref<boolean>(true);
  const topNavigation = ref<boolean>(true);
  const sideNavigation = ref<boolean>(false);

  const menuOptions = computed(() => {
    return [
      {
        label: () =>
          h(NInput, {
            placeholder: t("search..."),
            clearable: true,
          }),
        key: "search",
        show: false,
        icon: renderIcon(IconSearch)(menuOptionIconSize),
      },
      {
        label: () =>
          h(
            RouterLink,
            {
              to: {
                name: "home",
                params: {},
              },
            },
            { default: () => t("layouts.sidebarMenu.options.home") },
          ),
        key: "home",
        icon: renderIcon(IconPresentation)(menuOptionIconSize),
      },
      {
        label: () =>
          h(
            RouterLink,
            {
              to: {
                name: "manageProjects",
                params: {},
              },
            },
            { default: () => t("layouts.sidebarMenu.options.projects") },
          ),
        key: "projects",
        icon: renderIcon(IconSitemap)(menuOptionIconSize),
      },
      {
        label: () =>
          h(
            RouterLink,
            {
              to: {
                name: "manageTasks",
                params: {},
              },
            },
            { default: () => t("layouts.sidebarMenu.options.tasks") },
          ),
        key: "tasks",
        disabled: false,
        icon: renderIcon(IconBug)(menuOptionIconSize),
      },
      {
        label: t("layouts.sidebarMenu.options.reports"),
        key: "reports",
        disabled: true,
        icon: renderIcon(IconFileAnalytics)(menuOptionIconSize),
      },
      {
        label: t("layouts.sidebarMenu.options.charts"),
        key: "charts",
        disabled: true,
        icon: renderIcon(IconChartHistogram)(menuOptionIconSize),
      },
      {
        key: "divider-2",
        type: "divider",
        show: false,
        props: {
          style: {
            marginLeft: "32px",
          },
        },
      },
      {
        label: t("layouts.sidebarMenu.options.settings"),
        key: "settings",
        show: sessionStore.sessionUserIsAdmin,
        icon: renderIcon(IconSettings)(menuOptionIconSize),
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
            icon: renderIcon(IconUsers)(menuOptionIconSize),
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
            icon: renderIcon(IconUserCheck)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.projectSettings"),
            key: "projectSettings",
            icon: renderIcon(IconSettings)(menuOptionIconSize),
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
                icon: renderIcon(IconBookmark)(menuOptionIconSize),
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
                icon: renderIcon(IconFlagBolt)(menuOptionIconSize),
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
                icon: renderIcon(IconAdjustmentsBolt)(menuOptionIconSize),
              },
            ],
          },
          {
            label: t("layouts.sidebarMenu.options.taskSettings"),
            key: "taskSettings",
            icon: renderIcon(IconSettings)(menuOptionIconSize),
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
                icon: renderIcon(IconFlagBolt)(menuOptionIconSize),
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
                icon: renderIcon(IconAdjustmentsBolt)(menuOptionIconSize),
              },
            ],
          },
        ],
      },
      {
        key: "divider-3",
        type: "divider",
        props: {
          style: {
            marginLeft: "32px",
          },
        },
      },
      {
        label: sessionStore.sessionUserName,
        key: "myuser",
        icon: renderIcon(IconUserCircle)(menuOptionIconSize),
        children: [
          {
            label: "Side navigation",
            key: "sideNavigation",
            show: topNavigation.value,
            icon: renderIcon(IconLayoutSidebarLeftExpand)(menuOptionIconSize),
          },
          {
            label: "Top navigation",
            key: "topNavigation",
            show: sideNavigation.value,
            icon: renderIcon(IconLayoutNavbarExpand)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.disableNotifications"),
            key: "disableNotifications",
            show: notificationsDisabled.value,
            icon: renderIcon(IconBellOff)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.enableNotifications"),
            key: "enableNotifications",
            show: notificationsEnabled.value,
            icon: renderIcon(IconBell)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.switchToLightTheme"),
            key: "switchToLightTheme",
            show: darkTheme.value,
            icon: renderIcon(IconSun)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.switchToDarkTheme"),
            key: "switchToDarkTheme",
            show: lightTheme.value,
            icon: renderIcon(IconMoon)(menuOptionIconSize),
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
            icon: renderIcon(IconId)(menuOptionIconSize),
          },
          {
            label: t("layouts.sidebarMenu.options.signOut"),
            key: "signout",
            icon: renderIcon(IconLogout)(menuOptionIconSize),
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
