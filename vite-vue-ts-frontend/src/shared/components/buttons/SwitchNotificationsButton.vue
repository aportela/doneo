<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconBell, IconBellOff } from '@tabler/icons-vue';

    import { useUserSettingsStore } from "../../../stores/userSettings";

    import { BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    interface IProps {
        iconSize?: number,
        disabled?: boolean,
    };

    const props = withDefaults(defineProps<IProps>(), {
        iconSize: BUTTON_DEFAULT_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click.prevent="userSettingsStore.toggleNotifications" @mousedown.prevent
                :disabled="props.disabled">
                <n-icon :size="props.iconSize"
                    :component="userSettingsStore.disableNotifications ? IconBellOff : IconBell" />
            </n-button>
        </template>
        {{
            t(userSettingsStore.disableNotifications ?
                "shared.components.buttons.notifications.enable.toolTip" :
                "shared.components.buttons.notifications.disable.toolTip")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>