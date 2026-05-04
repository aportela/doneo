<script setup lang="ts">
    import { reactive, shallowRef, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";
    import { api } from '../../composables/api';
    import { type SearchWorkspacesResponse } from '../../types/apiResponses';
    import { type AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";
    import { type NaiveUISelectOptionWithColor } from '../../types/customNaiveUI';
    import { default as ColoredIconSelect } from './ColoredIconSelect.vue';

    const { t } = useI18n();

    const emit = defineEmits(['error'])

    const modelValue = defineModel<string | null>('value');

    const options = shallowRef<NaiveUISelectOptionWithColor[]>([]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const onRefresh = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.workspace.search().then((response: SearchWorkspacesResponse) => {
            options.value = response.data.workspaces.map((item) => ({
                value: item.id,
                label: item.name,
                disabled: false,
                color: item.hexColor,
            }));
        }).catch((error: any) => {
            emit("error", { error: error });
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    onMounted(() => {
        onRefresh();
    });
</script>

<template>
    <ColoredIconSelect :placeholder="t('Workspace')" v-model:value="modelValue" :options="options"
        :filterable="false" />
</template>

<style lang="css" scoped></style>