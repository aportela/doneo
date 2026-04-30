<script setup lang="ts">
    import type { MenuOption, } from 'naive-ui'
    import { NMenu } from 'naive-ui';
    import type { Component } from 'vue'
    import { IconBug, IconSitemap, IconHome, IconFileAnalytics, IconSettings, IconSquarePlus, IconListDetails, IconMatrix } from '@tabler/icons-vue';
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
                    { default: () => 'Welcome' }
                ),
            key: 'home',
            icon: renderIcon(IconHome)(commonIconSize)
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

            label: () =>
                h(
                    RouterLink,
                    {
                        to: {
                            name: 'workspaces',
                            params: {
                            }
                        }
                    },
                    { default: () => 'WorkSpaces' }
                ),
            key: 'workspaces',
            icon: renderIcon(IconMatrix)(commonIconSize),
            children: [
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
                            { default: () => 'Add' }
                        ),
                    key: 'add_project',
                    icon: renderIcon(IconSquarePlus)(commonIconSize)
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
                            { default: () => 'List' }
                        ),
                    key: 'list_project',
                    icon: renderIcon(IconListDetails)(commonIconSize)
                },
            ]
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
            children: [
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
                            { default: () => 'Add' }
                        ),
                    key: 'add_project',
                    icon: renderIcon(IconSquarePlus)(commonIconSize)
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
                            { default: () => 'List' }
                        ),
                    key: 'list_project',
                    icon: renderIcon(IconListDetails)(commonIconSize)
                },
            ]
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
            icon: renderIcon(IconFileAnalytics)(commonIconSize)
        },
        {
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
                    { default: () => 'Settings' }
                ),
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
                            { default: () => 'Manage users' }
                        ),
                    key: 'users',
                    icon: renderIcon(IconSettings)(commonIconSize)
                },
            ]
        },
    ]


    function handleUpdateValue(_key: string, _item: MenuOption) {
    }
</script>

<template>
    <n-menu :collapsed-width="64" :collapsed-icon-size="commonIconSize" :options="menuOptions"
        @update:value="handleUpdateValue" />
</template>

<style lang="css" scoped></style>