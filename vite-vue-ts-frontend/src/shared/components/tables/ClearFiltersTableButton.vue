<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, type ButtonSize } from 'naive-ui';
    import { IconFilterOff } from '@tabler/icons-vue';

    interface RefreshAddActionsColumnProps {
        disabled?: boolean;
        buttonSize?: ButtonSize
        iconSize?: number;
    }

    const emit = defineEmits(['clear'])

    const props = withDefaults(defineProps<RefreshAddActionsColumnProps>(), {
        disabled: false,
        size: "small",
        iconSize: 22,
    });

    const { t } = useI18n();

    const onClear = () => {
        emit("clear");
    };

</script>

<template>
    <n-button size="small" block @click="onClear" :disabled="props.disabled">
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