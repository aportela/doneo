<script setup lang="ts">
    import { ref, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NInputGroup, NInput, NButtonGroup, NButton, NIcon, NTooltip, type InputSize } from 'naive-ui';
    import { IconCheck, IconX } from '@tabler/icons-vue';

    import { DEFAULT_INPUT_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface Props {
        startupEditMode?: boolean;
        showCount?: boolean;
        maxLength?: number;
        disabled?: boolean;
        readOnly?: boolean;
        size?: InputSize;
        iconSize?: number;
        clearable?: boolean;
        placeholder?: string;
        onConfirm?: (newValue: string | null) => void;
        onCancel?: () => void;
    };

    const props = withDefaults(defineProps<Props>(), {
        startupEditMode: false,
        showCount: false,
        disabled: false,
        readOnly: false,
        size: DEFAULT_INPUT_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        clearable: false,
    });

    const { t } = useI18n();

    const value = defineModel<string | null>("value");

    const editValue = ref<string | null>(value.value ?? null);

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

    const confirmNewValue = () => {
        if (typeof props.onConfirm === 'function') {
            props.onConfirm(editValue.value);
        } else {
            editMode.value = !editMode.value;
            value.value = editValue.value;
        }
    };

    const cancelNewValue = () => {
        if (typeof props.onCancel === 'function') {
            props.onCancel();
        } else {
            editMode.value = !editMode.value;
        }
        editValue.value = value.value ?? null;
    };

    const onKeydown = (e: KeyboardEvent) => {
        switch (e.key) {
            case 'Enter':
                confirmNewValue();
                break;
            case 'Escape':
                cancelNewValue();
                break;
        }
    };
</script>

<template>
    <n-input-group>
        <n-input :size="props.size" :readonly="!editMode" v-model:value="editValue"
            :show-count="editMode && props.showCount" :maxlength="props.maxLength" :disabled="props.disabled"
            :clearable="props.clearable" :placeholder="props.placeholder"
            @click="() => { if (!editMode) { toggleMode(); } }" @keydown="onKeydown" />
        <n-button-group v-if="editMode" :size="props.size">
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button :size="props.size" @click="confirmNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :size="props.iconSize" :component="IconCheck" />
                        </template>
                    </n-button>
                </template>
                {{ t("shared.components.inputs.ToggleInput.buttons.confirm.toolTip") }}
            </n-tooltip>
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button :size="props.size" @click="cancelNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :size="props.iconSize" :component="IconX" />
                        </template>
                    </n-button>
                </template>
                {{ t("shared.components.inputs.ToggleInput.buttons.cancel.toolTip") }}
            </n-tooltip>
        </n-button-group>
    </n-input-group>
</template>

<style lang="css" scoped></style>