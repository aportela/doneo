<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NIcon, NButton, NButtonGroup, NSelect, type SelectOption } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { Attachment } from '../models/attachment.ts';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { attachmentService } from '../services/attachment.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';

    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { formatBytes } from '../../../shared/composables/format.ts';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';
    import { renderIcon, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';

    import UploadDialog from './UploadDialog.vue';
    import { bgDownload } from '../../../shared/composables/axios.ts';
    import ImagePreview from '../../../shared/components/widgets/preview/ImagePreview.vue';
    import AudioPreview from '../../../shared/components/widgets/preview/AudioPreview.vue';
    import PDFPreview from '../../../shared/components/widgets/preview/PDFPreview.vue';
    import type { SearchResponse } from '../types/dto.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_DOWNLOAD, DONEO_ICON_ACTION_PREVIEW } from '../../../shared/types/icons.ts';

    interface Props {
        id?: string;
        readOnly?: boolean;
        projectId: string;
        taskId?: string;
    }

    // TODO: this component is not completed, only warnings removed, missing i18n labels & code cleanups
    const props = withDefaults(defineProps<Props>(), { id: "AttachmentsTable" });

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const { t } = useI18n();
    const dialog = useDialog();
    const { notify } = useNotify();
    const userSettingsStore = useUserSettingsStore();
    const loadingStore = useLoadingStore();
    const tableSettingsStore = useTableSettingsStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<Attachment[]>([]);

    const tmpItem = ref<Attachment>(new Attachment());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                a.name.localeCompare(b.name)
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                b.name.localeCompare(a.name)
            );
        }
    };

    const uploadedFilesCount = ref<number>(0);

    const showUploadModal = ref<boolean>(false);


    watch(showUploadModal, (newValue) => {
        if (!newValue) {
            if (uploadedFilesCount.value > 0) {
                onRefresh();
            }
        }
    });

    const onShowUploadModal = () => {
        uploadedFilesCount.value = 0;
        showUploadModal.value = true;
    };

    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface AttachmentsTableFilters {
        name: string;
        createdByUserId: string | null;
        createdAt: TimestampRange;
        contentType: string | null;
    };

    const filters = reactive<AttachmentsTableFilters>(
        {
            name: "",
            createdByUserId: null,
            createdAt: {
                from: null,
                to: null,
            },
            contentType: null,
        }
    );

    const isFilteredByName = computed<boolean>(() => filters.name !== "");
    const isFilteredByCreator = computed<boolean>(() => filters.createdByUserId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByContentType = computed<boolean>(() => filters.contentType !== null);

    const onClearFilters = () => {
        filters.name = "";
        filters.createdByUserId = null;
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        filters.contentType = null;
    };

    const nameFilterLowerCase = computed(() =>
        filters.name.toLowerCase()
    );

    const localFilteredItems = computed(() => {
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

    const contentTypeOptions = ref<SelectOption[]>([]);

    const showImagePreviewModal = ref<boolean>(false);


    const imageSources = computed<string[]>(() => items.value.filter((item: Attachment) => item.allowImagePreview()).map((item: Attachment) => props.taskId ? item.getPreviewURL(props.projectId, props.taskId) : item.getPreviewURL(props.projectId)));

    const imageSourcesWithIds = computed(() => items.value.filter((item: Attachment) => item.allowImagePreview()).map((item: Attachment) => {
        return ({
            id: item.id,
            url: props.taskId ? item.getDownloadURL(props.projectId, props.taskId) : item.getDownloadURL(props.projectId),
        });
    }));

    const currentImagePreviewIndex = ref<number>(0);

    const showAudioPreviewModal = ref<boolean>(false);

    const audioSources = computed<Attachment[]>(() => items.value.filter((item: Attachment) => item.allowAudioPreview()));

    const currentAudioPreviewIndex = ref<number>(0);

    const showPDFPreviewModal = ref<boolean>(false);

    const pdfSources = computed<Attachment[]>(() => items.value.filter((item: Attachment) => item.allowPDFPreview()));

    const currentPDFPreviewIndex = ref<number>(0);

    const columnDefinitions = reactive<TableHeaderColumn<Attachment>[]>([
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
            render: (row: Attachment) => renderLabel(row.name),
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.size"),
            field: "size",
            visible: true,
            sortable: false,
            isFiltered: () => false,
            render: (row: Attachment) => renderLabel(formatBytes(row.size)),
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.contentType"),
            field: "contentType",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByContentType.value,
            render: (row: Attachment) => renderLabel(row.contentType),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreatedAt.value,
            render: (row: Attachment) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),

        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreator.value,
            render: (row: Attachment) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Attachment>[]>(() =>
        tableSettings.columns.map((column) => { // get saved ordered columns
            const definition = columnDefinitions.find((c) => c.field === column.field);
            return {
                label: definition?.label ?? "",
                field: column.field,
                visible: column.visible,
                sortable: definition!.sortable,
                align: definition?.align,
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );


    const onConfirmDelete = (attachment: Attachment) => {
        dialog.warning({
            title: t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
            content: () =>
                h('div', [
                    t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.message", { name: attachment.name, size: formatBytes(attachment.size) }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(attachment);
            },
        });
    };

    const onDelete = async (attachment: Attachment) => {
        if (props.taskId) {
            onDeleteTaskAttachment(attachment);
        } else {
            onDeleteProjectAttachment(attachment);
        }
    };

    const onDeleteProjectAttachment = async (attachment: Attachment) => {
        if (attachment.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await attachmentService.deleteProjectAttachment(props.projectId, attachment.id);
                notify('success', t("modules.projectAttachment.components.ProjectAttachmentsTab.notifications.projectAttachmentDeleted", { name: attachment.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = attachment;
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
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project attachment id not set", { file: "ProjectAttachmentsTab.vue", method: "onDelete" });
        }
    };

    const onDeleteTaskAttachment = async (projectAttachment: Attachment) => {
        if (projectAttachment.id && props.taskId) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await attachmentService.deleteTaskAttachment(props.projectId, props.taskId, projectAttachment.id);
                notify('success', t("modules.projectAttachment.components.AttachmentsTable.notifications.projectAttachmentDeleted", { name: projectAttachment.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = projectAttachment;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "AttachmentsTable.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "AttachmentsTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project attachment id not set", { file: "AttachmentsTable.vue", method: "onDelete" });
        }
    };

    const refreshContentItemsOptions = () => {
        contentTypeOptions.value = [...new Set(items.value.map((item: Attachment) => { return (item.contentType) }))].map((contentType) => { return ({ label: contentType, value: contentType }); });
    };

    const onRefresh = async () => {
        if (props.projectId) {
            if (props.taskId) {
                onRefreshTaskAttachments();
            } else {
                onRefreshProjectAttachments();
            }
        }
    };

    const onRefreshTaskAttachments = async () => {
        if (props.taskId) {
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
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "AttachmentsTable.onRefresh" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.refreshError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.refreshError");
                        console.error("Unhandled API error", { file: "AttachmentsTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        }
    };

    const onRefreshProjectAttachments = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await attachmentService.getProjectAttachments(props.projectId);
            items.value = results.attachments.map((attachment) => new Attachment(attachment));
            itemCount.value = items.value?.length ?? 0;
            refreshContentItemsOptions();
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
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };


    const onDownload = (attachment: Attachment) => {
        if (props.taskId) {
            bgDownload(attachment.getBgDownloadURL(props.projectId, props.taskId), attachment.name)
        } else {
            bgDownload(attachment.getBgDownloadURL(props.projectId), attachment.name)
        }
    };

    const onPreview = (attachment: Attachment) => {
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
            if (payload.to.includes("ProjectAttachmentsTab.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectAttachmentsTab.onDelete")) {
                onDelete(tmpItem.value);
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
    <AudioPreview v-model:show="showAudioPreviewModal" :project-id="props.projectId" :items="audioSources"
        :current-index="currentAudioPreviewIndex" />
    <PDFPreview v-model:show="showPDFPreviewModal" :project-id="props.projectId" :items="pdfSources"
        :current-index="currentPDFPreviewIndex" />

    <UploadDialog v-if="props.projectId && !props.readOnly" v-model:show="showUploadModal" :project-id="props.projectId"
        v-model:uploadedFilesCount="uploadedFilesCount" />

    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectAttachment.components.projectAttachmentsTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onShowUploadModal" @clear-filters="onClearFilters"
        :buttons="props.readOnly ? ['refresh', 'settings'] : ['refresh', 'add', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.name.placeholder')"
                    v-model:value="filters.name" />
                <n-select size="small" v-else-if="column.field === 'contentType'" clearable
                    :disabled="state.ajaxRunning" :options="contentTypeOptions" v-model:value="filters.contentType"
                    :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.contentType.placeholder')" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="state.ajaxRunning" />
                <UserSelector v-else-if="column.field === 'createdBy'" hideAvatar clearable
                    :disabled="state.ajaxRunning" size="small" v-model:id="filters.createdByUserId"
                    :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.user.placeholder')" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onDownload(row)" :disabled="state.ajaxRunning" class="doneo-table-actions-button">
                    {{ t("shared.buttons.Download.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_DOWNLOAD" />
                    </template>
                </n-button>
                <n-button @click="onPreview(row)"
                    :disabled="state.ajaxRunning || !(row.allowImagePreview() || row.allowAudioPreview() || row.allowPDFPreview())"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Preview.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_PREVIEW" />
                    </template>
                </n-button>
                <n-button @click="onConfirmDelete(row)" :disabled="state.ajaxRunning || props.readOnly"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Delete.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_DELETE" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>