<script setup lang="ts">
    import { ref, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NInputGroup, NSelect, NButtonGroup, NButton, NIcon, NTooltip, type SelectOption, NSpace, NTag } from 'naive-ui';
    import { IconCheck, IconX, IconTag } from '@tabler/icons-vue';

    interface ToggleInputProps {
        startupEditMode?: boolean;
        showCount?: boolean;
        maxLength?: number;
        disabled?: boolean;
        readOnly?: boolean;
        clearable?: boolean;
        placeholder?: string;
        onConfirm?: (newValue: string[]) => void;
        onCancel?: () => void;
    };

    const props = withDefaults(defineProps<ToggleInputProps>(), {
        startupEditMode: false,
        showCount: false,
        disabled: false,
        readOnly: false,
        clearable: false,
    });

    const { t } = useI18n();

    const value = defineModel<string[] | null>("value", { default: [] });

    const editValue = ref<string[]>(value.value ?? []);

    watch(value, (newValue) => {
        editValue.value = newValue ?? [];
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
        editValue.value = value.value ?? [];
    };

    const onAddTag = (label: string): SelectOption => {
        const normalized = label.trim().toLowerCase();
        return {
            label: normalized,
            value: normalized
        };
    }

</script>

<template>
    <n-input-group>
        <n-space v-if="!editMode" class="doneo-tag-selector-container" @click="toggleMode">
            <!-- TODO: router-link filter by tag -->
            <n-tag v-for="tag in editValue" :key="tag" class="doneo-cursor-pointer">
                {{ tag }}
                <template #icon>
                    <n-icon :component="IconTag" />
                </template>
            </n-tag>
        </n-space>
        <n-select v-else v-model:value="editValue" filterable multiple tag :show-arrow="false" :show="false"
            :on-create="onAddTag" :placeholder="props.placeholder"
            @click="() => { if (!editMode) { toggleMode(); } }" />
        <n-button-group v-if="editMode">
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button @click="confirmNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :component="IconCheck" />
                        </template>
                    </n-button>
                </template>
                {{ t("shared.components.selectors.ToggleTagSelector.buttons.confirm.toolTip") }}
            </n-tooltip>
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button @click="cancelNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :component="IconX" />
                        </template>
                    </n-button>
                </template>
                {{ t("shared.components.selectors.ToggleTagSelector.buttons.cancel.toolTip") }}
            </n-tooltip>
        </n-button-group>
    </n-input-group>
</template>

<style lang="css" scoped>
    .doneo-tag-selector-container {
        width: 100%;
        border: 1px solid #e0e0e6;
        border-radius: var(--n-border-radius);
        padding: 6px;
    }
</style>