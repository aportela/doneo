<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h, type Component } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup, NTooltip } from 'naive-ui';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { ProjectStatus } from '../models/project-status';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { renderColoredTag, renderIcon, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import { projectStatusService } from '../services/project-status.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT, DONEO_ICON_CLEAR_DATE, DONEO_ICON_FILL_DATE, DONEO_ICON_FILL_EMTPY_DATE, DONEO_ICON_STAR } from '../../../shared/types/icons.ts';
    import type { ProjectStatusResponse, SearchRequest } from '../types/dto.ts';
    import ProjectStatusForm from './ProjectStatusForm.vue';

    interface Props {
        id?: string;
    };

    const props = withDefaults(defineProps<Props>(), { id: "ProjectStatusesTable" });;

    const { t } = useI18n();
    const dialog = useDialog();

    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const tableSettingsStore = useTableSettingsStore();
    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<ProjectStatus[]>([]);

    const tmpItem = ref<ProjectStatus>(new ProjectStatus());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "name" ? a.name.localeCompare(b.name) : String(a.index).localeCompare(String(b.index))
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "name" ? b.name.localeCompare(a.name) : String(b.index).localeCompare(String(a.index))
            );
        }
    };

    interface ProjectStatusTableFilters {
        name: string;
    };

    const filters = reactive<ProjectStatusTableFilters>(
        {
            name: "",
        }
    );

    const isFilteredByName = computed<boolean>(() => filters.name.length > 0);

    const onClearFilters = () => {
        filters.name = "";
    };

    const nameFilterLowerCase = computed(() => filters.name.toLowerCase());

    // we have all results, use local filtering for avoiding server load
    const localFilteredItems = computed<ProjectStatus[]>(() => {
        return items.value.filter((projectType: ProjectStatus) => {
            const name = projectType.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const FLAG_ICON_SIZE = 22;

    const createFlagTooltip = (
        enabled: boolean,
        icon: Component,
        enabledKey: string,
        disabledKey: string
    ) =>
        h(
            NTooltip,
            { trigger: "hover" },
            {
                trigger: () =>
                    h(NIcon, {
                        component: icon,
                        size: FLAG_ICON_SIZE,
                        class: [
                            "doneo-cursor-help",
                            { "doneo-disabled-icon": !enabled },
                        ],
                    }),
                default: () => t(enabled ? enabledKey : disabledKey),
            }
        );

    // get project status column from project status row
    const renderProjectStatusFlagsColumn = (row: ProjectStatus) => {
        return h("div",
            {
                class: "project-status-flags",
            },
            [
                createFlagTooltip(
                    row.flags.defaultStatusOnCreation,
                    DONEO_ICON_STAR,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasDefaultStatusOnCreation",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotdefaultStatusOnCreation"
                ),
                createFlagTooltip(
                    row.flags.fillEmptyStartDate,
                    DONEO_ICON_FILL_EMTPY_DATE,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasFillEmptyStartDate",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotFillEmptyStartDate"
                ),
                createFlagTooltip(
                    row.flags.setStartDate,
                    DONEO_ICON_FILL_DATE,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasSetStartDate",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotSetStartDate"
                ),
                createFlagTooltip(
                    row.flags.fillEmptyFinishDate,
                    DONEO_ICON_FILL_EMTPY_DATE,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasFillEmptyFinishDate",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotFillEmptyFinishDate"
                ),
                createFlagTooltip(
                    row.flags.setFinishDate,
                    DONEO_ICON_FILL_DATE,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasSetFinishDate",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotSetFinishDate"
                ),
                createFlagTooltip(
                    row.flags.unsetFinishDateOnLeave,
                    DONEO_ICON_CLEAR_DATE,
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasUnsetFinishDateOnLeave",
                    "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotUnsetFinishDateOnLeave"
                ),
            ]
        );
    };

    const columnDefinitions = reactive<TableHeaderColumn<ProjectStatus>[]>([
        {
            label: t("modules.projectStatus.components.ProjectStatusesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: ProjectStatus) => renderColoredTag(row.name, row.hexColor, true),
        },
        {
            label: t("modules.projectStatus.components.ProjectStatusesTable.header.columns.index"),
            field: "index",
            visible: true,
            sortable: true,
            isFiltered: () => false,
            render: (row: ProjectStatus) => renderLabel(row.index),
        },
        {
            label: t("modules.projectStatus.components.ProjectStatusesTable.header.columns.flags"),
            field: "flags",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => false,
            render: (row: ProjectStatus) => renderProjectStatusFlagsColumn(row),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<ProjectStatus>[]>(() =>
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

    const onDelete = async (projectStatus: ProjectStatus) => {
        if (projectStatus.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectStatusService.delete(projectStatus.id);
                cacheStore.clearProjectStatusesCache();
                notify('success', t("modules.projectStatus.components.ProjectStatusesTable.notifications.projectStatusUpdated", { name: projectStatus.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = projectStatus;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusesTable.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.notFoundError");
                                break;
                            case 409:
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.deleteUsedError", { name: projectStatus.name });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectStatusesTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project status id not set", { file: "ProjectStatusesTable.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (projectStatus: ProjectStatus) => {
        dialog.warning({
            title: t("modules.projectStatus.components.ProjectStatusesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
            content: () =>
                h('div', [
                    t("modules.projectStatus.components.ProjectStatusesTable.dialogs.deleteConfirmation.message", { name: projectStatus.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(projectStatus);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    enabled: false,
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: currentOrder.field,
                    direction: currentOrder.direction,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                }
            };
            const response = await projectStatusService.search(payload);
            items.value = response.projectStatuses.map((projectStatus: ProjectStatusResponse) => new ProjectStatus(projectStatus))
            showNoItemsWarningMessage.value = items.value.length === 0;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusesTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusesTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectStatusesTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpItem.value = new ProjectStatus();
        showFormModal.value = true;
    };

    const onUpdate = (projectStatus: ProjectStatus) => {
        tmpItem.value = projectStatus;
        showFormModal.value = true;
    };

    const onProjectStatusAdded = (projectStatus: ProjectStatus) => {
        showFormModal.value = false;
        cacheStore.clearProjectStatusesCache();
        notify('success', t("modules.projectStatus.components.ProjectStatusesTable.notifications.projectStatusAdded", { name: projectStatus.name }));
        onRefresh();
    };

    const onProjectStatusUpdated = (projectStatus: ProjectStatus) => {
        showFormModal.value = false;
        cacheStore.clearProjectStatusesCache();
        notify('success', t("modules.projectStatus.components.ProjectStatusesTable.notifications.projectStatusUpdated", { name: projectStatus.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new ProjectStatus();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectStatusesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectStatusesTable.onDelete")) {
                onDelete(tmpItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showFormModal">
        <ProjectStatusForm class="project-status-form" :project-status-id="tmpItem.id" @add="onProjectStatusAdded"
            @update="onProjectStatusUpdated" @cancel="hideFormModal" v-if="showFormModal" />
        <span v-else />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectStatus.components.ProjectStatusesTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.projectStatus.components.ProjectStatusesTable.filters.name.placeholder')"
                    v-model:value="filters.name" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onUpdate(row)" :disabled="state.ajaxRunning" class="doneo-table-actions-button">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
                <n-button @click="onConfirmDelete(row)" :disabled="state.ajaxRunning"
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

<style lang="css" scoped>
    .project-status-form {
        width: 95%;
        max-width: 640px;
    }
</style>