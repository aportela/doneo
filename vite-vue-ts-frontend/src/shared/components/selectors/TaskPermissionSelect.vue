<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, type SelectOption } from 'naive-ui';

    import type { TaskPermissionSelectValue } from '../../types/task-permission-select-value';

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
    }

    const props = defineProps<Props>();

    const { t } = useI18n();

    const permission = defineModel<TaskPermissionSelectValue | null>("permission", { default: null });

    const options = computed<SelectOption[]>(() => [
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.updateTaskAllowed"),
            value: "updateTaskAllowed",
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.updateTaskDenied"),
            value: "updateTaskDenied",
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.deleteTaskAllowed"),
            value: "deleteTaskAllowed",
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.deleteTaskDenied"),
            value: "deleteTaskDenied",
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.viewTaskAllowed"),
            value: "viewTaskAllowed",
        },
        {
            label: t("shared.components.selectors.TaskPermissionSelect.options.viewTaskDenied"),
            value: "viewTaskDenied",
        },
    ]);

</script>

<template>
    <n-select filterable :clearable="props.clearable" v-model:value="permission" :options="options"
        :placeholder="props.placeholder" :disabled="props.disabled" />
</template>

<style lang="css" scoped></style>