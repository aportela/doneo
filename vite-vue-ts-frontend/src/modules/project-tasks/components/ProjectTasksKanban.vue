<script setup lang="ts">
    import { onMounted, ref } from 'vue'

    import { NCard, NButton, NInput, NScrollbar, NGrid, NGridItem } from 'naive-ui'

    import draggable from 'vuedraggable'

    import { taskStatusService } from '../../task-statuses/services/task-status'
    import type { SearchRequest } from '../../task-statuses/types/dto'
    import { Sort } from '../../../shared/types/models/sort'

    interface Task {
        id: string
        title: string
    }

    interface Column {
        id: string
        title: string
        tasks: Task[]
    }

    const columns = ref<Column[]>([
        { id: 'todo', title: 'To Do', tasks: [{ id: '1', title: 'Tarea A' }] },
        { id: 'doing', title: 'Doing', tasks: [{ id: '2', title: 'Tarea B' }] },
        { id: 'done', title: 'Done', tasks: [{ id: '3', title: 'Tarea C' }] }
    ])

    const newTaskTitle = ref('')

    function addTask(columnId: string) {
        const column = columns.value.find(c => c.id === columnId)
        if (column && newTaskTitle.value.trim() !== '') {
            column.tasks.push({ id: Date.now().toString(), title: newTaskTitle.value })
            newTaskTitle.value = ''
        }
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
                        <draggable v-model="column.tasks" group="tasks" item-key="id" animation="200"
                            ghost-class="drag-ghost">
                            <template #item="{ element }">
                                <NCard size="small" style="cursor: grab; margin-bottom: 8px;">
                                    {{ element.title }}
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
</style>