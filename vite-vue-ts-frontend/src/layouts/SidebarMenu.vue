<script setup lang="ts">
    import { nextTick, onMounted } from 'vue';
    import { NDivider, NMenu } from 'naive-ui';
    import { IconDatabaseStar } from '@tabler/icons-vue';
    import { useRoute } from 'vue-router'
    import { menuOptionIconSize, useMenu } from '../types/menu';

    import { useColorSchemeStore } from '../stores/colorScheme';
    import { useUserSettingsStore } from '../stores/userSettings';

    defineProps({
        collapsed: Boolean
    });

    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

    const route = useRoute();
    const showBrand = false;


    const { menuOptions, lightTheme, darkTheme, notificationsDisabled, notificationsEnabled } = useMenu();

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
        }
    }

    onMounted(() => {
        notificationsEnabled.value = userSettingsStore.hasNotificationsEnabled;
        notificationsDisabled.value = !userSettingsStore.hasNotificationsEnabled;
        lightTheme.value = colorSchemeStore.light;
        darkTheme.value = colorSchemeStore.dark;
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