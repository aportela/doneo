<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';

    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider, NButton, NDrawer, NDrawerContent, NIcon, NFlex, NCard } from 'naive-ui'

    import { useBreakpoints } from '@vueuse/core';
    import { TokenManager } from '../modules/auth/services/tokenManager';
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

    const loadingStore = useLoadingStore();

    const isMobile = breakpoints.smaller('mobile')

    const isCollapsed = ref<boolean>(false);

    const showSearchModal = ref(false)

    const showMobileMenu = ref(false)

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
            <SearchModal v-model:show="showSearchModal" />
            <n-layout>
                <!-- mobile menu on left drawer -->
                <n-drawer v-model:show="showMobileMenu" placement="left" :width="320">
                    <n-drawer-content closable>
                        <template #header>
                            <n-flex align-items="center">
                                <n-icon class="brand-icon" :component="Doneo" :size="isCollapsed ? 16 : 16" />
                                <span class="brand-name">Doneo</span>
                            </n-flex>
                        </template>
                        <NavigationMenu :show-brand-option="false" mode="vertical" />
                    </n-drawer-content>
                </n-drawer>
                <n-layout :has-sider="!isMobile" style="height: 100vh">
                    <!-- desktop menu on left slider -->
                    <n-layout-sider v-if="!isMobile" collapse-mode="width" :collapsed-width="72"
                        :collapsed="isCollapsed" @collapse="isCollapsed = true" @expand="isCollapsed = false" bordered>
                        <NavigationMenu mode="vertical" :collapsed="isCollapsed" @search="showSearchModal = true"
                            show-brand-option />
                    </n-layout-sider>
                    <n-layout>
                        <n-layout-header style="padding: 8px;">
                            <n-card>
                                <n-flex v-if="isMobile" align="center" justify="space-between">
                                    <n-button quaternary @click="showMobileMenu = true">
                                        <template #icon>
                                            <n-icon :component="Menu" />
                                        </template>
                                    </n-button>
                                    <div class="brand">
                                        <n-icon class="brand-icon" :component="Doneo" :size="isCollapsed ? 16 : 16" />
                                        <span class="brand-name">Doneo</span>
                                    </div>
                                </n-flex>
                                <n-flex align-items="center" v-else>
                                    <n-button text @click="isCollapsed = !isCollapsed">
                                        <template #icon>
                                            <n-icon :component="isCollapsed ? PanelLeftOpen : PanelLeftClose" />
                                        </template>
                                    </n-button>
                                    <NavigationBreadcrumb v-if="!isMobile" />
                                </n-flex>
                            </n-card>
                        </n-layout-header>
                        <n-layout-content style="padding: 8px;">
                            <router-view />
                        </n-layout-content>
                    </n-layout>
                </n-layout>
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
</style>