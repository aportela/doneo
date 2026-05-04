<script setup lang="ts">
    import { h, type CSSProperties, type VNodeChild, type Component } from 'vue';
    import { NSelect, NIcon, type SelectSize } from 'naive-ui';
    import { type NaiveUISelectOptionWithColor } from '../../types/customNaiveUI';

    interface WorkspaceSelectorProps {
        disabled?: boolean;
        loading?: boolean;
        size?: SelectSize;
        placeholder?: string;
        clearable?: boolean;
        filterable?: boolean;
        style?: string | CSSProperties;
        options: NaiveUISelectOptionWithColor[];
        icon: Component;
    };

    const props = withDefaults(defineProps<WorkspaceSelectorProps>(), {
        size: "medium",
        disabled: false,
        loading: false,
        clearable: false,
    });

    const modelValue = defineModel<string | null>('value');

    const renderLabel = (option: NaiveUISelectOptionWithColor): VNodeChild => {
        if (option.type === 'group')
            return `${option.label}`
        return [
            h(
                NIcon,
                {
                    color: option.color,
                    style: {
                        verticalAlign: '-0.15em',
                        marginRight: '4px'
                    }
                },
                {
                    default: () => h(props.icon)
                }
            ),
            option.label as string
        ]
    }
</script>

<template>
    <n-select :placeholder="placeholder" :size="size" v-model:value="modelValue" :options="options"
        :render-label="renderLabel" :style="style" :disabled="disabled" :loading="loading" :clearable="clearable"
        :filterable="filterable"></n-select>
</template>

<style lang=" css" scoped></style>