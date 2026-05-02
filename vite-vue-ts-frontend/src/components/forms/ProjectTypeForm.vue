<script setup lang="ts">
    import { NCard, NInput, NFlex, NButton } from 'naive-ui';
    import { IconCancel, IconPlus } from '@tabler/icons-vue';

    const emit = defineEmits(['cancel'])

    const onCancel = () => {
        emit('cancel')
    }
    const onAddProjectType = () => { };
</script>

<template>
    <n-card title="Add new project type" style="width: 400px;">
        <n-input placeholder="Project type name"></n-input>
        <template #action>
            <n-flex>
                <n-button @click="onAddProjectType">
                    <template #icon>
                        <IconPlus />
                    </template>
                    Add new
                </n-button>
                <n-button @click="onCancel">
                    <template #icon>
                        <IconCancel />
                    </template>
                    Cancel
                </n-button>
            </n-flex>
        </template>
    </n-card>
</template>

<style lang="css" scoped></style>