<script setup lang="ts">

    import { NLayoutHeader, NButton, NDropdown } from 'naive-ui'
    import { NIcon } from 'naive-ui'
    import { IconUserCircle, IconBrandGithub, IconDatabaseStar, IconLayout, IconMoon, IconId, IconLogout, IconHome } from '@tabler/icons-vue';
    import type { Component } from 'vue'
    import { ref, h } from 'vue'

    function renderIcon(icon: Component) {
        return () => {
            return h(NIcon, null, {
                default: () => h(icon)
            })
        }
    }


    const options = [
        {
            label: 'User Profile',
            key: 'profile',
            icon: renderIcon(IconId)
        },
        {
            label: 'Logout',
            key: 'logout',
            icon: renderIcon(IconLogout)
        }
    ]

    const headerClassXL = ref(true);
</script>

<template>
    <n-layout-header bordered class="header">
        <div class="header-content" :class="{ 'header-content-xl': headerClassXL }">
            <div class="logo">
                <IconDatabaseStar />
                <span class="title">Doneo</span>
            </div>
            <div class="actions">
                <n-button quaternary>
                    <IconBrandGithub :size="20" />
                </n-button>
                <n-button quaternary>
                    <IconMoon :size="20" />
                </n-button>
                <n-button quaternary @click="headerClassXL = !headerClassXL" @mousedown.prevent>
                    <IconLayout :size="20" />
                </n-button>
                <n-dropdown :options="options" placement="bottom-end" trigger="hover">
                    <n-button quaternary>
                        <IconUserCircle :size="20" />
                        Administrator
                    </n-button>
                </n-dropdown>
            </div>

        </div>
    </n-layout-header>
    <n-layout-header bordered class="header">
        <n-button quaternary>
            <template #icon>
                <NIcon>
                    <IconHome :size="33" />
                </NIcon>
            </template>
            Home
        </n-button>
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

    .header-content-xl {
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