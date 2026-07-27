<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { RouterLink } from "vue-router";
    import { useI18n } from "vue-i18n";

    import { NModal, NButtonGroup, NButton, NIcon, NDrawer, NDrawerContent, NFlex } from 'naive-ui';
    import { DONEO_ICON_ACTION_OPEN, DONEO_ICON_MAXIMIZE } from '../../../shared/types/icons.ts';

    import { useLoadingStore } from '../../../stores/loading';
    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { Task } from '../models/tasks.ts';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { PAGER_DEFAULT_RESULTS_PAGE, type Pagination } from '../../../shared/types/pager.ts';
    import type { Order } from '../../../shared/types/order.ts';
    import { taskService } from '../services/task.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { PatchRequest, SearchRequest, TaskResponse } from '../types/dto.ts';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';

    import NewTaskForm from './NewTaskForm.vue';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';

    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import TaskPrioritySelector from '../../task-priorities/components/TaskPrioritySelector.vue';
    import TaskStatusSelector from '../../task-statuses/components/TaskStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import type { TaskStatus } from '../../task-statuses/models/task-status.ts';
    import TaskResumeFloatingCard from './TaskResumeFloatingCard.vue';
    import { renderColoredTag, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';

    interface Props {
        id?: string;
        projectId?: string;
        readOnly?: boolean;
    }

    const props = withDefaults(defineProps<Props>(), { id: "TasksTable" });;

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const userSettingsStore = useUserSettingsStore();
    const tableSettingsStore = useTableSettingsStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<Task[]>([]);
    const tmpItem = ref<Task>(new Task());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        onRefresh();
    };

    const currentPagination = reactive<Pagination>({ enabled: true, currentPage: 1, resultsPage: PAGER_DEFAULT_RESULTS_PAGE, totalPages: 1, totalResults: 0 });
    const resetPager = ref<boolean>(false);

    const onPagerChanged = (newPagination: Pagination) => {
        currentPagination.enabled = newPagination.enabled;
        currentPagination.currentPage = newPagination.currentPage;
        currentPagination.resultsPage = newPagination.resultsPage;
        onRefresh();
    };

    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface TasksTableFilters {
        slug: string;
        priorityId: string | null;
        statusId: string | null;
        summary: string;
        createdAt: TimestampRange;
        createdByUserId: string | null;
    }

    const filters = reactive<TasksTableFilters>(
        {
            slug: "",
            priorityId: null,
            statusId: null,
            summary: "",
            createdAt: {
                from: null,
                to: null,
            },
            createdByUserId: null,
        }
    );

    watch(
        () => [
            filters.slug,
            filters.priorityId,
            filters.statusId,
            filters.summary,
            filters.createdAt.from,
            filters.createdAt.to,
            filters.createdByUserId,
        ],
        () => {
            resetPager.value = true;
        },
    );

    const isFilteredBySlug = computed<boolean>(() => filters.slug !== "");
    const isFilteredByPriority = computed<boolean>(() => filters.priorityId !== null);
    const isFilteredByStatus = computed<boolean>(() => filters.statusId !== null);
    const isFilteredBySummary = computed<boolean>(() => filters.summary !== "");
    const isFilteredByCreationDate = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByCreator = computed<boolean>(() => filters.createdByUserId !== null);

    const onClearFilters = () => {
        filters.slug = "";
        filters.priorityId = null;
        filters.statusId = null;
        filters.summary = "";
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        filters.createdByUserId = null;
    };


    const columnDefinitions = reactive<TableHeaderColumn<Task>[]>([
        {
            label: "Slug",
            field: "slug",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredBySlug.value,
            render: (row: Task) => renderLabel(row.slug),
        },
        {
            label: "Priority",
            field: "priority",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByPriority.value,
            render: (row: Task) => renderColoredTag(row.priority.name, row.priority.hexColor, true),
        },
        {
            label: "Status",
            field: "status",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByStatus.value,
            render: (row: Task) => renderColoredTag(row.status.name, row.status.hexColor, true),
        },
        {
            label: "Summary",
            field: "summary",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredBySummary.value,
            render: (row: Task) => renderLabel(row.summary),
        },
        {
            label: "Created at",
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreationDate.value,
            render: (row: Task) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
        },
        {
            label: "Created by",
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreator.value,
            render: (row: Task) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Task>[]>(() =>
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

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        showNoItemsWarningMessage.value = false;
        try {
            const payload: SearchRequest = {
                pager: { ...currentPagination, currentPage: resetPager.value ? 1 : currentPagination.currentPage },
                order: currentOrder,
                filter: {
                    projectId: props.projectId,
                    slug: filters.slug !== "" ? filters.slug : undefined,
                    summary: filters.summary !== "" ? filters.summary : undefined,
                    priorityId: filters.priorityId !== null ? filters.priorityId : undefined,
                    statusId: filters.statusId !== null ? filters.statusId : undefined,
                    createdAt: filters.createdAt,
                    createdByUserId: filters.createdByUserId !== null ? filters.createdByUserId : undefined,
                }
            };
            const response = await taskService.search(null, payload);
            currentPagination.enabled = response.pager.enabled;
            currentPagination.currentPage = response.pager.currentPage;
            currentPagination.resultsPage = response.pager.resultsPage;
            currentPagination.totalPages = response.pager.totalPages;
            currentPagination.totalResults = response.pager.totalResults;
            items.value = response.tasks.map((task: TaskResponse) => new Task(task))
            resetPager.value = false;
            showNoItemsWarningMessage.value = items.value.length === 0;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTasksPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageTasksPage.vue", method: "onRefresh" }, { err: fatalError });
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
        tmpItem.value = new Task();
        showFormModal.value = true;
    };


    const onTaskAdded = (task: TaskResponse) => {
        showFormModal.value = false;
        tmpItem.value = new Task();
        notify('success', t("modules.user.components.ProjectsTable.notifications.projectAdded", { name: task.summary }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new Task();
    };

    const showDrawer = ref<boolean>(false);

    const currentTask = ref<Task>(new Task());

    const onShowTaskResume = (task: Task) => {
        showDrawer.value = true;
        currentTask.value = task;
    };


    let updatedStatusTask: Task;
    let updatedStatus: TaskStatus;

    const onStatusChanged = async (updatedTask: Task, status: TaskStatus) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: PatchRequest = {
                id: updatedTask.id ?? "",
                status: {
                    id: status.id ?? ""
                },
            };
            const response: TaskResponse = await taskService.patch(updatedTask.projectId ?? '', payload);
            if (response.id === updatedTask.id) {
                onRefresh();
                notify('success', t("modules.task.components.ManageTasksPage.notifications.taskStatusUpdated", { summary: updatedTask.summary, status: status.name }));
            } else {
                state.ajaxErrorMessage = t("modules.task.components.ManageTasksPage.errors.statusUpdateError", { summary: updatedTask.summary });
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            updatedStatusTask = updatedTask;
                            updatedStatus = status;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTasksPage.onStatusChanged" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.task.components.ManageTasksPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.ManageTasksPage.errors.statusUpdateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.ManageTasksPage.errors.statusUpdateError");
                    console.error("Unhandled API error", { file: "ManageTasksPage.vue", method: "onStatusChanged" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageTasksPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageTasksPage.onStatusChanged")) {
                onStatusChanged(updatedStatusTask, updatedStatus);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="props.projectId && showFormModal">
        <NewTaskForm class="task-form" :project-id="props.projectId" @add="onTaskAdded" @cancel="hideFormModal" />
    </n-modal>
    <n-drawer v-model:show="showDrawer" :default-width="768" resizable placement="right">
        <n-drawer-content :native-scrollbar="false">
            <template #header>
                <n-flex align="center" justify="space-between">
                    Task {{ currentTask.slug }}
                    <router-link
                        :to="{ name: 'task', params: { projectId: currentTask.projectId, taskId: currentTask.id } }">
                        <n-icon :component="DONEO_ICON_MAXIMIZE" />
                    </router-link>
                </n-flex>
            </template>
            <TaskResumeFloatingCard v-if="showDrawer && currentTask.projectId && currentTask.id"
                :project-id="currentTask.projectId" :task-id="currentTask.id" />
        </n-drawer-content>
    </n-drawer>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning || props.readOnly" :rows="items"
        :row-key="row => row.id" :columns="columns" :order="currentOrder" :pager-data="currentPagination"
        pager-position="both" :show-no-items-warning-message="showNoItemsWarningMessage"
        :no-items-warning-message="t('modules.task.components.TasksTable.warnings.noItemsFound')" @sort="onSort"
        @refresh="onRefresh" @add="onAdd" @pager-changed="onPagerChanged" @clear-filters="onClearFilters"
        :buttons="props.projectId && !props.readOnly ? ['refresh', 'add', 'settings'] : ['refresh', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'slug'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.task.components.TasksTable.header.filters.slug.placeholder')"
                    v-model:value="filters.slug" @keydown-enter="onRefresh" />
                <TaskPrioritySelector v-else-if="column.field === 'priority'" :disabled="state.ajaxRunning" size="small"
                    v-model:id="filters.priorityId" :hide-prefix="true" clearable
                    :placeholder="t('modules.task.components.TasksTable.header.filters.priority.placeholder')" />
                <TaskStatusSelector v-else-if="column.field === 'status'" :disabled="state.ajaxRunning" size="small"
                    v-model:id="filters.statusId" :hide-prefix="true" clearable
                    :placeholder="t('modules.task.components.TasksTable.header.filters.status.placeholder')" />
                <TextFilterInput v-else-if="column.field === 'summary'" clearable :disabled="state.ajaxRunning"
                    size="small"
                    :placeholder="t('modules.task.components.TasksTable.header.filters.summary.placeholder')"
                    v-model:value="filters.summary" @keydown-enter="onRefresh" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable ref="createdAtFilterRef"
                    :disabled="state.ajaxRunning" size="small" v-model:range="filters.createdAt" />
                <UserSelector v-else-if="column.field === 'createdBy'" hideAvatar clearable
                    :disabled="state.ajaxRunning" size="small" v-model:id="filters.createdByUserId"
                    :placeholder="t('modules.task.components.TasksTable.header.filters.creator.placeholder')" />
            </th>
        </template>
        <!--
        <template #tbody>
            <tr v-for="task, index in items" :key="task.id ?? index">
                <td>
                    {{ task.slug }}
                </td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(task.priority.hexColor ?? '#888888')">{{
                    task.priority.name
                        }}</n-tag></td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(task.status.hexColor ?? '#888888')">{{
                    task.status.name }}</n-tag></td>
                <td><router-link
                        :to="{ name: 'taskTab', params: { taskId: task.id, projectId: task.projectId, tab: 'metadata' } }"
                        class="doneo-link-text-color-default">{{
                            task.summary
                        }}</router-link></td>
                <td>{{ task.createdAt.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</td>
                <td>
                    <AvatarUserName :user-id="task.createdBy.id" :user-name="task.createdBy.name" />
                </td>
                <td class="doneo-text-center">
                    <n-button-group :size="DEFAULT_BUTTON_SIZE">
                        <n-button :disabled="state.ajaxRunning" :size="DEFAULT_BUTTON_SIZE"
                            @click="onShowTaskResume(task)">
                            {{ t("shared.buttons.Open.label") }}
                            <template #icon>
                                <n-icon :size="22" :component="IconFilePencil" />
                            </template>
                        </n-button>
                        <ChangeTaskStatusDropdown :disabled="state.ajaxRunning" :read-only="props.readOnly"
                            :current-status="task.status"
                            @change="(status: TaskStatus) => onStatusChange(task, status)" />
                    </n-button-group>
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !state.ajaxRunning">
                    <n-empty :description="t('modules.task.components.TasksTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
        -->
        <template #rowactions="{ row }">
            <!-- TODO: use ManageTableActionButtons -->
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button class="doneo-table-actions-button" :disabled="state.ajaxRunning"
                    @click="onShowTaskResume(row)">
                    {{ t("shared.buttons.Open.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_OPEN" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>