<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { useNotify } from '../../../../shared/composables/notification';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../attachments/types/dto.ts";

    import { Attachment } from "../../../attachments/models/attachment.ts";

    import { attachmentService } from "../../../attachments/services/attachment.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import UploadDialog from "../../../attachments/components/UploadDialog.vue";
    import ImagePreview from "../../../../shared/components/ImagePreview.vue";

    import AttachmentsTable from "../../../attachments/components/AttachmentsTable.vue";
    import type { AttachmentsTableFilters } from "../../../attachments/types/attachments-table-filter.ts";
    import { bgDownload } from "../../../../shared/composables/axios.ts";
    import AudioPreview from "../../../../shared/components/AudioPreview.vue";
    import PDFPreview from "../../../../shared/components/PDFPreview.vue";

    interface ProjectAttachmentsProps {
        disabled?: boolean;
        readOnly?: boolean;
        style?: string | CSSProperties;
        projectId: string;
        taskId: string;
    }

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const props = defineProps<ProjectAttachmentsProps>();

    const emit = defineEmits(['delete']);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Attachment[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const uploadCount = ref<number>(0);

    const filters = reactive<AttachmentsTableFilters>({
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
        return items.value.filter((attachment: Attachment) => {
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

    const imageSources = computed<string[]>(() => items.value.filter((item: Attachment) => item.allowImagePreview()).map((item: Attachment) => item.getPreviewURL(props.projectId, props.taskId)));

    const imageSourcesWithIds = computed(() => items.value.filter((item: Attachment) => item.allowImagePreview()).map((item: Attachment) => {
        return ({
            id: item.id,
            url: item.getDownloadURL(props.projectId, props.taskId),
        });
    }));

    const currentImagePreviewIndex = ref<number>(0);

    const showAudioPreviewModal = ref<boolean>(false);

    const audioSources = computed<Attachment[]>(() => items.value.filter((item: Attachment) => item.allowAudioPreview()));

    const currentAudioPreviewIndex = ref<number>(0);

    const showPDFPreviewModal = ref<boolean>(false);

    const pdfSources = computed<Attachment[]>(() => items.value.filter((item: Attachment) => item.allowPDFPreview()));

    const currentPDFPreviewIndex = ref<number>(0);

    const selectedItem = ref<Attachment>(new Attachment());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
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
            const results: SearchResponse = await attachmentService.getTaskAttachments(props.projectId, props.taskId);
            items.value = results.attachments.map((attachment) => new Attachment(attachment));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskAttachmentsTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectAttachment.components.TaskAttachmentsTab.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectAttachment.components.TaskAttachmentsTab.errors.refreshError");
                    console.error("Unhandled API error", { file: "TaskAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (projectAttachment: Attachment, _index?: number) => {
        if (projectAttachment.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await attachmentService.deleteTaskAttachment(props.projectId, props.taskId, projectAttachment.id);
                notify('success', t("modules.projectAttachment.components.TaskAttachmentsTab.notifications.projectAttachmentDeleted", { name: projectAttachment.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = projectAttachment;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskAttachmentsTab.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.TaskAttachmentsTab.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.TaskAttachmentsTab.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectAttachment.components.TaskAttachmentsTab.errors.deleteError");
                        console.error("Unhandled API error", { file: "TaskAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project attachment id not set", { file: "TaskAttachmentsTab.vue", method: "onDelete" });
        }
    };

    const onDownload = (attachment: Attachment, _index: number) => {
        bgDownload(attachment.getBgDownloadURL(props.projectId, props.taskId), attachment.name)
    };

    const onPreview = (attachment: Attachment, _index: number) => {
        if (attachment.allowImagePreview()) {
            currentImagePreviewIndex.value = imageSourcesWithIds.value.findIndex((item) => item.id == attachment.id);
            showImagePreviewModal.value = true;
        } else if (attachment.allowAudioPreview()) {
            currentAudioPreviewIndex.value = audioSources.value.findIndex((item) => item.id == attachment.id);
            showAudioPreviewModal.value = true;
        } else if (attachment.allowPDFPreview()) {
            currentPDFPreviewIndex.value = pdfSources.value.findIndex((item) => item.id == attachment.id);
            showPDFPreviewModal.value = true;
        } else {
            console.error("Invalid preview");
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskAttachmentsTab.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("TaskAttachmentsTab.onDelete")) {
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
    <AudioPreview v-model:show="showAudioPreviewModal" :project-id="props.projectId" :task-id="props.taskId"
        :items="audioSources" :current-index="currentAudioPreviewIndex" />
    <PDFPreview v-model:show="showPDFPreviewModal" :project-id="props.projectId" :task-id="props.taskId"
        :items="pdfSources" :current-index="currentPDFPreviewIndex" />

    <!-- TODO: onupload notification -->
    <UploadDialog v-if="props.projectId && props.taskId && !props.readOnly" v-model:show="showUploadModal"
        :project-id="props.projectId" :task-id="props.taskId" v-model:upload-count="uploadCount" />
    <n-card bordered :style="props.style">
        <AttachmentsTable :project-id="props.projectId" :items="filteredItems" :disabled="state.ajaxRunning"
            v-model:filters="filters" @refresh="onRefresh" @add="onShowUploadModal" @delete="onDelete"
            @download="onDownload" @preview="onPreview" :read-only="props.readOnly" />
    </n-card>
</template>

<style lang="css" scoped></style>