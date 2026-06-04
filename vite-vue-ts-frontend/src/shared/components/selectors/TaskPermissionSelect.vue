<script setup lang="ts">
    import { ref, computed, watch, onMounted, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';

    import type { TaskPermissionSelectValue } from '../../types/task-permission-select-value';

    interface RoleSelectorProps {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        disabled?: boolean;
    }

    const { t } = useI18n();

    const selectInstRef = ref<SelectInst | null>(null)

    const selectedValue = ref<number | null>(null);

    const props = defineProps<RoleSelectorProps>();

    const permission = defineModel<TaskPermissionSelectValue | null>("permission", { default: null });

    const options = computed<SelectOption[]>(() => [
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.updateTaskAllowed"),
            value: 1
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.updateTaskDenied"),
            value: 2
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.deleteTaskAllowed"),
            value: 3
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.deleteTaskDenied"),
            value: 4
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.viewTaskAllowed"),
            value: 5
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.viewTaskDenied"),
            value: 6
        },
    ]);

    watch(selectedValue, (newValue: number | null) => {
        switch (newValue) {
            case 1:
                permission.value = "updateTaskAllowed"
                break;
            case 2:
                permission.value = "updateTaskDenied"
                break;
            case 3:
                permission.value = "deleteTaskAllowed"
                break;
            case 4:
                permission.value = "deleteTaskDenied"
                break;
            case 5:
                permission.value = "viewTaskAllowed"
                break;
            case 6:
                permission.value = "viewTaskDenied"
                break;
            default:
                permission.value = null;
                break;
        }
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
        v-model:value="selectedValue" :options="options" :placeholder="props.placeholder" :size="props.size"
        :disabled="props.disabled" />
</template>

<style lang="css" scoped></style>