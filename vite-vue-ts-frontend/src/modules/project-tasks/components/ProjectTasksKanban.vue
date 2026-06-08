<script setup lang="ts">
    import { onMounted, ref } from 'vue'

    import { NCard, NButton, NInput, NScrollbar, NGrid, NGridItem, NIcon } from 'naive-ui'

    import draggable from 'vuedraggable'

    import { taskStatusService } from '../../task-statuses/services/task-status'
    import type { SearchRequest as SearchRequestStatus } from '../../task-statuses/types/dto'

    import { projectTaskService } from '../services/task'
    import type { SearchRequest as SearchRequestTask } from '../types/dto'

    import { Sort } from '../../../shared/types/models/sort'
    import { IconPaperclip } from '@tabler/icons-vue'

    interface Task {
        id: string
        slug: string;
        summary: string
        attachmentCount: number;
    }


    interface Column {
        id: string
        title: string
        hexColor: string;
        tasks: Task[]
    }

    const columns = ref<Column[]>([
    ]);

    const newTaskTitle = ref('')

    function addTask(columnId: string) {
        const column = columns.value.find(c => c.id === columnId)
        if (column && newTaskTitle.value.trim() !== '') {
            column.tasks.push({ id: Date.now().toString(), slug: "00", summary: newTaskTitle.value, attachmentCount: 0 })
            newTaskTitle.value = ''
        }
    }

    const statusSort = new Sort("index", "ASC");

    const tasksSort = new Sort("createdAt", "DESC");

    const refreshStatus = async () => {
        const payload: SearchRequestStatus = {
            pager: {
                currentPage: 1,
                resultsPage: 0,
            },
            order: {
                field: statusSort.field,
                sort: statusSort.order,
            },
            filter: {
                //name: filters.name.length > 0 ? filters.name : undefined,
            }
        };
        try {
            const response = await taskStatusService.search(payload);
            columns.value = response.taskStatuses.map((taskStatus) => { return { id: taskStatus.id, title: taskStatus.name, hexColor: taskStatus.hexColor, tasks: [] } });
        } catch (e) {
            console.error(e);
        }
    }

    const refreshTasks = async () => {
        const payload: SearchRequestTask = {
            pager: {
                currentPage: 1,
                resultsPage: 0,
            },
            order: {
                field: tasksSort.field,
                sort: tasksSort.order,
            },
            filter: {
                //name: filters.name.length > 0 ? filters.name : undefined,
            }
        };
        try {
            const response = await projectTaskService.search(null, payload);
            columns.value.forEach((column) => {
                column.tasks = response.tasks.filter((task) => task.status.id == column.id).map((task) => { return { id: task.id, slug: task.slug, summary: task.summary, attachmentCount: Math.floor(Math.random() * 3) }; });
            });
        } catch (e) {
            console.error(e);
        }
    }

    onMounted(async () => {
        try {
            await refreshStatus();
            await refreshTasks();
        } catch (e) {
            console.error(e);
        }
    });
</script>

<template>
    <NGrid :cols="columns.length" x-gap="16">
        <NGridItem v-for="column in columns" :key="column.id">
            <NCard :title="column.title" size="large" class="doneo-kanban-column"
                :style="'background-color: ' + column.hexColor + ';'">
                <div style="flex: 1; overflow: hidden;">
                    <NScrollbar style="height: 100%;">
                        <draggable v-model="column.tasks" group="tasks" item-key="id" animation="200"
                            ghost-class="drag-ghost">
                            <template #item="{ element }">
                                <NCard size="small" style="cursor: grab; margin-bottom: 8px;">
                                    {{ element.slug }}
                                    <p>{{ element.summary }}</p>
                                    <img v-if="element.attachmentCount > 0"
                                        :src="'https://loremflickr.com/320/240?random=' + element.id">
                                    <p v-if="element.attachmentCount > 0"><n-icon :component="IconPaperclip" /> {{
                                        element.attachmentCount }}</p>
                                </NCard>
                            </template>
                        </draggable>
                    </NScrollbar>
                </div>
                <div style="margin-top: 8px; display: flex; gap: 6px;">
                    <NInput v-model:value="newTaskTitle" placeholder="New task..." size="small" />
                    <NButton size="small" @click="addTask(column.id)">+</NButton>
                </div>
            </NCard>
        </NGridItem>
    </NGrid>
</template>

<style scoped>
    .drag-ghost {
        opacity: 0.5;
        background-color: #f0f0f0;
    }

    .doneo-kanban-column {
        min-height: 500px;
        display: flex;
        flex-direction: column;
    }
</style>