<script setup lang="ts">
    import { NInput, NIcon } from 'naive-ui';
    import type { InputSize } from "naive-ui";
    import { IconSearch } from '@tabler/icons-vue';

    import { DEFAULT_INPUT_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface Props {
        disabled?: boolean;
        readOnly?: boolean;
        size?: InputSize,
        iconSize?: number;
        placeholder?: string;
        clearable?: boolean;
    };

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        readOnly: false,
        size: DEFAULT_INPUT_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        clearable: false,
    });

    const emit = defineEmits(['keydownEnter']);

    const model = defineModel<string>({
        default: "",
    });

    function onKeyDownEnter() {
        emit('keydownEnter');
    };
</script>

<template>
    <n-input :size="props.size" :disabled="props.disabled" v-model:value="model" :placeholder="props.placeholder"
        :clearable="props.clearable" @keydown.enter="onKeyDownEnter">
        <template #prefix>
            <n-icon :size="props.iconSize" :component="IconSearch" />
        </template>
    </n-input>
</template>

<style lang="css" scoped></style>