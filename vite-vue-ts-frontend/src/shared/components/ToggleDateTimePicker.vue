<script setup lang="ts">
    import { ref, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NDatePicker } from 'naive-ui';

    interface ToggleInputProps {
        startupEditMode?: boolean;
        disabled?: boolean;
        readOnly?: boolean;
        clearable?: boolean;
        placeholder?: string;
        onConfirm?: (newValue: number | null) => void;
        onCancel?: () => void;
    };

    const props = withDefaults(defineProps<ToggleInputProps>(), {
        startupEditMode: false,
        disabled: false,
        readOnly: false,
        clearable: false,
    });

    const { t } = useI18n();

    const value = defineModel<number | null>("value");

    const editValue = ref<number | null>(value.value ?? null);

    watch(value, (newValue) => {
        editValue.value = newValue ?? null;
    });

    const editMode = ref<boolean>(props.startupEditMode);

    const toggleMode = () => {
        if (!props.readOnly) {
            editMode.value = !editMode.value;
        }
    };

    const setEditMode = () => {
        editMode.value = true;
    };

    const setViewMode = () => {
        editMode.value = false;
    };

    defineExpose({ setEditMode, setViewMode });

</script>

<template>
    <n-date-picker :readonly="!editMode" v-model:value="editValue" type="datetime" :disabled="props.disabled"
        :clearable="props.clearable"
        :placeholder="props.placeholder || t('shared.components.pickers.ToggleDateTimePicker.placeholder')"
        @click="() => { if (!editMode) { toggleMode(); } }" />
</template>

<style lang="css" scoped></style>