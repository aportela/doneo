<script setup lang="ts">
    import type { MenuOption } from 'naive-ui'
    import { NMenu, NInput } from 'naive-ui';
    import type { Component } from 'vue'
    import { IconPresentation, IconUserCircle, IconBug, IconDatabaseStar, IconSitemap, IconFileAnalytics, IconUserCheck, IconSettings, IconUsers, IconChartHistogram, IconBookmark, IconFlagBolt, IconAdjustmentsBolt, IconLogout, IconId, IconSearch } from '@tabler/icons-vue';
    // IconListDetails
    import { NIcon } from 'naive-ui'
    import { h } from 'vue'
    import { RouterLink } from 'vue-router'

    import { useRoute } from 'vue-router'

    const route = useRoute()


    const commonIconSize = 20;

    function renderIcon(icon: Component) {
        return (size = commonIconSize) =>
            () =>
                h(NIcon, { size }, {
                    default: () => h(icon)
                })
    }

    const menuOptions: MenuOption[] = [
        {
            key: 'divider-1',
            type: 'divider',
            props: {
                style: {
                    marginLeft: '32px'
                }
            }
        },
        {
            label: () =>
                h(
                    NInput,
                    {
                        placeholder: "search...",
                        clearable: true,
                    },
                    { default: () => 'Overview' }
                ),
            key: 'search',
            show: false,
            icon: renderIcon(IconSearch)(commonIconSize),
        },
        {
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'home',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Overview' }
                ),
            key: 'home',
            icon: renderIcon(IconPresentation)(commonIconSize),
        },
        {
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'projects',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Projects' }
                ),
            key: 'projects',
            icon: renderIcon(IconSitemap)(commonIconSize),
        },
        {
            label: "Tasks",
            key: 'tasks',
            disabled: true,
            icon: renderIcon(IconBug)(commonIconSize)
        },
        {
            label: "Reports",
            key: "reports",
            disabled: true,
            icon: renderIcon(IconFileAnalytics)(commonIconSize)
        },
        {
            label: "Charts",
            key: 'charts',
            disabled: true,
            icon: renderIcon(IconChartHistogram)(commonIconSize)
        },
        {
            key: 'divider-2',
            type: 'divider',
            show: false,
            props: {
                style: {
                    marginLeft: '32px'
                }
            }
        },
        {
            label: "Settings",
            key: 'settings',
            show: true,
            icon: renderIcon(IconSettings)(commonIconSize),
            children: [
                {
                    label: () =>
                        h(
                            RouterLink,
                            {
                                to: {
                                    name: 'users',
                                    params: {
                                    }
                                }
                            },
                            { default: () => 'Users' }
                        ),
                    key: 'users',
                    icon: renderIcon(IconUsers)(commonIconSize)
                },
                {
                    label: 'Roles',
                    key: 'roles',
                    disabled: true,
                    icon: renderIcon(IconUserCheck)(commonIconSize),
                },
                {
                    label: "Project",
                    key: 'projectSettings',
                    icon: renderIcon(IconSettings)(commonIconSize),
                    children: [
                        {
                            label: () =>
                                h(
                                    RouterLink,
                                    {
                                        to: {
                                            name: 'manageProjectTypes',
                                            params: {
                                            }
                                        }
                                    },
                                    { default: () => 'Type' }
                                ),
                            key: "projectTypes",
                            icon: renderIcon(IconBookmark)(commonIconSize)
                        },
                        {
                            label: () =>
                                h(
                                    RouterLink,
                                    {
                                        to: {
                                            name: 'manageProjectPriorities',
                                            params: {
                                            }
                                        }
                                    },
                                    { default: () => 'Priority' }
                                ),
                            key: "projectPriorities",
                            icon: renderIcon(IconFlagBolt)(commonIconSize)
                        },
                        {
                            label: () =>
                                h(
                                    RouterLink,
                                    {
                                        to: {
                                            name: 'manageProjectStatuses',
                                            params: {
                                            }
                                        }
                                    },
                                    { default: () => 'Status' }
                                ),
                            key: "projectStatuses",
                            icon: renderIcon(IconAdjustmentsBolt)(commonIconSize)

                        }
                    ]
                },
            ]
        },
        {
            key: 'divider-3',
            type: 'divider',
            props: {
                style: {
                    marginLeft: '32px'
                }
            }
        },
        {
            label: "John Doe",
            key: 'myuser',
            icon: renderIcon(IconUserCircle)(commonIconSize),
            children: [
                {
                    label: "Profile",
                    key: "profile",
                    disabled: true,
                    icon: renderIcon(IconId)(commonIconSize),
                },
                {
                    label: "Logout",
                    key: "signout",
                    icon: renderIcon(IconLogout)(commonIconSize),
                }
            ]
        },
    ];


</script>

<template>
    <div class="brand-container">
        <IconDatabaseStar :size="32" class="brand-icon" />
        <span class="brand-name">Doneo</span>
    </div>
    <n-menu :collapsed-width="64" :collapsed-icon-size="commonIconSize" :options="menuOptions"
        :value="route.name as string" accordion />
</template>

<style lang="css" scoped>
    .brand-container {
        padding-left: 32px;
        display: flex;
        align-items: center;
    }

    .brand-name {
        font-weight: bold;
        font-size: 30px;
    }

    .brand-icon {
        margin-right: 8px;
    }
</style>