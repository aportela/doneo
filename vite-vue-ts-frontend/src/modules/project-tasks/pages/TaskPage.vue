<script setup lang="ts">
    import { ref, computed, watch, nextTick } from 'vue';
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane, type TabsInst } from 'naive-ui';

    //import { IconAlertTriangle } from '@tabler/icons-vue';

    import TaskMetadataTab from '../components/TaskMetadataTab.vue';
    //import ProjectAttachmentsTab from '../components/ProjectAttachmentsTab.vue';
    //import ProjectNotesTab from '../components/ProjectNotesTab.vue';
    //import ProjectHistoryOperationsTab from '../components/ProjectHistoryOperationsTab.vue';

    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.projectId as string ?? ""
    const taskId = route.params.id as string

    // TODO: set tab with type (type tab = "metadata" | "notes"....)
    const tab = computed({
        get: () => route.params.tab as string,
        set: (value: string) => {
            router.push({
                name: 'projectTab',
                params: {
                    id: route.params.id,
                    tab: value
                }
            });
        }
    });

    const noteCount = ref<number>(0);
    const attachmentCount = ref<number>(0);
    const historyOperationCount = ref<number>(0);

    // TODO: i18n
    const attachmentsTabLabel = computed(() => attachmentCount.value > 0 ? `Attachments (${attachmentCount.value})` : 'Attachments')
    const notesTabLabel = computed(() => noteCount.value > 0 ? `Notes (${noteCount.value})` : 'Notes')
    const historyTabLabel = computed(() => historyOperationCount.value > 0 ? `History (${historyOperationCount.value})` : 'History')

    const tabsRef = ref<TabsInst>();


    // recalc bar position on dynamic tab labels changes
    watch(
        () => [attachmentsTabLabel.value, notesTabLabel.value, historyTabLabel.value],
        async () => {
            await nextTick();
            tabsRef.value?.syncBarPosition();
        }
    );
</script>

<template>
    <n-tabs placement="top" type="line" ref="tabsRef" animated v-model:value="tab">
        <n-tab-pane name="metadata" display-directive="show" key="metadata">
            <template #tab>
                Metadata
            </template>
            <!-- TODO: mode ???-->
            <TaskMetadataTab mode="add" :project-id="projectId" :task-id="taskId" v-model:note-count="noteCount"
                v-model:attachment-count="attachmentCount" v-model:history-operation-count="historyOperationCount" />
        </n-tab-pane>
        <n-tab-pane name="notes" :tab="notesTabLabel" display-directive="show:lazy">
            <!--
            <ProjectNotesTab :project-id="taskId" v-model:item-count="noteCount" />
            -->
        </n-tab-pane>
        <n-tab-pane name="attachments" :tab="attachmentsTabLabel" display-directive="show:lazy">
            <!--
            <ProjectAttachmentsTab :project-id="taskId" v-model:item-count="attachmentCount" />
            -->
        </n-tab-pane>
        <n-tab-pane name="history" :tab="historyTabLabel" display-directive="show:lazy">
            <!--
            <ProjectHistoryOperationsTab :project-id="taskId" v-model:item-count="historyOperationCount"
                :key="historyOperationCount" />
                -->
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
    }
</style>