<script setup lang="ts">
    import { nextTick, onMounted } from 'vue';
    import { useRoute, useRouter } from "vue-router";

    import { NDivider, NMenu } from 'naive-ui';
    import { IconDatabaseStar } from '@tabler/icons-vue';

    import { useColorSchemeStore } from '../stores/colorScheme';
    import { useLoadingStore } from '../stores/loading';
    import { useSessionStore } from '../stores/session';
    import { useUserSettingsStore } from '../stores/userSettings';
    import { menuOptionIconSize, useMenu } from '../shared/types/menu';
    import { useCacheStore } from '../stores/cache';
    import { authService } from '../modules/auth/services/auth';

    defineProps({
        collapsed: Boolean
    });

    const route = useRoute();
    const router = useRouter();
    const loadingStore = useLoadingStore();
    const sessionStore = useSessionStore();
    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();
    const cacheStore = useCacheStore();

    const showBrand = false;


    const { menuOptions, lightTheme, darkTheme, notificationsDisabled, notificationsEnabled, topNavigation, sideNavigation } = useMenu();

    const onSignOut = () => {
        loadingStore.set(true);
        authService.signOut().then(() => {
            sessionStore.removeAccessToken();
            cacheStore.clearAllCaches();
            router.push(
                { name: "login" }
            ).catch((e) => {
                console.error(e);
            });
        }).catch(() => {
            sessionStore.removeAccessToken();
            router.push(
                { name: "login" }
            ).catch((e) => {
                console.error(e);
            });
        }).finally(() => {
            loadingStore.set(false);
        });
    };

    const handleMenuSelect = (menuOptionKey: string) => {
        switch (menuOptionKey) {
            case "disableNotifications":
            case "enableNotifications":
                userSettingsStore.toggleNotifications();
                nextTick(() => {
                    notificationsEnabled.value = userSettingsStore.hasNotificationsEnabled;
                    notificationsDisabled.value = !userSettingsStore.hasNotificationsEnabled;
                });
                break;
            case "switchToLightTheme":
            case "switchToDarkTheme":
                colorSchemeStore.toggle();
                nextTick(() => {
                    lightTheme.value = colorSchemeStore.light;
                    darkTheme.value = colorSchemeStore.dark;
                });
                break;
            case "topNavigation":
            case "sideNavigation":
                userSettingsStore.toggleNavigationMode();
                nextTick(() => {
                    topNavigation.value = userSettingsStore.topNavigationMode;
                    sideNavigation.value = userSettingsStore.sideNavigationMode;
                });
                break;
            case "signout":
                onSignOut();
                break;
        }
    }

    onMounted(() => {
        notificationsEnabled.value = userSettingsStore.hasNotificationsEnabled;
        notificationsDisabled.value = !userSettingsStore.hasNotificationsEnabled;
        lightTheme.value = colorSchemeStore.light;
        darkTheme.value = colorSchemeStore.dark;
        topNavigation.value = userSettingsStore.topNavigationMode;
        sideNavigation.value = userSettingsStore.sideNavigationMode;
    });

</script>

<template>
    <div class="brand-container" v-if="showBrand">
        <IconDatabaseStar :size="22" class="brand-icon" />
        <span class="brand-name" v-if="!collapsed">Doneo</span>
    </div>
    <n-divider class="brand-divider" v-if="showBrand" />
    <n-menu :collapsed-width="64" :collapsed-icon-size="menuOptionIconSize" :options="menuOptions"
        :value="route.name as string" accordion :collapsed="collapsed" @update:value="handleMenuSelect" />
</template>

<style lang="css" scoped>
    .brand-container {
        padding-left: 32px;
        padding-top: 16px;
        padding-bottom: 16px;
        display: flex;
        align-items: center;
    }

    .brand-name {
        font-size: --n-font-size;
    }

    .brand-icon {
        margin-right: 8px;
    }

    .brand-divider {
        margin: 6px 18px;
    }
</style>