import { h, ref, computed, type Component } from "vue";
import { RouterLink } from "vue-router";
import { NIcon, NInput, type MenuOption } from "naive-ui";
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
} from "@tabler/icons-vue";

const menuOptionIconSize = 20;

const renderIcon = (icon: Component) => {
  return (size = menuOptionIconSize) =>
    () =>
      h(
        NIcon,
        { size },
        {
          default: () => h(icon),
        },
      );
};

export { menuOptionIconSize };

export function useMenu() {
  const lightTheme = ref(true);
  const darkTheme = ref(false);
  const notificationsDisabled = ref(true);
  const notificationsEnabled = ref(true);
  const menuOptions = computed(() => {
    return [
      {
        label: () =>
          h(
            NInput,
            {
              placeholder: "search...",
              clearable: true,
            },
            { default: () => "Overview" },
          ),
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
            { default: () => "Overview" },
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
                name: "projects",
                params: {},
              },
            },
            { default: () => "Projects" },
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
                name: "tasks",
                params: {},
              },
            },
            { default: () => "Tasks" },
          ),
        key: "tasks",
        disabled: false,
        icon: renderIcon(IconBug)(menuOptionIconSize),
      },
      {
        label: "Reports",
        key: "reports",
        disabled: true,
        icon: renderIcon(IconFileAnalytics)(menuOptionIconSize),
      },
      {
        label: "Charts",
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
        label: "Settings",
        key: "settings",
        show: true,
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
                { default: () => "Users" },
              ),
            key: "manageUsers",
            icon: renderIcon(IconUsers)(menuOptionIconSize),
          },
          {
            label: "Roles",
            key: "roles",
            disabled: true,
            icon: renderIcon(IconUserCheck)(menuOptionIconSize),
          },
          {
            label: "Project",
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
                    { default: () => "Type" },
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
                    { default: () => "Priority" },
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
                    { default: () => "Status" },
                  ),
                key: "manageProjectStatuses",
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
        label: "John Doe",
        key: "myuser",
        icon: renderIcon(IconUserCircle)(menuOptionIconSize),
        children: [
          {
            label: "Disable notifications",
            key: "disableNotifications",
            show: notificationsDisabled.value,
            icon: renderIcon(IconBellOff)(menuOptionIconSize),
          },
          {
            label: "Enable notifications",
            key: "enableNotifications",
            show: notificationsEnabled.value,
            icon: renderIcon(IconBell)(menuOptionIconSize),
          },
          {
            label: "Light theme",
            key: "switchToLightTheme",
            show: darkTheme.value,
            icon: renderIcon(IconSun)(menuOptionIconSize),
          },
          {
            label: "Dark theme",
            key: "switchToDarkTheme",
            show: lightTheme.value,
            icon: renderIcon(IconMoon)(menuOptionIconSize),
          },
          {
            label: "Profile",
            key: "profile",
            disabled: true,
            icon: renderIcon(IconId)(menuOptionIconSize),
          },
          {
            label: "Logout",
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
  };
}
