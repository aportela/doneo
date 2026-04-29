<script setup lang="ts">
    import { NLayoutHeader, NButton, NDropdown, NTabs, NTab } from 'naive-ui'
    import { NIcon } from 'naive-ui'
    import { IconBug, IconSitemap, IconUserCircle, IconBrandGithub, IconDatabaseStar, IconLayout, IconMoon, IconId, IconLogout, IconHome, IconFileAnalytics, IconSettings } from '@tabler/icons-vue';
    import type { Component } from 'vue'
    import { ref, h } from 'vue'

    function renderIcon(icon: Component) {
        return () => {
            return h(NIcon, null, {
                default: () => h(icon)
            })
        }
    }


    const options = [
        {
            label: 'Profile',
            key: 'profile',
            icon: renderIcon(IconId)
        },
        {
            label: 'Logout',
            key: 'logout',
            icon: renderIcon(IconLogout)
        }
    ]

    const headerClassXL = ref(true);

    const currentTab = ref<string | number>("home");


</script>

<template>
    <n-layout-header bordered class="header">
        <div class="header-content" :class="{ 'header-content-xl': headerClassXL }">
            <div class="logo">
                <IconDatabaseStar />
                <span class="title">Doneo</span>
            </div>
            <div class="actions">
                <n-button quaternary>
                    <IconBrandGithub :size="20" />
                </n-button>
                <n-button quaternary>
                    <IconMoon :size="20" />
                </n-button>
                <n-button quaternary @click="headerClassXL = !headerClassXL" @mousedown.prevent>
                    <IconLayout :size="20" />
                </n-button>
                <n-dropdown :options="options" placement="bottom-end" trigger="hover">
                    <n-button quaternary>
                        <IconUserCircle :size="20" />
                        Administrator
                    </n-button>
                </n-dropdown>
            </div>
        </div>
    </n-layout-header>
    <n-layout-header bordered class="header">
        <div :class="{ 'header-content-xl': headerClassXL }">
            <n-tabs type="line" v-model:value="currentTab" animated>
                <n-tab name="home" tab="Home">
                    <RouterLink to="/home" class="nav-link">
                        <IconHome :size="18" />
                        Home
                    </RouterLink>
                </n-tab>
                <n-tab name="projects" tab="Projects">
                    <RouterLink to="/projects" class="nav-link">
                        <IconSitemap :size="18" />
                        Projects
                    </RouterLink>
                </n-tab>
                <n-tab name="tasks" tab="Tasks">
                    <RouterLink to="/tasks" class="nav-link">
                        <IconBug :size="18" />
                        Tasks
                    </RouterLink>
                </n-tab>
                <n-tab name="reports" tab="Reports">
                    <RouterLink to="/reports" class="nav-link">
                        <IconFileAnalytics :size="18" />
                        Reports
                    </RouterLink>
                </n-tab>
                <n-tab name="settings" tab="Settings">
                    <RouterLink to="/settings" class="nav-link">
                        <IconSettings :size="18" />
                        Settings
                    </RouterLink>
                </n-tab>
            </n-tabs>
        </div>
    </n-layout-header>
</template>

<style lang="css" scoped>
    .header {
        height: 64px;
        display: flex;
        align-items: center;
        padding: 0 20px;
    }

    .header-content {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .header-content-xl {
        max-width: 1320px;
        margin: 0px auto;
    }

    .actions {
        display: flex;
        gap: 6px;
        align-items: center;
    }

    .logo {
        display: flex;
        align-items: center;
    }

    .title {
        margin-left: 8px;
        font-size: 18px;
        font-weight: 600;
    }
</style>