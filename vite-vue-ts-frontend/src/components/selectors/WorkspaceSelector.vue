<script setup lang="ts">
    import { ref } from 'vue';
    import { NSelect } from 'naive-ui';

    interface WorkspaceOptionInterface {
        id: string;
        name: string;
    }

    class WorkspaceOption implements WorkspaceOptionInterface {
        id: string;
        name: string;
        constructor(option: WorkspaceOptionInterface) {
            this.id = option.id;
            this.name = option.name;
        }
    }

    const selectedWorkspace = ref<WorkspaceOption | null>(null);

</script>

<template>
    <n-select placeholder="Select workspace" v-model="selectedWorkspace"></n-select>
</template>

<style lang="css" scoped></style>