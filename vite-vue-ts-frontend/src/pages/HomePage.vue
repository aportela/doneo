<script setup lang="ts">
    import { NGrid, NGridItem, NCard } from 'naive-ui';
    import { default as ProjectTypeSelect } from '../components/selectors/ProjectTypeSelect.vue';
    import { ref } from 'vue';


    const selected = ref<string | null>(null);

</script>

<template>
    <h1>Overview</h1>
    <h2>Dashboard</h2>

    <p>Es: {{ selected }}</p>
    <ProjectTypeSelect placeholder="Seleccione un tipo" v-model:value="selected" filterable
        workspace-id="019df284-0e87-7883-96e3-5945412a87bf" />


    <n-grid cols="1 1600:2" :x-gap="32" :y-gap="32" responsive="self" item-responsive style="margin-top: 1em;">
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
    </n-grid>
</template>

<style lang="css" scoped></style>