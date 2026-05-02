<script setup lang="ts">
    import { ref, onMounted } from 'vue'
    import { NSpin, NTable, NColorPicker, NTag, NButton, NFlex } from 'naive-ui'
    import { api } from '../../../composables/api';
    import { IconArrowDown, IconArrowUp } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../../../composables/color';

    interface ProjectPriorityInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    const projectPriorities = ref<ProjectPriorityInterface[]>([]);

    const loading = ref<boolean>(false);

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
                <tr v-for="projectPriority in projectPriorities" :key="projectPriority.id">
                    <td>{{ projectPriority.name }}</td>
                    <td><n-color-picker v-model:value="projectPriority.hexColor" :modes="['hex']"
                            :show-alpha="false"></n-color-picker></td>
                    <td><n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)">{{ projectPriority.name
                            }}</n-tag></td>
                    <td>
                        <n-flex>
                            <n-button>
                                <IconArrowUp />
                            </n-button>
                            <n-button>
                                <IconArrowDown />
                            </n-button>
                        </n-flex>
                    </td>
                </tr>
            </tbody>
        </n-table>
    </n-spin>
</template>

<style lang="css" scoped></style>