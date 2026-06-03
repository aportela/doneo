<script setup lang="ts">
    import { NInput, NIcon } from 'naive-ui';
    import type { InputSize } from "naive-ui";
    import { IconSearch } from '@tabler/icons-vue';

    interface Props {
        disabled?: boolean;
        size?: InputSize,
        placeholder?: string;
        clearable?: boolean;
    }

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        size: "medium",
        clearable: false,
    });

    const emit = defineEmits(['keydownEnter']);

    const model = defineModel<string>({
        default: "",
    });

    function onKeyDownEnter() {
        emit('keydownEnter');
    }
</script>

<template>
    <n-input :size="props.size" :disabled="props.disabled" v-model:value="model" :placeholder="props.placeholder"
        :clearable="props.clearable" @keydown.enter="onKeyDownEnter">
        <template #prefix>
            <n-icon>
                <IconSearch />
            </n-icon>
        </template>
    </n-input>
</template>

<style lang="css" scoped></style>