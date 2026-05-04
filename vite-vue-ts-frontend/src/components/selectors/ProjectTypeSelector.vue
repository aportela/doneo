<script setup lang="ts">
    import { h, reactive, onMounted, type CSSProperties, shallowRef, type VNodeChild } from 'vue';
    import { useI18n } from "vue-i18n";
    import { type SelectOption, NSelect, NIcon } from 'naive-ui';
    import { api } from '../../composables/api';
    import { type SearchProjectTypesResponse } from '../../types/apiResponses';
    import { type AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";
    import { IconSquareFilled } from '@tabler/icons-vue';

    const { t } = useI18n();

    interface ProjectTypeSelectorProps {
        workspaceId: string;
        style?: string | CSSProperties;
    };

    const props = defineProps<ProjectTypeSelectorProps>();

    const emit = defineEmits(['error'])

    const selectedProjectType = defineModel<string | null>('value');

    export interface SelectOptionWithColor extends SelectOption {
        color: string;
    }

    const options = shallowRef<SelectOptionWithColor[]>([]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const findValidProjectType = (options: SelectOptionWithColor[], stored?: string | null) => {
        if (!options.length) return null;

        if (stored && options.some(o => o.value === stored)) {
            return stored;
        } else {
            return null
        }
    };

    const onRefresh = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.projectTypes.search(props.workspaceId).then((response: SearchProjectTypesResponse) => {
            options.value = response.data.projectTypes.map((item) => ({
                value: item.id,
                label: item.name,
                disabled: false,
                color: item.hexColor,
            }));
            selectedProjectType.value = findValidProjectType(
                options.value,
                props.workspaceId
            );
        }).catch((error: any) => {
            emit("error", { error: error });
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const renderLabel = (option: SelectOptionWithColor): VNodeChild => {
        if (option.type === 'group')
            return `${option.label}(Cool!)`
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
                    default: () => h(IconSquareFilled)
                }
            ),
            option.label as string
        ]
    }

    onMounted(() => {
        onRefresh();
    });
</script>

<template>
    <n-select :placeholder="t('Project type')" v-model:value="selectedProjectType" :options="options"
        :render-label="renderLabel" :style="style" :disabled="state.ajaxRunning"
        :loading="state.ajaxRunning"></n-select>
</template>

<style lang="css" scoped></style>