<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup } from 'naive-ui';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { ProjectType } from '../models/project-type';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { renderColoredTag, renderIcon } from '../../../shared/composables/naive-ui-helpers.ts';
    import { projectTypeService } from '../services/project-type.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT } from '../../../shared/types/icons.ts';
    import type { ProjectTypeResponse, SearchRequest } from '../types/dto.ts';
    import ProjectTypeForm from './ProjectTypeForm.vue';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "ProjectTypesTable" });;

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

    const items = shallowRef<ProjectType[]>([]);

    const tmpItem = ref<ProjectType>(new ProjectType());

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

    interface ProjectTypesTableFilters {
        name: string;
    }

    const filters = reactive<ProjectTypesTableFilters>(
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
    const localFilteredItems = computed<ProjectType[]>(() => {
        return items.value.filter((projectType: ProjectType) => {
            const name = projectType.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const columnDefinitions = reactive<TableHeaderColumn<ProjectType>[]>([
        {
            label: t("modules.projectType.components.ProjectTypesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: ProjectType) => renderColoredTag(row.name, row.hexColor, true),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<ProjectType>[]>(() =>
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

    const onDelete = async (projectType: ProjectType) => {
        if (projectType.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectTypeService.delete(projectType.id);
                cacheStore.clearProjectTypesCache();
                notify('success', t("modules.projectType.components.ProjectTypesTable.notifications.projectTypeDeleted", { name: projectType.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = projectType;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypesTable.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypesTable.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypesTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypesTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectTypesTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project type id not set", { file: "ProjectTypesTable.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (projectType: ProjectType) => {
        dialog.warning({
            title: t("modules.projectType.components.ProjectTypesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
            content: () =>
                h('div', [
                    t("modules.projectType.components.ProjectTypesTable.dialogs.deleteConfirmation.message", { name: projectType.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(projectType);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        showNoItemsWarningMessage.value = false;
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
            const response = await projectTypeService.search(payload);
            items.value = response.projectTypes.map((projectType: ProjectTypeResponse) => new ProjectType(projectType))
            showNoItemsWarningMessage.value = items.value.length === 0;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypesTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypesTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypesTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectTypesTable.vue", method: "onRefresh" }, { err: fatalError });
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
        tmpItem.value = new ProjectType();
        showFormModal.value = true;
    };

    const onUpdate = (projectType: ProjectType) => {
        tmpItem.value = projectType;
        showFormModal.value = true;
    };

    const onProjectTypeAdded = (projectType: ProjectType) => {
        showFormModal.value = false;
        cacheStore.clearProjectTypesCache();
        notify('success', t("modules.projectType.components.ProjectTypesTable.notifications.projectTypeAdded", { name: projectType.name }));
        onRefresh();
    };

    const onProjectTypeUpdated = (projectType: ProjectType) => {
        showFormModal.value = false;
        cacheStore.clearProjectTypesCache();
        notify('success', t("modules.projectType.components.ProjectTypesTable.notifications.projectTypeUpdated", { name: projectType.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new ProjectType();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectTypesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectTypesTable.onDelete")) {
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
        <ProjectTypeForm class="project-type-form" :project-type-id="tmpItem.id" @add="onProjectTypeAdded"
            @update="onProjectTypeUpdated" @cancel="hideFormModal" v-if="showFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectType.components.ProjectTypesTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.projectType.components.ProjectTypesTable.filters.name.placeholder')"
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
    .project-type-form {
        width: 95%;
        max-width: 640px;
    }
</style>