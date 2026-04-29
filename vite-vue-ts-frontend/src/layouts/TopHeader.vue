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
        return () => {
            // TODO: size prop
            return h(NIcon, null, {
                default: () => h(icon)
            })
        }
    }

    const userDropdownOptions = [
        {
            label: 'Profile',
            key: 'profile',
            icon: renderIcon(IconId)
        },
        {
            label: 'Logout',
            key: 'logout',
            icon: renderIcon(IconLogout)
        }
    ];

</script>

<template>
    <n-layout-header bordered class="top-header-layout-container">
        <div class="top-header-container"
            :class="userSettingsStore.hasFluidLayout ? 'header-fluid' : 'header-contained'">
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
    .top-header-layout-container {
        height: 64px;
        display: flex;
        align-items: center;
        padding: 0 20px;
    }

    .top-header-container {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .header-fluid {
        max-width: 100%;
    }

    .header-contained {
        max-width: 1320px;
        margin: 0px auto;
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