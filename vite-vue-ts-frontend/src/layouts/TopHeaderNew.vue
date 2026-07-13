<script setup lang="ts">
    import { NDivider, NIcon } from 'naive-ui'
    import Doneo from '../shared/components/icons/Doneo.vue';
    import { Search } from '@lucide/vue';

    import { useUserSettingsStore } from '../stores/userSettings';

    import BreadCrumb from './BreadCrumb.vue';
    import NavigationMenu from '../shared/components/menus/NavigationMenu.vue';

    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['openSearchModal']);

    const commonIconSize = 18;

    const onSearch = () => {
        emit('openSearchModal')
    };

    const showBreadCrumb = false;
</script>

<template>
    <div class="top-header">
        <div class="top-header__container top-header__container--fluid">
            <div class="brand-container">
                <n-icon :size="commonIconSize" :component="Doneo" />
                <span class="brand-name">Doneo</span>
            </div>
            <n-divider vertical style="margin: 0px 48px;" />
            <BreadCrumb style="margin: 0px 48px;" v-if="showBreadCrumb" />
            <n-divider vertical />
            <div class=" search-container" @click="onSearch">
                <span class="shortcut">
                    <n-icon :size="16" :component="Search" />
                    <kbd>Crtl</kbd>+<kbd>K</kbd> to open search
                </span>
            </div>
            <NavigationMenu mode="horizontal" v-if="userSettingsStore.topNavigationMode" style="margin: 0px 48px;" />
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
        margin-left: 23px;
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