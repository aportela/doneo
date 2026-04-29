<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent } from 'naive-ui'
    import { default as TopHeader } from './TopHeader.vue';
    import { default as TopMenu } from './TopMenu.vue';
    import { default as SidebarMenu } from './SidebarMenu.vue';
    import { useUserSettingsStore } from '../stores/userSettings';

    const userSettingsStore = useUserSettingsStore();

    const withSidebar = computed(() => userSettingsStore.hasFluidLayout);
    const isCollapsed = ref<boolean>(false);
</script>

<template>
    <div :class="`main-container-${userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'}`">
        <n-layout>
            <n-layout-header bordered>
                <TopHeader />
                <TopMenu v-if="!withSidebar" />
            </n-layout-header>
            <n-layout :has-sider="true" v-if="withSidebar">
                <n-layout-sider :collapsed-width="64" :width="240" v-model:collapsed="isCollapsed" show-trigger>
                    <SidebarMenu />
                </n-layout-sider>
                <n-layout-content>
                    <router-view />
                </n-layout-content>
            </n-layout>
            <n-layout-content>
                <router-view />
            </n-layout-content>
        </n-layout>
    </div>
</template>

<style lang="css" scoped>

    .main-container-contained {
        width: 60%;
        margin: 16px auto;
        border: 1px solid rgb(239, 239, 245);
        border-radius: 4px;
    }

    .main-container-fluid {
        width: 100%;
    }

    .n-layout-content {
        height: calc(100vh - 64px);
        overflow: auto;
    }
</style>