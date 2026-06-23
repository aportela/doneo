<script setup lang="ts">
    import { nextTick } from 'vue'
    import { useRoute } from 'vue-router'

    import { NMenu } from 'naive-ui'

    import { menuOptionIconSize, useMenu } from '../shared/types/menu';

    import { useColorSchemeStore } from '../stores/colorScheme';
    import { useUserSettingsStore } from '../stores/userSettings';

    defineProps({
        collapsed: Boolean
    });


    const route = useRoute();

    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

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
</script>

<template>
    <div class="header">
        <n-menu mode="horizontal" :collapsed-width="64" :collapsed-icon-size="menuOptionIconSize" :options="menuOptions"
            :value="route.name as string" accordion :collapsed="collapsed" @update:value="handleMenuSelect" />
    </div>
</template>

<style lang="css" scoped>
    .header {
        height: 64px;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 0 10px;
        box-sizing: border-box;
        width: 100%;
    }

</style>