<script setup lang="ts">
    import { NTabs, NTab } from 'naive-ui'
    import { IconBug, IconSitemap, IconHome, IconFileAnalytics, IconSettings } from '@tabler/icons-vue';
    import { ref, watch } from 'vue'
    import { useRouter } from "vue-router";
    import { useUserSettingsStore } from '../stores/userSettings';

    const router = useRouter();
    const userSettingsStore = useUserSettingsStore();
    const currentTab = ref<string | number>("home");

    watch(() => currentTab.value, (newValue) => {
        router.push(
            { name: newValue.toString() }
        ).catch((e) => {
            console.error(e);
        });
    });
</script>

<template>
    <div class="header">
        <div class="header__container"
            :class="`header__container--${userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'}`">
            <n-tabs type="bar" v-model:value="currentTab" animated>
                <n-tab name="home" tab="Home">
                    <IconHome :size="18" />
                    Home
                </n-tab>
                <n-tab name="projects" tab="Projects">
                    <IconSitemap :size="18" />
                    Projects
                </n-tab>
                <n-tab name="tasks" tab="Tasks">
                    <IconBug :size="18" />
                    Tasks
                </n-tab>
                <n-tab name="reports" tab="Reports">
                    <IconFileAnalytics :size="18" />
                    Reports
                </n-tab>
                <n-tab name="settings" tab="Settings">
                    <IconSettings :size="18" />
                    Settings
                </n-tab>
            </n-tabs>
        </div>
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