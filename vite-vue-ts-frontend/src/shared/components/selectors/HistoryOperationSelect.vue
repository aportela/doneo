<script setup lang="ts">
    import { ref, computed, watch, onMounted, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';


    interface HistoryOperationSelectProps {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        disabled?: boolean;
        showOnlyTaskHistoryOperations?: boolean;
    }

    const { t } = useI18n();

    const selectInstRef = ref<SelectInst | null>(null)

    const selectedValue = ref<number | null>(null);

    const props = defineProps<HistoryOperationSelectProps>();

    const historyOperationType = defineModel<number | null>("historyOperationType", { default: null });

    const options = computed<SelectOption[]>(() => [
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectCreated"),
            value: 100
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectUpdated"),
            value: 101
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectDeleted"),
            value: 102
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectNoteAdded"),
            value: 110
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectNoteUpdated"),
            value: 111
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectNoteDeleted"),
            value: 112
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectAttachmentAdded"),
            value: 120
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectAttachmentDeleted"),
            value: 122
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectPermissionAdded"),
            value: 130
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.projectPermissionDeleted"),
            value: 132
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskCreated"),
            value: 200
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskUpdated"),
            value: 201
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskDeleted"),
            value: 202
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskNoteAdded"),
            value: 210
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskNoteUpdated"),
            value: 211
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskNoteDeleted"),
            value: 212
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskAttachmentAdded"),
            value: 220
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskAttachmentDeleted"),
            value: 222
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskTimeEntryAdded"),
            value: 230
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskTimeEntryUpdated"),
            value: 231
        },
        {
            label: t("shared.components.selectors.HistoryOperationSelect.options.taskTimeEntryDeleted"),
            value: 232
        },

    ]);

    const taskOptions = computed(() => options.value.filter((opt: SelectOption) => typeof opt.value === 'number' && opt.value >= 200));

    watch(selectedValue, (newValue: number | null) => {
        historyOperationType.value = newValue;
    });

    const focus = () => {
        nextTick(() => {
            selectInstRef.value?.focus();
        });
    };

    const reset = () => {
        selectedValue.value = null;
    }

    defineExpose({ reset });

    onMounted(() => {
        if (props.autoFocus) {
            focus();
        }
    });
</script>

<template>
    <n-select filterable ref="selectInstRef" :required="props.required" :clearable="props.clearable"
        v-model:value="selectedValue" :options="props.showOnlyTaskHistoryOperations ? taskOptions : options"
        :placeholder="props.placeholder" :size="props.size" :disabled="props.disabled" />
</template>

<style lang="css" scoped></style>