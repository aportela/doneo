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
    import ImagePreview from "../../../shared/components/ImagePreview.vue";

    interface ProjectAttachmentsProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectAttachmentsProps>();

    const emit = defineEmits(['delete']);


    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectAttachment[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const uploadCount = ref<number>(0);

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

    const showImagePreviewDialog = ref<boolean>(false);

    const imageSources = computed<string[]>(() => items.value.filter((item: ProjectAttachment) => item.allowImagePreview()).map((item: ProjectAttachment) => item.getDownloadURL(props.projectId)));
    const imageSourcesWithIds = computed(() => items.value.filter((item: ProjectAttachment) => item.allowImagePreview()).map((item: ProjectAttachment) => {
        return (
            {
                id: item.id,
                url: item.getDownloadURL(props.projectId),
            });
    }
    ));

    const currentImagePreviewIndex = ref<number>(0);

    const selectedItem = ref<ProjectAttachment>(new ProjectAttachment());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    watch(showUploadDialog, (newValue) => {
        if (!newValue) {
            if (uploadCount.value > 0) {
                onRefresh();
            }
        }
    });

    const onShowUploadDialog = () => {
        uploadCount.value = 0;
        showUploadDialog.value = true;
    };

    const onRefresh = async () => {
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
    };

    const onDelete = async (projectAttachment: ProjectAttachment, _index?: number) => {
        if (projectAttachment.id) {
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
            console.error("project attachment id not set", { file: "ProjectAttachmentsTab.vue", method: "onDelete" });
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectAttachmentsTab.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectAttachmentsTab.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    const onDownload = (_projectAttachment: ProjectAttachment, _index: number) => { };

    const onPreview = (_projectAttachment: ProjectAttachment, _index: number) => {
        currentImagePreviewIndex.value = imageSourcesWithIds.value.findIndex((item) => item.id == _projectAttachment.id);
        showImagePreviewDialog.value = true;
    };

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <ImagePreview v-model:show="showImagePreviewDialog" :sources="imageSources"
        :current-index="currentImagePreviewIndex" />
    <UploadDialog v-if="props.projectId" v-model:show="showUploadDialog" :project-id="props.projectId"
        v-model:upload-count="uploadCount" />
    <n-card bordered :style="props.style">
        <ProjectAttachmentsTable :project-id="props.projectId" :project-attachments="filteredAttachments"
            :loading="state.ajaxRunning" v-model:name-filter="filterByName" v-model:user-filter="filterByUser"
            @refresh="onRefresh" @add="onShowUploadDialog" @delete="onDelete" @download="onDownload"
            @preview="onPreview" />
    </n-card>
</template>

<style lang="css" scoped></style>