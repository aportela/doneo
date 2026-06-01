<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';

    import { appBus } from '../../../shared/composables/bus';

    import { projectAttachmentService } from "../../attachments/services/project-attachment.ts";
    import { handleAPIError } from '../../../api/client/errorHandler';

    import type { SearchResponse } from "../../attachments/types/dto.ts";

    import { ProjectAttachment } from "../../attachments/models/project-attachment.ts";

    import UploadDialog from "../../attachments/components/UploadDialog.vue";
    import ProjectAttachmentsTable from "../../attachments/components/ProjectAttachmentsTable.vue";

    interface ProjectAttachmentsProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectAttachmentsProps>();
    const emit = defineEmits(['delete']);


    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectAttachment[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filterByUser = ref<string>("");

    const filterByName = ref<string>("");

    const nameFilter = computed(() =>
        filterByName.value?.toLowerCase() ?? ''
    );

    const filteredAttachments = computed(() => {
        return items.value.filter((permission) => {
            const fileName = permission.name?.toLowerCase();
            return (
                (!nameFilter.value || fileName?.includes(nameFilter.value)) &&
                (!filterByUser.value || filterByUser.value === permission.createdBy.id)
            );
        });
    });

    const showUploadDialog = ref<boolean>(false);

    const selectedItem = ref<ProjectAttachment>(new ProjectAttachment());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onShowUploadDialog = () => {
        showUploadDialog.value = true;
    };

    const onRefresh = async () => {
        if (props.projectId) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const results: SearchResponse = await projectAttachmentService.getProjectAttachments(props.projectId);
                items.value = results.attachments.map((attachment) => new ProjectAttachment(attachment));
                itemCount.value = items.value?.length ?? 0;
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectAttachmentsTab.onRefresh" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                        console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("project id not set", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" });
        }
    };

    const onDelete = async (projectAttachment: ProjectAttachment, _index?: number) => {
        if (props.projectId && projectAttachment.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectAttachmentService.deleteProjectAttachment(props.projectId, projectAttachment.id);
                notify('success', t("modules.projectPermission.components.projectPermissions.notifications.projectPermissionDeleted", {}));
                onRefresh();
                // TODO: not refreshing permission count on parent component
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = projectAttachment;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectAttachmentsTab.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "ProjectAttachmentsTab.vue", method: "onDelete" });
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        if (props.projectId) {
            onRefresh();
        }
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectAttachmentsTab.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectAttachmentsTab.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>

    <UploadDialog v-if="props.projectId" v-model:show="showUploadDialog" :project-id="props.projectId" />
    <n-card bordered :style="props.style">
        <ProjectAttachmentsTable :project-attachments="filteredAttachments" :loading="state.ajaxRunning"
            v-model:name-filter="filterByName" v-model:user-filter="filterByUser" @refresh="onRefresh"
            @add="onShowUploadDialog" @delete="onDelete" />
    </n-card>
</template>

<style lang="css" scoped></style>