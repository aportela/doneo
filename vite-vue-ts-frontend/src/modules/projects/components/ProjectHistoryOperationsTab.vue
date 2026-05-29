<script setup lang="ts">
    import type { CSSProperties } from "vue";

    import { NCard, NTimeline, NTimelineItem } from "naive-ui";

    interface ProjectNotesProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectNotesProps>();
</script>

<template>
    <n-card bordered :style="props.style">
        <n-timeline>
            <n-timeline-item content="Oops" />
            <n-timeline-item type="success" title="Success" content="success content" time="2018-04-03 20:46" />
            <n-timeline-item type="error" content="Error content" time="2018-04-03 20:46" />
            <n-timeline-item type="warning" title="Warning" content="warning content" time="2018-04-03 20:46" />
            <n-timeline-item type="info" title="Info" content="info content" time="2018-04-03 20:46"
                line-type="dashed" />
            <n-timeline-item content="Oops" />
        </n-timeline>
    </n-card>
</template>

<style lang="css" scoped></style>