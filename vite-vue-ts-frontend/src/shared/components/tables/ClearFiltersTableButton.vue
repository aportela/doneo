<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, type ButtonSize } from 'naive-ui';
    import { IconFilterOff } from '@tabler/icons-vue';

    import { DEFAULT_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from "../../../constants";

    interface IProps {
        size?: ButtonSize;
        iconSize?: number;
        disabled?: boolean;
    }

    const props = withDefaults(defineProps<IProps>(), {
        size: DEFAULT_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
    });

    const emit = defineEmits(['clear'])

    const { t } = useI18n();

    const onClear = () => {
        emit("clear");
    };

</script>

<template>
    <n-button :size="props.size" block @click="onClear" :disabled="props.disabled">
        <template #icon>
            <n-icon :size="props.iconSize" :component="IconFilterOff" />
        </template>
        {{ t("shared.components.table.filters.button.clearFilters.label") }}
    </n-button>
</template>

<style lang="css" scoped>
    .doneo-table-actions-button-group {
        width: 100%;
    }

    .doneo-table-actions-button {
        flex: 1;
    }
</style>