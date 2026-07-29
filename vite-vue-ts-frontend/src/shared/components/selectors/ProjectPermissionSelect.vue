<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, type SelectOption } from 'naive-ui';

    import type { ProjectPermissionSelectValue } from '../../types/project-permission-select-value';

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
    }

    const props = defineProps<Props>();

    const { t } = useI18n();

    const permission = defineModel<ProjectPermissionSelectValue | null>("permission", { default: null });

    const options = computed<SelectOption[]>(() => [
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.updateProjectAllowed"),
            value: "updateProjectAllowed",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.updateProjectDenied"),
            value: "updateProjectDenied",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.deleteProjectAllowed"),
            value: "deleteProjectAllowed",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.deleteProjectDenied"),
            value: "deleteProjectDenied",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.viewProjectAllowed"),
            value: "viewProjectAllowed",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.viewProjectDenied"),
            value: "viewProjectDenied",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.addTaskAllowed"),
            value: "addTaskAllowed",
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.addTaskDenied"),
            value: "addTaskDenied",
        },
    ]);
</script>

<template>
    <n-select filterable :clearable="props.clearable" v-model:value="permission" :options="options"
        :placeholder="props.placeholder" :disabled="props.disabled" />
</template>

<style lang="css" scoped></style>