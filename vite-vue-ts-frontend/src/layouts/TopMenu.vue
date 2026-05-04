<script setup lang="ts">
    import { NMenu, NIcon, type MenuOption } from 'naive-ui'
    import { IconBug, IconSitemap, IconFileAnalytics, IconSettings, IconPresentation, IconChartHistogram, IconUsers } from '@tabler/icons-vue';
    import { ref, watch, h } from 'vue'
    import type { Component } from 'vue'
    import { useRouter } from "vue-router";
    import { RouterLink } from 'vue-router'


    const router = useRouter();
    const currentTab = ref<string | number>("home");

    const commonIconSize = 20;

    watch(() => currentTab.value, (newValue) => {
        router.push(
            { name: newValue.toString() }
        ).catch((e) => {
            console.error(e);
        });
    });

    function renderIcon(icon: Component) {
        return (size = commonIconSize) =>
            () =>
                h(NIcon, { size }, {
                    default: () => h(icon)
                })
    }

    function handleUpdateValue(_key: string, _item: MenuOption) {
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
            icon: renderIcon(IconPresentation)(commonIconSize)
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
                    label: () =>
                        h(
                            RouterLink,
                            {
                                to: {
                                    name: 'projectSettings',
                                    params: {
                                    }
                                }
                            },
                            { default: () => 'Project settings' }
                        ),
                    key: 'projectSettings',
                    icon: renderIcon(IconSettings)(commonIconSize)
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

    ]
</script>

<template>
    <div class="header">
        <n-menu :collapsed-width="64" :collapsed-icon-size="commonIconSize" :options="menuOptions" mode="horizontal"
            @update:value="handleUpdateValue" />
    </div>
</template>

<style lang="css" scoped>
    .header {
        height: 64px;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 0 10px;
        box-sizing: border-box;
        /*
        width: 100%;
        */
    }

    .header__container {
        /*
        width: 100%;
        */
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .header__container--contained {
        max-width: 1320px;
        margin: 0 auto;
    }

    .header__container--fluid {
        max-width: 100%;
        margin: 0;
    }
</style>