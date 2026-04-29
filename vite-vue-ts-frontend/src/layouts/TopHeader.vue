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

    function renderIcon(icon: Component) {
        return () => {
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
    ]

</script>

<template>
    <n-layout-header bordered class="header">
        <div class="header-content" :class="userSettingsStore.hasFluidLayout ? 'fluid' : 'contained'">
            <div class="logo">
                <IconDatabaseStar />
                <span class="title">Doneo</span>
            </div>
            <n-space>
                <GithubButton />
                <SwitchColorSchemeButton />
                <SwitchFluidLayoutButton />
                <n-dropdown :options="userDropdownOptions" placement="bottom-end" trigger="hover">
                    <n-button quaternary>
                        <IconUserCircle :size="20" />
                        Administrator
                    </n-button>
                </n-dropdown>
            </n-space>
        </div>
    </n-layout-header>

</template>

<style lang="css" scoped>
    .header {
        height: 64px;
        display: flex;
        align-items: center;
        padding: 0 20px;
    }

    .header-content {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .fluid {
        max-width: 100%;
        margin: 0px auto;
    }

    .contained {
        max-width: 1320px;
        margin: 0px auto;
    }

    .actions {
        display: flex;
        gap: 6px;
        align-items: center;
    }

    .logo {
        display: flex;
        align-items: center;
    }

    .title {
        margin-left: 8px;
        font-size: 18px;
        font-weight: 600;
    }
</style>