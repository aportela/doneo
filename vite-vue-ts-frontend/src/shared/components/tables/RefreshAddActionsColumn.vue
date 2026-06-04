<script setup lang="ts">

    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon, type ButtonSize } from 'naive-ui';
    import { IconRefresh, IconPlus, IconSettings } from '@tabler/icons-vue';

    interface RefreshAddActionsColumnProps {
        disabled?: boolean;
        buttonSize?: ButtonSize
        iconSize?: number;
        hideRefresh?: boolean;
        hideAdd?: boolean;
        hideSettings?: boolean;
    }

    const emit = defineEmits(['refresh', 'add', 'settings'])

    const props = withDefaults(defineProps<RefreshAddActionsColumnProps>(), {
        disabled: false,
        buttonSize: "small",
        iconSize: 22,
        hideRefresh: false,
        hideAdd: false,
        hideSettings: false,
    });

    const { t } = useI18n();

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onSettings = () => {
        emit("settings");
    };
</script>

<template>
    <n-button-group :size="props.buttonSize" class="doneo-table-actions-button-group">
        <n-button @click="onRefresh" :disabled="props.disabled" v-if="!props.hideRefresh"
            class="doneo-table-actions-button">
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconRefresh" />
            </template>
            {{ t("shared.buttons.Refresh.label") }}
        </n-button>
        <n-button @click="onAdd" :disabled="props.disabled" v-if="!props.hideAdd" class="doneo-table-actions-button">
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconPlus" />
            </template>
            {{ t("shared.buttons.Add.label") }}
        </n-button>
        <n-button @click="onSettings" :disabled="props.disabled || true" v-if="!props.hideSettings"
            class="doneo-table-actions-button">
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconSettings" />
            </template>
            {{ t("shared.buttons.Settings.label") }}
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped>
    .doneo-table-actions-button-group {
        width: 100%;
    }

    .doneo-table-actions-button {
        flex: 1;
    }
</style>