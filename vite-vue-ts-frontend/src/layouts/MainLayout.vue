<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';

    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider, NButton, NDrawer, NDrawerContent, NIcon, NFlex } from 'naive-ui'

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
    import { Menu, PanelLeftClose, PanelLeftOpen } from '@lucide/vue';
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
                    <n-button v-if="isMobile" quaternary circle @click="mobileMenuOpen = true">☰</n-button>
                </n-layout-header>
                <n-layout :has-sider="true" v-if="userSettingsStore.sideNavigationMode">
                    <n-layout-sider v-if="!isMobile" collapse-mode="width" :collapsed-width="72" :width="220"
                        :collapsed="isCollapsed" @collapse="isCollapsed = true" @expand="isCollapsed = false" bordered>
                        <div class="brand">
                            <n-icon class="brand-icon" :component="Doneo" :size="isCollapsed ? 16 : 16" />
                            <span class="brand-name" v-if="!isCollapsed">Doneo</span>
                        </div>
                        <NavigationMenu mode="vertical" :collapsed="isCollapsed" @search="showSearchModal = true" />
                    </n-layout-sider>
                    <n-layout>
                        <n-layout-header bordered>
                            <n-flex align-items="center">
                                <n-button v-if="isMobile" quaternary @click="mobileMenuOpen = true">
                                    <template #icon>
                                        <n-icon :component="Menu" />
                                    </template>
                                </n-button>
                                <div class="brand" v-if="isMobile">
                                    <n-icon class="brand-icon" :component="Doneo" :size="isCollapsed ? 16 : 16" />
                                    <span class="brand-name">Doneo</span>
                                </div>
                                <n-button quaternary v-else @click="isCollapsed = !isCollapsed">
                                    <template #icon>
                                        <n-icon :component="isCollapsed ? PanelLeftOpen : PanelLeftClose" />
                                    </template>
                                </n-button>
                                <NavigationBreadcrumb v-if="!isMobile" />
                            </n-flex>
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
</style>