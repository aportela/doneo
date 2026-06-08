<script setup lang="ts">
    import { ref } from 'vue'
    import { NCalendar, NCard, NTag } from 'naive-ui'

    interface Task {
        id: string
        title: string
        date: string // YYYY-MM-DD
    }

    const tasks = ref<Task[]>([
        { id: '1', title: 'Diseñar UI', date: '2026-06-08' },
        { id: '2', title: 'Fix bugs', date: '2026-06-08' },
        { id: '3', title: 'Deploy', date: '2026-06-10' }
    ])

    function getTasksByDate(date: string) {
        return tasks.value.filter(t => t.date === date)
    }

</script>

<template>
    <NCalendar>
        <template #default="{ year, month, date }">
            <NCard size="small" style="min-height: 80px">
                <!--
                <div style="font-size: 12px; opacity: 0.6">
                    {{ date }}
                </div>
                -->
                <div v-for="task in getTasksByDate(`${year}-${String(month).padStart(2, '0')}-${String(date).padStart(2, '0')}`)"
                    :key="task.id">
                    <NTag size="small" type="info">
                        {{ task.title }}
                    </NTag>
                </div>
            </NCard>
        </template>
    </NCalendar>
</template>