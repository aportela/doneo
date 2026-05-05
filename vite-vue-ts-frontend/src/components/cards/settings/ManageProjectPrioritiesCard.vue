<script setup lang="ts">
    import { ref, onMounted } from 'vue'
    import { useI18n } from "vue-i18n";
    import { NSpin, NTable, NColorPicker, NTag, NButton, NButtonGroup } from 'naive-ui'
    import { api } from '../../../composables/api';
    import { IconArrowDown, IconArrowUp } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../../../composables/color';

    const { t } = useI18n();

    interface ProjectPriorityInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    const projectPriorities = ref<ProjectPriorityInterface[]>([]);

    const loading = ref<boolean>(false);

    const onChangeIndex = (projectPriorityId: string, index: number) => {
        api.projectPriorities.patchIndex(projectPriorityId, index);
    }

    onMounted(() => {
        loading.value = true;
        api.projectPriorities.search().then((successResponse: any) => {
            projectPriorities.value = successResponse.data.projectPriorities;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });
</script>

<template>
    <n-spin :show="loading">
        <n-table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Color</th>
                    <th>Preview</th>
                    <th>Move Index/Position</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="projectPriority, index in projectPriorities" :key="projectPriority.id">
                    <td>{{ projectPriority.name }}</td>
                    <td><n-color-picker v-model:value="projectPriority.hexColor" :modes="['hex']"
                            :show-alpha="false"></n-color-picker></td>
                    <td><n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)">{{ projectPriority.name
                            }}</n-tag></td>
                    <td>
                        <n-button-group>
                            <n-button size="small" @click="onChangeIndex(projectPriority.id, index - 1)"
                                :disabled="index === 0">
                                {{ t("Move up") }}
                                <IconArrowUp />
                            </n-button>
                            <n-button size="small" @click="onChangeIndex(projectPriority.id, index + 1)"
                                :disabled="index >= projectPriorities.length - 1">
                                {{ t("Move down") }}
                                <IconArrowDown />
                            </n-button>
                        </n-button-group>
                    </td>
                </tr>
            </tbody>
        </n-table>
    </n-spin>
</template>

<style lang="css" scoped></style>