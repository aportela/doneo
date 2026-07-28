<script setup lang="ts">
    import { ref, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NInputGroup, NSelect, NButtonGroup, NButton, NIcon, NTooltip, type SelectOption, NFlex, NTag, type SelectSize } from 'naive-ui';
    import { IconCheck, IconX, IconTag } from '@tabler/icons-vue';

    import { DEFAULT_SELECTOR_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface Props {
        startupEditMode?: boolean;
        disabled?: boolean;
        readOnly?: boolean;
        size?: SelectSize;
        iconSize?: number;
        placeholder?: string;
        onConfirm?: (newValue: string[]) => void;
        onCancel?: () => void;
    };

    const props = withDefaults(defineProps<Props>(), {
        startupEditMode: false,
        disabled: false,
        readOnly: false,
        size: DEFAULT_SELECTOR_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
    });

    const { t } = useI18n();

    const value = defineModel<string[]>("value", {
        default: () => []
    });

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
    };
</script>

<template>
    <n-input-group>
        <n-flex v-if="!editMode" class="doneo-tag-selector-container doneo-cursor-pointer" @click="toggleMode">
            <!-- TODO: router-link filter by tag -->
            <n-tag :size="props.size" v-for="tag in editValue" :key="tag" class="doneo-cursor-pointer">
                {{ tag }}
                <template #icon>
                    <n-icon :size="props.iconSize" :component="IconTag" />
                </template>
            </n-tag>
        </n-flex>
        <n-select v-else :size="props.size" v-model:value="editValue" filterable multiple tag :show-arrow="false"
            :show="false" :on-create="onAddTag" :placeholder="props.placeholder"
            @click="() => { if (!editMode) { toggleMode(); } }" />
        <n-button-group :size="props.size" v-if="editMode">
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button :size="props.size" @click="confirmNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :size="props.iconSize" :component="IconCheck" />
                        </template>
                    </n-button>
                </template>
                {{ t("shared.components.selectors.ToggleTagSelector.buttons.confirm.toolTip") }}
            </n-tooltip>
            <n-tooltip trigger="hover">
                <template #trigger>
                    <n-button :size="props.size" @click="cancelNewValue" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :size="props.iconSize" :component="IconX" />
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
        min-height: 20px;
        border: 1px solid #e0e0e6;
        border-radius: var(--n-border-radius);
        padding: 6px;
    }
</style>