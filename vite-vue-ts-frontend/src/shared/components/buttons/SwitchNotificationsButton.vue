<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip, type ButtonSize } from 'naive-ui';

    import { Bell, BellOff } from "@lucide/vue";

    import { useUserSettingsStore } from "../../../stores/userSettings";

    import { DEFAULT_SINGLE_ACTION_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface IProps {
        size?: ButtonSize;
        iconSize?: number,
        disabled?: boolean,
    };

    const props = withDefaults(defineProps<IProps>(), {
        size: DEFAULT_SINGLE_ACTION_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button :size="props.size" quaternary @click.prevent="userSettingsStore.toggleNotifications"
                @mousedown.prevent :disabled="props.disabled">
                <n-icon :size="props.iconSize" :component="userSettingsStore.notifications ? BellOff : Bell" />
            </n-button>
        </template>
        {{
            t(userSettingsStore.notifications ?
                "shared.components.buttons.notifications.enable.toolTip" :
                "shared.components.buttons.notifications.disable.toolTip")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>