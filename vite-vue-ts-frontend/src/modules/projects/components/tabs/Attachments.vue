<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { useNotify } from '../../../../shared/composables/notification';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../attachments/types/dto.ts";

    import { ProjectAttachment } from "../../../attachments/models/project-attachment.ts";

    import { projectAttachmentService } from "../../../attachments/services/project-attachment.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import UploadDialog from "../../../attachments/components/UploadDialog.vue";
    import ImagePreview from "../../../../shared/components/ImagePreview.vue";

    import ProjectAttachmentsTable from "../../../attachments/components/ProjectAttachmentsTable.vue";
    import type { ProjectAttachmentsTableFilters } from "../../../attachments/types/project-attachments-table-filter.ts";

    interface ProjectAttachmentsProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const props = defineProps<ProjectAttachmentsProps>();

    const emit = defineEmits(['delete']);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectAttachment[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const uploadCount = ref<number>(0);

    const filters = reactive<ProjectAttachmentsTableFilters>({
        name: "",
        createdByUserId: null,
        createdAt: {
            from: null,
            to: null,
        },
        contentType: null,
    });

    const nameFilterLowerCase = computed(() =>
        filters.name.toLowerCase()
    );

    const filteredItems = computed(() => {
        return items.value.filter((attachment: ProjectAttachment) => {
            const name = attachment.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value)) &&
                (filters.createdByUserId === null || filters.createdByUserId === attachment.createdBy.id) &&
                (filters.contentType === null || filters.contentType === attachment.contentType) &&
                ((filters.createdAt.from === null && filters.createdAt.to === null) || (attachment.createdAt.msTimestamp != null && filters.createdAt.from != null && filters.createdAt.from <= attachment.createdAt.msTimestamp && filters.createdAt.to != null && filters.createdAt.to >= attachment.createdAt.msTimestamp))
            );
        });
    });

    const showUploadModal = ref<boolean>(false);

    const showImagePreviewModal = ref<boolean>(false);

    const imageSources = computed<string[]>(() => items.value.filter((item: ProjectAttachment) => item.allowImagePreview()).map((item: ProjectAttachment) => item.getDownloadURL(props.projectId)));

    const imageSourcesWithIds = computed(() => items.value.filter((item: ProjectAttachment) => item.allowImagePreview()).map((item: ProjectAttachment) => {
        return ({
            id: item.id,
            url: item.getDownloadURL(props.projectId),
        });
    }));

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

    watch(showUploadModal, (newValue) => {
        if (!newValue) {
            if (uploadCount.value > 0) {
                onRefresh();
            }
        }
    });

    const onShowUploadModal = () => {
        uploadCount.value = 0;
        showUploadModal.value = true;
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
                            state.ajaxErrorMessage = t("modules.projectAttachment.components.ProjectAttachmentsTab.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectAttachment.components.ProjectAttachmentsTab.errors.refreshError");
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
                notify('success', t("modules.projectAttachment.components.ProjectAttachmentsTab.notifications.projectAttachmentDeleted", { name: projectAttachment.name }));
                onRefresh();
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
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.ProjectAttachmentsTab.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.ProjectAttachmentsTab.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectAttachment.components.ProjectAttachmentsTab.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("project attachment id not set", { file: "ProjectAttachmentsTab.vue", method: "onDelete" });
        }
    };

    const onDownload = (_projectAttachment: ProjectAttachment, _index: number) => { };

    const onPreview = (_projectAttachment: ProjectAttachment, _index: number) => {
        currentImagePreviewIndex.value = imageSourcesWithIds.value.findIndex((item) => item.id == _projectAttachment.id);
        showImagePreviewModal.value = true;
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

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <ImagePreview v-model:show="showImagePreviewModal" :sources="imageSources"
        :current-index="currentImagePreviewIndex" />
    <!-- TODO: onupload notification -->
    <UploadDialog v-if="props.projectId" v-model:show="showUploadModal" :project-id="props.projectId"
        v-model:upload-count="uploadCount" />
    <n-card bordered :style="props.style">
        <ProjectAttachmentsTable :project-id="props.projectId" :items="filteredItems" :disabled="state.ajaxRunning"
            v-model:filters="filters" @refresh="onRefresh" @add="onShowUploadModal" @delete="onDelete"
            @download="onDownload" @preview="onPreview" />
    </n-card>
</template>

<style lang="css" scoped></style>