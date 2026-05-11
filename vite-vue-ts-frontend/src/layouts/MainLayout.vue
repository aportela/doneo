<script setup lang="ts">
    import { ref, reactive, onMounted, onBeforeUnmount } from 'vue';
    import { NLayout, NLayoutHeader, NLayoutSider, NLayoutContent, NSpin, NDialogProvider, NButton, NDrawer, NDrawerContent, NModal } from 'naive-ui'
    import { useBreakpoints } from '@vueuse/core';
    import { useAppBus, type AppBusEvent } from '../shared/composables/bus';
    import { TokenManager } from '../modules/auth/services/tokenManager';
    import { useUserSettingsStore } from '../stores/userSettings';
    import { useLoadingStore } from '../stores/loading';
    import { useSessionStore } from '../stores/session';
    import TopHeader from './TopHeader.vue';
    import TopMenu from './TopMenu.vue';
    import SearchModal from '../shared/components/modals/SearchModal.vue';
    import LoginForm from '../modules/auth/components/LoginForm.vue';
    import SidebarMenu from './SidebarMenu.vue';

    const appBus = useAppBus();

    const reAuthEmitters = reactive<Array<string>>([]);

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

    const accessTokenCheckInterval = 300; // check every 5 min (300 seconds)

    const refreshAccessTokenIfNeeded = async (): Promise<boolean> => {
        if (
            sessionStore.hasAccessToken &&
            sessionStore.accessTokenExpiresBeforeInterval(accessTokenCheckInterval)
        ) {
            return await TokenManager.refreshAccessToken(sessionStore);
        } else {
            return false;
        }
    };


    const onSuccessReauth = () => {
        visibleReauthDialog.value = false;
        appBus.emitReauthValidNotify(reAuthEmitters);
        reAuthEmitters.length = 0;
    };

    let refreshInterval: number;

    onMounted(async () => {
        window.addEventListener('keydown', onGlobalKeydown)
        appBus.on(async (event: AppBusEvent) => {
            if (event.type === "reauthRequired") {
                reAuthEmitters.push(event.emitter);
                try {
                    const success = await TokenManager.refreshAccessToken(sessionStore);
                    if (success) {
                        visibleReauthDialog.value = false;
                        appBus.emitReauthValidNotify(reAuthEmitters);
                        reAuthEmitters.length = 0;
                    } else {
                        visibleReauthDialog.value = true;
                    }
                } catch (error: unknown) {
                    console.error("An unhandled exception occurred during access token refresh", error);
                    visibleReauthDialog.value = true;
                };
            }
        });
        refreshInterval = setInterval(() => {
            refreshAccessTokenIfNeeded()
                .catch((e: Error) => {
                    console.error(
                        "An unhandled exception occurred during access token refresh",
                        e,
                    );
                })
                .finally(() => { });
        }, accessTokenCheckInterval * 1000);
    });

    onBeforeUnmount(() => {
        clearInterval(refreshInterval);
        window.removeEventListener('keydown', onGlobalKeydown)
    });
    const visibleReauthDialog = ref(false);
</script>

<template>
    <n-dialog-provider>

        <n-spin style="height: 100vh;" :show="loadingStore.isLoading">
            <n-modal title="reauth" v-model:show="visibleReauthDialog" preset="card" style="width: 400px;">
                <LoginForm @success="onSuccessReauth" />
            </n-modal>
            <n-layout>
                <n-drawer v-model:show="mobileMenuOpen" placement="left" :width="220">
                    <n-drawer-content closable>
                        <SidebarMenu />
                    </n-drawer-content>
                </n-drawer>
                <SearchModal v-model:show="showSearchModal" />
                <n-layout-header bordered>
                    <TopHeader />
                    <TopMenu v-if="userSettingsStore.topNavigationMode" />
                    <n-button v-if="isMobile" quaternary circle @click="mobileMenuOpen = true">☰</n-button>
                </n-layout-header>
                <n-layout :has-sider="true" v-if="userSettingsStore.sideNavigationMode">
                    <n-layout-sider v-if="!isMobile" collapse-mode="width" :collapsed-width="62" :width="220"
                        :collapsed="isCollapsed" @collapse="isCollapsed = true" @expand="isCollapsed = false"
                        show-trigger="arrow-circle">
                        <SidebarMenu :collapsed="isCollapsed" />
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