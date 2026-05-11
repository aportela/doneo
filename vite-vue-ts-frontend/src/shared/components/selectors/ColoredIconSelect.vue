<script setup lang="ts">
    import { h, type CSSProperties, type VNodeChild, type Component, watch } from 'vue';
    import { NSelect, NIcon, type SelectSize } from 'naive-ui';
    import { type NaiveUISelectOptionWithColor } from '../../types/customNaiveUI';
    import { IconSquareFilled } from '@tabler/icons-vue';


    interface ColoredIconSelectProps {
        size?: SelectSize;
        disabled?: boolean;
        loading?: boolean;
        clearable?: boolean;
        filterable?: boolean;
        placeholder?: string;
        style?: string | CSSProperties;
        options: NaiveUISelectOptionWithColor[];
        icon?: Component;
    };


    const props = withDefaults(defineProps<ColoredIconSelectProps>(), {
        size: "medium",
        disabled: false,
        loading: false,
        clearable: false,
        filterable: false,
    });

    const modelValue = defineModel<string | null>('value');

    watch(modelValue, (val) => {
        const exists = props.options.some(opt => opt.value === val)
        if (!exists && val !== null) {
            console.warn(`modelValue "${val}" does not exist in options`, props.options)
            modelValue.value = null
        }
    }, { immediate: true })

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
                    default: () => h(props.icon || IconSquareFilled)
                }
            ),
            option.label as string
        ]
    }
</script>

<template>
    <n-select :placeholder="placeholder" :size="size" v-model:value="modelValue" :options="options"
        :render-label="renderLabel" :style="style" :disabled="disabled" :loading="loading" :clearable="clearable"
        :filterable="filterable">
    </n-select>
</template>

<style lang=" css" scoped></style>