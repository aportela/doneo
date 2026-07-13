<script setup lang="ts">
    import { NIcon } from 'naive-ui'
    import { Search } from '@lucide/vue';

    import { useUserSettingsStore } from '../stores/userSettings';

    import Doneo from '../shared/components/icons/Doneo.vue';
    import BreadCrumb from './BreadCrumb.vue';
    import NavigationMenu from '../shared/components/menus/NavigationMenu.vue';

    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['openSearchModal']);

    const commonIconSize = 18;

    const onSearch = () => {
        emit('openSearchModal')
    };

    const showBreadCrumb = true;
</script>

<template>
    <div class="top-header">
        <div class="top-header__container top-header__container--fluid">

            <div class="brand-container">
                <n-icon :size="commonIconSize" :component="Doneo" />
                <span class="brand-name">Doneo</span>
            </div>

            <div class="header-actions">
                <div class="header-actions__left">
                    <BreadCrumb v-if="showBreadCrumb" />
                </div>

                <div class="header-actions__center">
                    <div @click="onSearch">
                        <span class="shortcut">
                            <n-icon :size="16" :component="Search" />
                            <kbd>Ctrl</kbd>+<kbd>K</kbd> to open search
                        </span>
                    </div>
                </div>

                <div class="header-actions__right">
                    <NavigationMenu mode="horizontal" v-if="userSettingsStore.topNavigationMode" />
                </div>

            </div>

        </div>
    </div>
</template>

<style lang="css" scoped>
    .top-header__container {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
    }

    .brand-container {
        display: flex;
        align-items: center;
        flex-shrink: 0;
        margin-left: 23px;
    }

    .header-actions {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-left: 172px;
        margin-right: 16px;
    }

    .header-actions__left,
    .header-actions__center,
    .header-actions__right {
        display: flex;
        align-items: center;
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