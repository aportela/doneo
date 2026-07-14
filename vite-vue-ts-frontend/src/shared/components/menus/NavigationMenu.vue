<script setup lang="ts">
    import { h, ref, computed, onMounted } from 'vue';
    import { useRoute, useRouter, RouterLink } from "vue-router";
    import { useI18n } from "vue-i18n";

    import { NMenu, NIcon } from 'naive-ui';
    import type { MenuMixedOption } from "naive-ui/es/menu/src/interface";

    import { useColorSchemeStore } from '../../../stores/colorScheme';
    import { useLoadingStore } from '../../../stores/loading';
    import { useSessionStore } from '../../../stores/session';
    import { useUserSettingsStore } from '../../../stores/userSettings';
    import { useCacheStore } from '../../../stores/cache';

    import { authService } from '../../../modules/auth/services/auth';
    import { navigationMenuService } from './services/navigation-menu';

    import {
        Home,
        FolderKanban,
        ListTodo,
        Users,
        UserKey,
        Settings,
        FolderCog,
        FileCog,
        Bookmark,
        Goal,
        Route,
        CircleUser,
        UserCog,
        LogOut,
        Search,
        Bell,
        BellOff,
        Sun,
        Moon,
        Notebook,
        Folder,
        //FileText,
        PanelLeftOpen,
        PanelTopOpen,
        LibraryBig,
        //SquarePlus,
    } from "@lucide/vue";
    import type { Project } from '../../../modules/projects/models/project';

    interface IProps {
        collapsed?: boolean;
        mode: "horizontal" | "vertical";
    };

    const props = defineProps<IProps>();

    const route = useRoute();

    const router = useRouter();
    const { t } = useI18n();

    const loadingStore = useLoadingStore();
    const sessionStore = useSessionStore();
    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();
    const cacheStore = useCacheStore();

    const currentProjects = ref<Project[]>([]);

    const currentProjectsMenuItems = computed<MenuMixedOption[]>(() =>
        currentProjects.value.map((project: Project) => {
            return (
                {
                    key: project.id || "",
                    label: project.slug,
                    children: [
                        {
                            key: "T" + project.id,
                            label: "Tasks",
                        },
                        {
                            key: "P" + project.id,
                            label: "Pages",
                        }
                    ]
                }
            );
        })
    );

    const menuOptions = computed<MenuMixedOption[]>(() =>
        [
            // home
            {
                key: "home",
                label: () =>
                    h(
                        RouterLink,
                        { to: { name: "home" } },
                        {
                            default: () => t("shared.components.menus.NavigationMenu.items.home"),
                        },
                    ),
                icon: () => h(NIcon, null, { default: () => h(Home) }),
                disabled: false,
                show: true,
            },
            // search
            {
                key: "search",
                label: "Search",
                icon: () => h(NIcon, null, { default: () => h(Search) }),
                disabled: false,
                show: false,
            },
            // workspace
            {
                key: "workspace",
                label: t("shared.components.menus.NavigationMenu.items.workspace"),
                icon: () => h(NIcon, null, { default: () => h(Notebook) }),
                disabled: false,
                show: true,
                children: [
                    {
                        key: "workspaceProjects",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "manageProjects" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.projects"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(FolderKanban) }),
                        disabled: false,
                        show: true,
                    },
                    {
                        key: "workspaceTasks",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "manageTasks" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.tasks"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(ListTodo) }),
                        disabled: false,
                        show: true,
                    },
                    {
                        key: "workspaceWiki",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "manageWiki" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.wiki"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(LibraryBig) }),
                        disabled: false,
                        show: true,
                    }
                ],
            },
            // current projects (not archived)
            {
                key: "currentProjects",
                label: t("shared.components.menus.NavigationMenu.items.currentProjects") + " (" + currentProjects.value.length + ")",
                /*
                label: () =>
                    h(
                        NSpace,
                        {
                            justify: "space-between",
                            align: "center",
                            style: { width: "100%" }
                        },
                        {
                            default: () => [
                                t("shared.components.menus.NavigationMenu.items.currentProjects"),
                                h(
                                    NButton,
                                    {
                                        quaternary: true,
                                        circle: true,
                                        size: "tiny",
                                        onClick: (e: MouseEvent) => {
                                            e.stopPropagation();
                                            console.log("Add");
                                        }
                                    },
                                    {
                                        icon: () => h(SquarePlus)
                                    }
                                )
                            ]
                        }
                    ),
                */
                icon: () => h(NIcon, null, { default: () => h(Folder) }),
                disabled: false,
                show: currentProjects.value.length > 0,
                children: currentProjects.value.length > 0 ? currentProjectsMenuItems.value : undefined,
            },
            // settings
            {
                key: "settings",
                label: t("shared.components.menus.NavigationMenu.items.settings"),
                icon: () => h(NIcon, null, { default: () => h(Settings) }),
                disabled: false,
                show: true,
                children: [
                    // manage users
                    {
                        key: "settingsManageUsers",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "manageUsers" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.manageUsers"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(Users) }),
                        disabled: false,
                        show: true,
                    },
                    // manage roles
                    {
                        key: "settingsManageRoles",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "manageRoles" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.manageRoles"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(UserKey) }),
                        disabled: false,
                        show: true,
                    },
                    // project settings
                    {
                        key: "projectSettings",
                        label: t("shared.components.menus.NavigationMenu.items.projectSettings"),
                        icon: () => h(NIcon, null, { default: () => h(FolderCog) }),
                        disabled: false,
                        show: true,
                        children: [
                            // project types
                            {
                                key: "manageProjectTypes",
                                label: () =>
                                    h(
                                        RouterLink,
                                        { to: { name: "manageProjectTypes" } },
                                        {
                                            default: () => t("shared.components.menus.NavigationMenu.items.manageProjectTypes"),
                                        },
                                    ),
                                icon: () => h(NIcon, null, { default: () => h(Bookmark) }),
                                disabled: false,
                                show: true,
                            },
                            // project priorities
                            {
                                key: "manageProjectPriorities",
                                label: () =>
                                    h(
                                        RouterLink,
                                        { to: { name: "manageProjectPriorities" } },
                                        {
                                            default: () => t("shared.components.menus.NavigationMenu.items.manageProjectPriorities"),
                                        },
                                    ),
                                icon: () => h(NIcon, null, { default: () => h(Goal) }),
                                disabled: false,
                                show: true,
                            },
                            // project statuses
                            {
                                key: "manageProjectStatuses",
                                label: () =>
                                    h(
                                        RouterLink,
                                        { to: { name: "manageProjectStatuses" } },
                                        {
                                            default: () => t("shared.components.menus.NavigationMenu.items.manageProjectStatuses"),
                                        },
                                    ),
                                icon: () => h(NIcon, null, { default: () => h(Route) }),
                                disabled: false,
                                show: true,
                            },
                        ],
                    },
                    // task settings
                    {
                        key: "taskSettings",
                        label: t("shared.components.menus.NavigationMenu.items.taskSettings"),
                        icon: () => h(NIcon, null, { default: () => h(FileCog) }),
                        disabled: false,
                        show: true,
                        children: [
                            // task priorities
                            {
                                key: "manageTaskPriorities",
                                label: () =>
                                    h(
                                        RouterLink,
                                        { to: { name: "manageTaskPriorities" } },
                                        {
                                            default: () => t("shared.components.menus.NavigationMenu.items.manageTaskPriorities"),
                                        },
                                    ),
                                icon: () => h(NIcon, null, { default: () => h(Goal) }),
                                disabled: false,
                                show: true,
                            },
                            // project statuses
                            {
                                key: "manageTaskStatuses",
                                label: () =>
                                    h(
                                        RouterLink,
                                        { to: { name: "manageTaskStatuses" } },
                                        {
                                            default: () => t("shared.components.menus.NavigationMenu.items.manageTaskStatuses"),
                                        },
                                    ),
                                icon: () => h(NIcon, null, { default: () => h(Route) }),
                                disabled: false,
                                show: true,
                            },
                        ],
                    }
                ]
            },
            {
                key: "divider",
                type: "divider",
            },
            // current user
            {
                key: "currentUser",
                label: sessionStore.sessionUserName,
                icon: () => h(NIcon, null, { default: () => h(CircleUser) }),
                disabled: false,
                show: true,
                children: [
                    {
                        key: "switchNavigation",
                        label:
                            t(userSettingsStore.currentNavigationMode === "side" ? "shared.components.menus.NavigationMenu.items.switchTopNavigation" : "shared.components.menus.NavigationMenu.items.switchSidebarNavigation"),
                        icon: () => h(NIcon, null, { default: () => h(userSettingsStore.currentNavigationMode === "side" ? PanelTopOpen : PanelLeftOpen) }),
                        disabled: false,
                        show: true,
                    },
                    {
                        key: "switchNotifications",
                        label:
                            t(userSettingsStore.hasNotificationsEnabled ? "shared.components.menus.NavigationMenu.items.disableNotifications" : "shared.components.menus.NavigationMenu.items.enableNotifications"),
                        icon: () => h(NIcon, null, { default: () => h(userSettingsStore.hasNotificationsEnabled ? BellOff : Bell) }),
                        disabled: false,
                        show: true,
                    },
                    {
                        key: "switchTheme",
                        label:
                            t(colorSchemeStore.light ? "shared.components.menus.NavigationMenu.items.switchToDarkTheme" : "shared.components.menus.NavigationMenu.items.switchToLightTheme"),
                        icon: () => h(NIcon, null, { default: () => h(colorSchemeStore.light ? Moon : Sun) }),
                        disabled: false,
                        show: true,
                    },
                    // profile
                    {
                        key: "profile",
                        label: () =>
                            h(
                                RouterLink,
                                { to: { name: "profile" } },
                                {
                                    default: () => t("shared.components.menus.NavigationMenu.items.profile"),
                                },
                            ),
                        icon: () => h(NIcon, null, { default: () => h(UserCog) }),
                        disabled: false,
                        show: true,
                    },
                    {
                        key: "signout",
                        label:
                            t("shared.components.menus.NavigationMenu.items.signOut"),
                        icon: () => h(NIcon, null, { default: () => h(LogOut) }),
                        disabled: false,
                        show: true,
                    },
                ]
            }
        ]
    );

    const getCurrentProjects = async () => {
        try {
            const response = await navigationMenuService.getCurrentProjects();
            currentProjects.value = response.projects;
        } catch (e) {
            console.log(e);
        }
    };

    const onSignOut = () => {
        loadingStore.set(true);
        authService.signOut().then(() => {
            sessionStore.removeAccessToken();
            cacheStore.clearAllCaches();
            router.push(
                { name: "login" }
            ).catch((e) => {
                console.error(e);
            });
        }).catch(() => {
            sessionStore.removeAccessToken();
            router.push(
                { name: "login" }
            ).catch((e) => {
                console.error(e);
            });
        }).finally(() => {
            loadingStore.set(false);
        });
    };

    const handleMenuSelect = (menuOptionKey: string) => {
        switch (menuOptionKey) {
            case "switchNotifications":
            case "disableNotifications":
            case "enableNotifications":
                userSettingsStore.toggleNotifications();
                break;
            case "switchTheme":
            case "switchDarkTheme":
            case "switchLightTheme":
                colorSchemeStore.toggle();
                break;
            case "switchNavigation":
            case "switchTopNavigation":
            case "switchSidebarNavigation":
                userSettingsStore.toggleNavigationMode();
                break;
            case "signout":
                onSignOut();
                break;
        }
    }

    const currentMenuValue = computed<string>(() => route.name as string);

    onMounted(async () => {
        getCurrentProjects();
    });
</script>

<template>
    <n-menu :mode="props.mode" :collapsed-width="64" :icon-size="16" :collapsed-icon-size="24" :options="menuOptions"
        :value="currentMenuValue" accordion :collapsed="props.collapsed" @update:value="handleMenuSelect" />
</template>

<style lang="css" scoped></style>