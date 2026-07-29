<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup } from 'naive-ui';

    import { IconTrash } from '@tabler/icons-vue';

    import { useLoadingStore } from '../../../stores/loading';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { TimeTracking } from '../models/time-tracking.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { timeTrackingService } from '../services/time-tracking.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';

    import TimeTrackingForm from './TimeTrackingForm.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';
    import { renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import type { SearchResponse } from '../types/dto.ts';
    import { DONEO_ICON_ACTION_DELETE } from '../../../shared/types/icons.ts';

    interface Props {
        id?: string;
        readOnly?: boolean;
        projectId: string;
        taskId: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "TimeTrackingsTable" });

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

    const items = shallowRef<TimeTracking[]>([]);

    const tmpItem = ref<TimeTracking>(new TimeTracking());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "createdAt", direction: "DESC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        onRefresh();
    };

    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface TimeTrackingsTableFilters {
        createdByUserId: string | null;
        createdAt: TimestampRange;
        summary: string;
        // TODO: spentTime
    }


    const filters = reactive<TimeTrackingsTableFilters>(
        {
            createdByUserId: null,
            createdAt: {
                from: null,
                to: null,
            },
            summary: "",
        }
    );


    const isFilteredByCreator = computed<boolean>(() => filters.createdByUserId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredBySummary = computed<boolean>(() => filters.summary !== "");

    const onClearFilters = () => {
        filters.createdByUserId = null;
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        filters.summary = "";
    };

    const columnDefinitions = reactive<TableHeaderColumn<TimeTracking>[]>([
        {
            label: t("modules.timeTracking.components.TimeTrackingsTable.header.columns.summary"),
            field: "summary",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredBySummary.value,
            render: (row: TimeTracking) => renderLabel(row.summary),
        },
        {
            label: t("modules.timeTracking.components.TimeTrackingsTable.header.columns.spentTime"),
            field: "spentTime",
            visible: true,
            sortable: false,
            isFiltered: () => false,
            render: (row: TimeTracking) => renderLabel(row.geti18nTimeParts().map(({ key, count }) => `${count} ${t(key, count)}`).join(", ")),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreatedAt.value,
            render: (row: TimeTracking) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreator.value,
            render: (row: TimeTracking) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<TimeTracking>[]>(() =>
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

    const onDelete = async (timeTracking: TimeTracking) => {
        if (timeTracking.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await timeTrackingService.deleteTaskTimeTracking(props.projectId, props.taskId, timeTracking.id);
                items.value = items.value.filter((item) => item.id != timeTracking.id)
                itemCount.value = items.value?.length ?? 0;
                notify('success', t("modules.timeTracking.components.taskTimeTrackingsTab.notifications.timeTrackingDeleted", { summary: timeTracking.summary }));
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = timeTracking;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPermissions.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.deleteError");
                        console.error("Unhandled API error", { file: "TimeTrackings.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "TimeTrackings.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (timeTracking: TimeTracking) => {
        dialog.warning({
            title: t("modules.timeTracking.components.TimeTrackingsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.timeTracking.components.TimeTrackingsTable.dialogs.deleteConfirmation.message", { summary: timeTracking.summary }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(timeTracking);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await timeTrackingService.getTaskTimeTrackings(props.projectId, props.taskId);
            items.value = results.timeTrackings.map((timeTracking) => new TimeTracking(timeTracking));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TrackTimeTrackingsTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.TimeTrackingsTab.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.TimeTrackingsTab.errors.refreshError");
                    console.error("Unhandled API error", { file: "TrackTimeTrackingsTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpItem.value = new TimeTracking();
        showFormModal.value = true;
    };

    const onTaskTimeTrackingAdded = (timeTracking: TimeTracking) => {
        showFormModal.value = false;
        notify('success', t("modules.timeTracking.components.taskTimeTrackingsTab.notifications.timeTrackingAdded", { summary: timeTracking.summary }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new TimeTracking();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TrackTimeTrackingsTab.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="showFormModal">
        <TimeTrackingForm :project-id="props.projectId" :task-id="props.taskId" mode="add" style="width: 42%;"
            @add="onTaskTimeTrackingAdded" @cancel="hideFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="items" :row-key="row => row.id"
        :columns="columns" :order="currentOrder" :show-no-items-warning-message="showNoItemsWarningMessage"
        :no-items-warning-message="t('modules.timeTracking.components.TimeTrackingsTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters"
        :buttons="props.readOnly ? ['refresh', 'settings'] : ['refresh', 'add', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'summary'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.timeTracking.components.TimeTrackingsTable.filters.summary.placeholder')"
                    v-model:value="filters.summary" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="state.ajaxRunning" />
                <UserSelector v-else-if="column.field === 'createdBy'" clearable :disabled="state.ajaxRunning"
                    size="small" v-model:id="filters.createdByUserId"
                    :placeholder="t('modules.timeTracking.components.TimeTrackingsTable.filters.user.placeholder')" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
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