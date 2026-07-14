<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';

    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider, NButton, NDrawer, NDrawerContent } from 'naive-ui'

    import { useBreakpoints } from '@vueuse/core';
    import { TokenManager } from '../modules/auth/services/tokenManager';
    import { useUserSettingsStore } from '../stores/userSettings';
    import { useLoadingStore } from '../stores/loading';
    import { useSessionStore } from '../stores/session';
    import TopHeader from './TopHeader.vue';
    import SearchModal from '../shared/components/modals/SearchModal.vue';
    import ReAuthModal from '../modules/auth/components/ReAuthModal.vue';
    import RemoteAPIAlertModal from '../shared/components/modals/RemoteAPIAlertModal.vue';
    import NavigationMenu from '../shared/components/menus/NavigationMenu.vue';

    const sessionStore = useSessionStore();

    const breakpoints = useBreakpoints({
        mobile: 768
    });

    const userSettingsStore = useUserSettingsStore();

    const loadingStore = useLoadingStore();

    const isMobile = breakpoints.smaller('mobile')

    const isCollapsed = ref<boolean>(false);

    const showSearchModal = ref(false)

    const mobileMenuOpen = ref(false)

    function onGlobalKeydown(e: KeyboardEvent) {
        if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
            e.preventDefault()
            showSearchModal.value = true
        }
    }

    const accessTokenSecondsCheckInterval = 300; // check every 5 min (300 seconds)

    const refreshAccessTokenIfNeeded = async (): Promise<boolean> => {
        if (
            sessionStore.hasAccessToken &&
            sessionStore.accessTokenExpiresBeforeInterval(accessTokenSecondsCheckInterval)
        ) {
            return await TokenManager.refreshAccessToken(sessionStore);
        } else {
            return false;
        }
    };


    let refreshInterval: number;

    onMounted(async () => {
        window.addEventListener('keydown', onGlobalKeydown);

        refreshInterval = setInterval(() => {
            refreshAccessTokenIfNeeded()
                .catch((e: Error) => {
                    console.error(
                        "An unhandled exception occurred during access token refresh",
                        e,
                    );
                })
                .finally(() => { });
        }, accessTokenSecondsCheckInterval * 1000);
    });

    onBeforeUnmount(() => {
        clearInterval(refreshInterval);
        window.removeEventListener('keydown', onGlobalKeydown);
    });
</script>

<template>
    <n-dialog-provider>
        <n-spin style="height: 100vh;" :show="loadingStore.isLoading">
            <ReAuthModal />
            <RemoteAPIAlertModal />
            <n-layout>
                <n-drawer v-model:show="mobileMenuOpen" placement="left" :width="320">
                    <n-drawer-content closable>
                        <NavigationMenu mode="vertical" />
                    </n-drawer-content>
                </n-drawer>
                <SearchModal v-model:show="showSearchModal" />
                <n-layout-header bordered v-if="false">
                    <TopHeader @open-search-modal="showSearchModal = true;" />
                    <n-button v-if="isMobile" quaternary circle @click="mobileMenuOpen = true">☰</n-button>
                </n-layout-header>
                <n-layout :has-sider="true" v-if="userSettingsStore.sideNavigationMode">
                    <n-layout-sider v-if="!isMobile" collapse-mode="width" :collapsed-width="62" :width="220"
                        :collapsed="isCollapsed" @collapse="isCollapsed = true" @expand="isCollapsed = false"
                        show-trigger="arrow-circle" bordered>
                        <NavigationMenu mode="vertical" :collapsed="isCollapsed" />
                    </n-layout-sider>
                    <n-layout>
                        <n-layout-header bordered>
                            <TopHeader @open-search-modal="showSearchModal = true;" />
                            <n-button v-if="isMobile" quaternary circle @click="mobileMenuOpen = true">☰</n-button>
                        </n-layout-header>
                        <n-layout-content>
                            <router-view />
                        </n-layout-content>
                    </n-layout>
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