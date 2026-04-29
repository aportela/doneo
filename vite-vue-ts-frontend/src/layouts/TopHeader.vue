<script setup lang="ts">
    import { NButton, NDropdown } from 'naive-ui'
    import { NIcon, NFlex } from 'naive-ui'
    import { IconUserCircle, IconDatabaseStar, IconId, IconLogout } from '@tabler/icons-vue';
    import type { Component } from 'vue'
    import { h } from 'vue'
    import { default as SwitchFluidLayoutButton } from '../components/buttons/SwitchFluidLayoutButton.vue';
    import { default as GithubButton } from '../components/buttons/GithubButton.vue';
    import { default as SwitchColorSchemeButton } from '../components/buttons/SwitchColorSchemeButton.vue';
    import { useUserSettingsStore } from '../stores/userSettings';

    const userSettingsStore = useUserSettingsStore();

    const commonIconSize = 18;

    function renderIcon(icon: Component) {
        return (size = commonIconSize) =>
            () =>
                h(NIcon, { size }, {
                    default: () => h(icon)
                })
    }

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
</script>

<template>
    <div class="top-header">
        <div class="top-header__container"
            :class="`top-header__container--${userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'}`">
            <div class="brand-container">
                <IconDatabaseStar :size="commonIconSize" />
                <span class="brand-name">Doneo</span>
            </div>
            <n-flex>
                <GithubButton :size="commonIconSize" />
                <SwitchColorSchemeButton :size="commonIconSize" />
                <SwitchFluidLayoutButton :size="commonIconSize" />
                <n-dropdown :options="userDropdownOptions" placement="bottom-end" trigger="hover">
                    <n-button quaternary>
                        <IconUserCircle :size="commonIconSize" />
                        <span class="username">Administrator</span>
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

    .username {
        margin-left: 8px;
    }
</style>