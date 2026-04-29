<script setup lang="ts">
    import { NLayoutHeader, NTabs, NTab } from 'naive-ui'
    import { IconBug, IconSitemap, IconHome, IconFileAnalytics, IconSettings } from '@tabler/icons-vue';
    import { ref } from 'vue'

    import { useUserSettingsStore } from '../stores/userSettings';

    const userSettingsStore = useUserSettingsStore();

    const currentTab = ref<string | number>("home");
</script>

<template>
    <n-layout-header bordered class="header">
        <div class="header__container"
            :class="`header__container--${userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'}`">
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
        justify-content: center;
        align-items: center;
        padding: 0 10px;
        box-sizing: border-box;
        width: 100%;
    }

    .header__container {
        width: 100%;
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