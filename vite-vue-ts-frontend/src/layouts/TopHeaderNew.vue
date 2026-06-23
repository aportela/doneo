<script setup lang="ts">
    import { nextTick, onMounted, watch } from 'vue';
    import { useRoute, useRouter } from "vue-router";

    import { NButton, NDropdown, NMenu } from 'naive-ui'
    import { NFlex, NInput } from 'naive-ui'
    import { IconUserCircle, IconDatabaseStar, IconId, IconLogout, IconSearch } from '@tabler/icons-vue';

    import SwitchNotificationsButton from '../shared/components/buttons/SwitchNotificationsButton.vue';
    import SwitchNavigationModeButton from '../shared/components/buttons/SwitchNavigationModeButton.vue';
    import GithubButton from '../shared/components/buttons/GithubButton.vue';
    import SwitchColorSchemeButton from '../shared/components/buttons/SwitchColorSchemeButton.vue';
    import SwitchLocaleButton from '../shared/components/buttons/SwitchLocaleButton.vue';
    import { api } from '../shared/composables/api';
    import { useSessionStore } from "../stores/session";
    import { useLoadingStore } from '../stores/loading';
    import { useCacheStore } from "../stores/cache.ts";
    import { renderIcon } from '../shared/composables/naive-ui-icon';
    import TimerButton from "../shared/components/buttons/TimerButton.vue";

    import { menuOptionIconSize, useMenu } from '../shared/types/menu';

    import { useColorSchemeStore } from '../stores/colorScheme';
    import { useUserSettingsStore } from '../stores/userSettings';

    const route = useRoute();
    const router = useRouter();

    const sessionStore = useSessionStore();
    const loadingStore = useLoadingStore();
    const cacheStore = useCacheStore();

    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

    const { menuOptions, lightTheme, darkTheme, notificationsDisabled, notificationsEnabled, topNavigation, sideNavigation } = useMenu();

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
        }
    }

    const commonIconSize = 18;

    const userDropdownOptions = [
        {
            label: 'Profile',
            key: 'profile',
            icon: renderIcon(IconId)(commonIconSize)
        },
        {
            label: 'Logout',
            key: 'logout',
            icon: renderIcon(IconLogout)(commonIconSize)
        }
    ];

    const onUserDropDownSelect = (key: string | number) => {
        switch (key) {
            case "profile":
                break;
            case "logout":
                onSignOut();
                break;
        }
    };

    const onSignOut = () => {
        loadingStore.set(true);
        api.auth.signOut().then(() => {
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

    watch(
        [
            () => userSettingsStore.hasNotificationsEnabled,
            () => colorSchemeStore.light,
            () => colorSchemeStore.dark,
            () => userSettingsStore.topNavigationMode,
            () => userSettingsStore.sideNavigationMode
        ], () => {
            notificationsEnabled.value = userSettingsStore.hasNotificationsEnabled;
            notificationsDisabled.value = !userSettingsStore.hasNotificationsEnabled;
            lightTheme.value = colorSchemeStore.light;
            darkTheme.value = colorSchemeStore.dark;
            topNavigation.value = userSettingsStore.topNavigationMode;
            sideNavigation.value = userSettingsStore.sideNavigationMode;
        });

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
    <div class="top-header">
        <div class="top-header__container top-header__container--fluid">
            <div class="brand-container">
                <IconDatabaseStar :size="commonIconSize" />
                <span class="brand-name">Doneo</span>
            </div>
            <div class="search-container">
                <n-input placeholder="Search..." style="min-width: 50%;" round v-if="false">
                    <template #prefix>
                        <IconSearch :size="16" />
                    </template>
                </n-input>
                <span class="shortcut">
                    <IconSearch :size="16" />
                    <kbd>Crtl</kbd>+<kbd>K</kbd> to open search
                </span>
            </div>
            <n-menu :collapsed-width="64" :collapsed-icon-size="menuOptionIconSize" :options="menuOptions"
                :value="route.name as string" mode="horizontal" @update:value="handleMenuSelect"
                v-if="userSettingsStore.topNavigationMode" />
            <n-flex v-if="false">
                <SwitchLocaleButton :icon-size="commonIconSize" />
                <SwitchNavigationModeButton :icon-size="commonIconSize" />
                <SwitchColorSchemeButton :icon-size="commonIconSize" />
                <SwitchNotificationsButton :icon-size="commonIconSize" />
                <TimerButton :icon-size="commonIconSize" />
                <GithubButton :icon-size="commonIconSize" />
                <n-dropdown v-if="false" :options="userDropdownOptions" placement="bottom-end" trigger="hover"
                    @select="onUserDropDownSelect">
                    <n-button quaternary>
                        <IconUserCircle :size="commonIconSize" />
                        <span class="username">{{ sessionStore.sessionUserName }}</span>
                    </n-button>
                </n-dropdown>
            </n-flex>
        </div>
    </div>
</template>

<style lang="css" scoped>
    .top-header {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 48px;
        padding: 0px 10px;
        box-sizing: border-box;
        width: 100%;
        border-bottom: 1px solid rgb(239, 239, 245)
    }

    .top-header__container {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .top-header__container--contained {
        max-width: 1320px;
        margin: 0 auto;
    }

    .top-header__container--fluid {
        max-width: 100%;
        margin: 0;
    }

    .brand-container {
        display: flex;
        align-items: center;
    }


    .brand-name {
        margin-left: 8px;
        font-size: 18px;
        font-weight: 600;
    }

    .search-container {
        display: flex;
        align-items: center;
    }

    .username {
        margin-left: 8px;
    }


    .shortcut {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;

        padding: 4px 8px;
        border: 1px solid rgb(239, 239, 245);
        border-radius: 17px;
        width: 300px;
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