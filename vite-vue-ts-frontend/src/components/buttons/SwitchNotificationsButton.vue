<script setup lang="ts">
    import { useI18n } from "vue-i18n";
    import { NButton, NTooltip } from 'naive-ui';
    import { IconBell, IconBellOff } from '@tabler/icons-vue';
    import { useUserSettingsStore } from '../../stores/userSettings';

    const { t } = useI18n();

    interface SwitchNotificationsButtonProps {
        size?: number,
    };

    withDefaults(defineProps<SwitchNotificationsButtonProps>(), {
        size: 20
    });

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary hint>
                <IconBellOff :size="size" v-if="userSettingsStore.disableNotifications"
                    @click.prevent="userSettingsStore.toggleNotifications" @mousedown.prevent />
                <IconBell :size="size" v-else @click="userSettingsStore.toggleNotifications" @mousedown.prevent />
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