<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';

    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider, NButton, NDrawer, NDrawerContent, NIcon } from 'naive-ui'

    import { useBreakpoints } from '@vueuse/core';
    import { TokenManager } from '../modules/auth/services/tokenManager';
    import { useUserSettingsStore } from '../stores/userSettings';
    import { useLoadingStore } from '../stores/loading';
    import { useSessionStore } from '../stores/session';
    import SearchModal from '../shared/components/modals/SearchModal.vue';
    import ReAuthModal from '../modules/auth/components/ReAuthModal.vue';
    import RemoteAPIAlertModal from '../shared/components/modals/RemoteAPIAlertModal.vue';
    import NavigationMenu from '../shared/components/menus/NavigationMenu.vue';
    import Doneo from '../shared/components/icons/Doneo.vue';
    import { Search } from '@lucide/vue';
    import NavigationBreadcrumb from '../shared/components/breadcrumbs/NavigationBreadcrumb.vue';

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
                        <div class="brand">
                            <n-icon class="brand-icon" :component="Doneo" :size="isCollapsed ? 24 : 16" />
                            <span class="brand-name">Doneo</span>
                        </div>
                        <div style="margin: 8px 0px;">
                            <span class="shortcut" style="margin-left: 24px;" @click="showSearchModal = true;">
                                <n-icon :size="16" :component="Search" />
                                Search... <kbd>Ctrl</kbd>+<kbd>K</kbd>
                            </span>
                        </div>
                        <NavigationMenu mode="vertical" :collapsed="isCollapsed" />
                    </n-layout-sider>
                    <n-layout>
                        <n-layout-header bordered>
                            <n-button v-if="isMobile" quaternary circle @click="mobileMenuOpen = true">☰</n-button>
                            <NavigationBreadcrumb />
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

    /*
    .n-layout-content {
        height: calc(100vh - 64px);
        overflow: auto;
        padding: 16px;
    }
*/
    .brand {
        display: flex;
        align-items: center;
        border-bottom: solid 1px var(--n-border-color);
        padding: 7px 0px;
    }

    .brand-icon {
        margin-left: 34px;
    }

    .brand span {
        margin-left: 12px;
    }

    .shortcut {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        padding: 2px 6px;
        border: 1px solid rgba(20, 20, 22, 0.274);
        border-radius: 16px;
        width: 172px;
        cursor: pointer;
    }

    kbd {
        padding: 2px 8px;
        border-radius: 4px;
        border: 1px solid #ccc;
        font-family: monospace;
        font-size: 12px;
        color: var(--n-text-color);
    }

</style>