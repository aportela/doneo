<script setup lang="ts">
    import { NLayoutHeader, NButton, NDropdown } from 'naive-ui'
    import { NIcon, NSpace } from 'naive-ui'
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
    <n-layout-header bordered class="top-header">
        <div class="top-header__container"
            :class="`top-header__container--${userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'}`">
            <div class="brand-container">
                <IconDatabaseStar :size="commonIconSize" />
                <span class="brand-name">Doneo</span>
            </div>
            <n-space>
                <GithubButton :size="commonIconSize" />
                <SwitchColorSchemeButton :size="commonIconSize" />
                <SwitchFluidLayoutButton :size="commonIconSize" />
                <n-dropdown :options="userDropdownOptions" placement="bottom-end" trigger="hover">
                    <n-button quaternary>
                        <IconUserCircle :size="commonIconSize" />
                        Administrator
                    </n-button>
                </n-dropdown>
            </n-space>
        </div>
    </n-layout-header>
</template>

<style lang="css" scoped>
    .top-header {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 48px;
        box-sizing: border-box;
        width: 100%;
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
</style>