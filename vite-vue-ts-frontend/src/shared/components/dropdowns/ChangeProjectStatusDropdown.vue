<script setup lang="ts">
    import { computed } from 'vue';
    import { NDropdown, NButton, NIcon, type DropdownOption } from 'naive-ui';
    import { IconStatusChange } from '@tabler/icons-vue';
    import type { TaskStatus } from '../../../modules/task-statuses/models/task-status';
    import { useCacheStore } from '../../../stores/cache';

    interface ChangeProjectStatusDropdownProps {
        currentStatus: TaskStatus
    };

    const cacheStore = useCacheStore();

    const props = defineProps<ChangeProjectStatusDropdownProps>();

    const options = computed<DropdownOption[]>(() =>
        cacheStore.projectStatuses.map(item => ({
            label: item.name ?? '',
            key: item.id ?? '',
            disabled: item.id === props.currentStatus.id
        }))
    );

    const onChange = (key: string | number) => {
        const status = cacheStore.projectStatuses.find(
            item => item.id === key
        );
        console.log(status);
    }
</script>

<template>
    <n-dropdown trigger="click" :options="options" @select="onChange">
        <n-button>
            <template #icon>
                <n-icon :size="22" :component="IconStatusChange" />
            </template>
            Change status
        </n-button>
    </n-dropdown>
</template>

<style lang="css" scoped></style>