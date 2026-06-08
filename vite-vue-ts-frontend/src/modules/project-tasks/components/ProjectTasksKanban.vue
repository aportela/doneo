<script setup lang="ts">
    import { onMounted, ref } from 'vue'

    import { NCard, NButton, NInput, NScrollbar, NIcon, NGrid, NGridItem } from 'naive-ui'
    import { IconPlus } from '@tabler/icons-vue'

    import { taskStatusService } from '../../task-statuses/services/task-status'
    import type { SearchRequest } from '../../task-statuses/types/dto'
    import { Sort } from '../../../shared/types/models/sort'

    interface Task {
        id: number
        title: string
        description?: string
    }

    interface Column {
        id: number
        title: string
        tasks: Task[]
    }

    const sort = new Sort("name", "ASC");

    const refreshStatus = async () => {
        const payload: SearchRequest = {
            pager: {
                currentPage: 1,
                resultsPage: 0,
            },
            order: {
                field: sort.field,
                sort: sort.order,
            },
            filter: {
                //name: filters.name.length > 0 ? filters.name : undefined,
            }
        };
        taskStatusService.search(payload);
    }

    const columns = ref<Column[]>([
        { id: 1, title: 'To Do', tasks: [{ id: 1, title: 'Task 1' }, { id: 2, title: 'Task 2' }] },
        { id: 2, title: 'In Progress', tasks: [{ id: 3, title: 'Task 3' }] },
        { id: 3, title: 'Done', tasks: [{ id: 4, title: 'Task 4' }] }
    ])

    const newTaskTitle = ref('')
    const addTask = (columnId: number) => {
        if (!newTaskTitle.value.trim()) return
        const column = columns.value.find(c => c.id === columnId)
        if (!column) return
        column.tasks.push({
            id: Date.now(),
            title: newTaskTitle.value
        })
        newTaskTitle.value = ''
    }

    onMounted(() => {
        refreshStatus();
    });
</script>

<template>
    <NGrid cols="3" x-gap="16">
        <NGridItem v-for="column in columns" :key="column.id">
            <NCard :title="column.title" size="large" style="height: 500px; display: flex; flex-direction: column;">
                <div style="flex: 1; overflow: hidden;">
                    <NScrollbar style="height: 100%;">
                        <div style="display: flex; flex-direction: column; gap: 8px;">
                            <NCard v-for="task in column.tasks" :key="task.id" size="small" style="cursor: grab;">
                                {{ task.title }}
                            </NCard>
                        </div>
                    </NScrollbar>
                </div>
                <div style="margin-top: 8px; display: flex; gap: 6px;">
                    <NInput v-model:value="newTaskTitle" placeholder="Nueva tarea..." size="small" />
                    <NButton size="small" @click="addTask(column.id)">
                        +
                    </NButton>
                </div>
            </NCard>
        </NGridItem>
    </NGrid>
</template>

<style scoped>

    .n-scrollbar__wrap {
        max-height: 300px;
    }
</style>