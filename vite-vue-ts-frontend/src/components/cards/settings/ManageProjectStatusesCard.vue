<script setup lang="ts">
    import { ref, onMounted } from 'vue'
    import { NSpin, NTable, NColorPicker, NTag, NButton, NFlex } from 'naive-ui'
    import { api } from '../../../composables/api';
    import { IconArrowDown, IconArrowUp } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../../../composables/color';

    interface ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    class ProjectStatus implements ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
        constructor(item: ProjectStatusInterface) {
            this.id = item.id;
            this.name = item.name;
            this.index = item.index;
            this.hexColor = item.hexColor;
        }
    }

    const projectStatuses = ref<ProjectStatus[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.projectStatuses.search().then((successResponse: any) => {
            projectStatuses.value = successResponse.data.projectStatuses;
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
                <tr v-for="projectStatus in projectStatuses" :key="projectStatus.id">
                    <td>{{ projectStatus.name }}</td>
                    <td><n-color-picker v-model:value="projectStatus.hexColor" :modes="['hex']"
                            :show-alpha="false"></n-color-picker></td>
                    <td><n-tag :color="getNaiveUITagColorProperty(projectStatus.hexColor)">{{ projectStatus.name
                            }}</n-tag>
                    </td>
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