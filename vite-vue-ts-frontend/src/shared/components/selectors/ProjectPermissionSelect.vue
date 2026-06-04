<script setup lang="ts">
    import { ref, computed, watch, onMounted, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';

    import type { ProjectPermissionSelectValue } from '../../types/project-permission-select-value';

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

    const permission = defineModel<ProjectPermissionSelectValue | null>("permission", { default: null });

    const options = computed<SelectOption[]>(() => [
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.updateProjectAllowed"),
            value: 1
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.updateProjectDenied"),
            value: 2
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.deleteProjectAllowed"),
            value: 3
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.deleteProjectDenied"),
            value: 4
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.viewProjectAllowed"),
            value: 5
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.viewProjectDenied"),
            value: 6
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.addTaskAllowed"),
            value: 7
        },
        {
            label: t("shared.components.selectors.ProjectPermissionSelect.options.addTaskDenied"),
            value: 8
        },
    ]);

    watch(selectedValue, (newValue: number | null) => {
        switch (newValue) {
            case 1:
                permission.value = "updateProjectAllowed"
                break;
            case 2:
                permission.value = "updateProjectDenied"
                break;
            case 3:
                permission.value = "deleteProjectAllowed"
                break;
            case 4:
                permission.value = "deleteProjectDenied"
                break;
            case 5:
                permission.value = "viewProjectAllowed"
                break;
            case 6:
                permission.value = "viewProjectDenied"
                break;
            case 7:
                permission.value = "addTaskAllowed"
                break;
            case 8:
                permission.value = "addTaskDenied"
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