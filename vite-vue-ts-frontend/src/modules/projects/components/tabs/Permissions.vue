<script setup lang="ts">
    import { NCard } from "naive-ui";

    import ProjectPermissionsTable from "../../../project-permissions/components/ProjectPermissionsTable.vue";

    interface Props {
        readOnly?: boolean;
        projectId: string;
    }

    const props = defineProps<Props>();

    const itemCount = defineModel<number>("itemCount", { default: 0 });
</script>

<template>
    <n-card bordered>
        <ProjectPermissionsTable :read-only="props.readOnly" :project-id="props.projectId"
            v-model:item-count="itemCount" />
    </n-card>
</template>

<style lang="css" scoped></style>