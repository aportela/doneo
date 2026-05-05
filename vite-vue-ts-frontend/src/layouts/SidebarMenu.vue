<script setup lang="ts">
    import type { MenuOption } from 'naive-ui'
    import { NMenu } from 'naive-ui';
    import type { Component } from 'vue'
    import { IconPresentation, IconUserCircle, IconBug, IconSitemap, IconFileAnalytics, IconSettings, IconUsers, IconChartHistogram, IconBookmark, IconFlagBolt, IconAdjustmentsBolt } from '@tabler/icons-vue';
    // IconListDetails
    import { NIcon } from 'naive-ui'
    import { h } from 'vue'
    import { RouterLink } from 'vue-router'



    const commonIconSize = 20;

    function renderIcon(icon: Component) {
        return (size = commonIconSize) =>
            () =>
                h(NIcon, { size }, {
                    default: () => h(icon)
                })
    }

    const menuOptions: MenuOption[] = [
        /*
        {
            key: 'divider-1',
            type: 'divider',
            props: {
                style: {
                    marginLeft: '32px'
                }
            }
        },
        */
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
            /*
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'tasks',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Tasks' }
                ),
            */
            key: 'tasks',
            icon: renderIcon(IconBug)(commonIconSize)
        },
        {
            label: "Reports",
            /*
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'reports',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Reports' }
                ),
            key: 'reports',
            */
            icon: renderIcon(IconFileAnalytics)(commonIconSize)
        },
        {
            label: "Charts",
            /*
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'charts',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Charts' }
                ),
            */
            key: 'charts',
            icon: renderIcon(IconChartHistogram)(commonIconSize)
        },
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
            label: "Admin",
            /*
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'settings',
                            params: {
                            }
                        }
                    },
                    { default: () => 'Admin' }
                ),
            */
            key: 'settings',
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
                /*
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
                            { default: () => 'Workspaces' }
                        ),
                    key: 'users',
                    icon: renderIcon(IconMatrix)(commonIconSize)
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
                    icon: renderIcon(IconSitemap)(commonIconSize)
                },
                {
                    label: () =>
                        h(
                            RouterLink,
                            {
                                to: {
                                    name: 'tasks',
                                    params: {
                                    }
                                }
                            },
                            { default: () => 'Tasks' }
                        ),
                    key: 'tasks',
                    icon: renderIcon(IconBug)(commonIconSize)
                },
                {
                    label: () =>
                        h(
                            RouterLink,
                            {
                                to: {
                                    name: 'reports',
                                    params: {
                                    }
                                }
                            },
                            { default: () => 'Reports' }
                        ),
                    key: 'reports',
                    icon: renderIcon(IconBug)(commonIconSize)
                },
                */
                {
                    label: "Project settings",
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
                                    { default: () => 'Project types' }
                                ),
                            key: "111",
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
                                    { default: () => 'Project priorities' }
                                ),
                            key: "222",
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
                                    { default: () => 'Project statuses' }
                                ),
                            key: "333",
                            icon: renderIcon(IconAdjustmentsBolt)(commonIconSize)

                        }
                    ]
                },
            ]
        },
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
            label: "My profile",
            /*
            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'profile',
                            params: {
                            }
                        }
                    },
                    { default: () => 'My profile' }
                ),
            */
            key: 'profile',
            icon: renderIcon(IconUserCircle)(commonIconSize)
        },

    ]

    function handleUpdateValue(_key: string, _item: MenuOption) {
    }

    /*
    const currentWorkspace = ref<string>("w1");

    const workspaceOptions = [
        {
            label: 'Workspace 1',
            value: 'w1',
        },
        {
            label: 'Workspace 2',
            value: 'w2',
        },
        {
            label: 'Workspace 3',
            value: 'w3',
        },
        {
            label: 'Workspace 4',
            value: 'w4',
        },
    ]
    */
</script>

<template>
    <n-menu :collapsed-width="64" :collapsed-icon-size="commonIconSize" :options="menuOptions"
        @update:value="handleUpdateValue" />
</template>

<style lang="css" scoped></style>