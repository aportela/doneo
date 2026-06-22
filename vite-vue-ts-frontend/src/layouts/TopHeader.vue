<script setup lang="ts">
    import { useRouter } from "vue-router";


    import { NButton, NDropdown } from 'naive-ui'
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


    const router = useRouter();

    const sessionStore = useSessionStore();
    const loadingStore = useLoadingStore();
    const cacheStore = useCacheStore();

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
            <n-flex>

                <SwitchLocaleButton :icon-size="commonIconSize" />
                <GithubButton :icon-size="commonIconSize" />
                <SwitchNotificationsButton :icon-size="commonIconSize" />
                <SwitchColorSchemeButton :icon-size="commonIconSize" />
                <SwitchNavigationModeButton :icon-size="commonIconSize" />
                <TimerButton :icon-size="commonIconSize" />
                <n-dropdown :options="userDropdownOptions" placement="bottom-end" trigger="hover"
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