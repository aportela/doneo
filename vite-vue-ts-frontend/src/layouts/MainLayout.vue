<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';
    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider } from 'naive-ui'
    import { default as TopHeader } from './TopHeader.vue';
    import { default as TopMenu } from './TopMenu.vue';
    import { default as SidebarMenu } from './SidebarMenu.vue';
    import { useUserSettingsStore } from '../stores/userSettings';
    import { useLoadingStore } from '../stores/loading';

    import SearchModal from '../components/modals/SearchModal.vue';

    const userSettingsStore = useUserSettingsStore();

    const loadingStore = useLoadingStore();

    const isCollapsed = ref<boolean>(false);

    const showSearchModal = ref(false)

    function onGlobalKeydown(e: KeyboardEvent) {
        if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
            e.preventDefault()
            showSearchModal.value = true
        }
    }

    onMounted(() => {
        window.addEventListener('keydown', onGlobalKeydown)
    })

    onBeforeUnmount(() => {
        window.removeEventListener('keydown', onGlobalKeydown)
    })
</script>

<template>
    <n-dialog-provider>
        <n-spin style="height: 100vh;" :show="loadingStore.isLoading">
            <n-layout>
                <SearchModal v-model:show="showSearchModal" />
                <n-layout-header bordered>
                    <TopHeader />
                    <TopMenu v-if="userSettingsStore.topNavigationMode" />
                </n-layout-header>
                <n-layout :has-sider="true" v-if="userSettingsStore.sideNavigationMode">
                    <n-layout-sider collapse-mode="width" :collapsed-width="62" :width="220" :collapsed="isCollapsed"
                        @collapse="isCollapsed = true" @expand="isCollapsed = false" show-trigger>
                        <SidebarMenu />
                    </n-layout-sider>
                    <n-layout-content>
                        <router-view />
                    </n-layout-content>
                </n-layout>
                <n-layout-content v-else>
                    <router-view />
                </n-layout-content>
            </n-layout>
        </n-spin>
    </n-dialog-provider>
</template>

<style lang="css" scoped>
    .n-layout-content {
        height: calc(100vh - 64px);
        overflow: auto;
        padding: 16px;
    }
</style>