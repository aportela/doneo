<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconBell, IconBellOff } from '@tabler/icons-vue';

    import { useUserSettingsStore } from "../../../stores/userSettings";

    interface SwitchNotificationsButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchNotificationsButtonProps>(), {
        iconSize: 20
    });

    const { t } = useI18n();

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click.prevent="userSettingsStore.toggleNotifications" @mousedown.prevent>
                <n-icon :size="iconSize" :component="userSettingsStore.disableNotifications ? IconBellOff : IconBell" />
            </n-button>
        </template>
        {{
            t(userSettingsStore.disableNotifications ?
                "Enable notifications" :
                "Disable notifications")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>