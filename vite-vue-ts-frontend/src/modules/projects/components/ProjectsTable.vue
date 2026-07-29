<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useRouter } from "vue-router";
    import { useI18n } from "vue-i18n";

    import { NModal, NButtonGroup, NButton, NIcon, NDrawer, NDrawerContent, NFlex } from 'naive-ui';

    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_EDIT, DONEO_ICON_ACTION_PREVIEW } from '../../../shared/types/icons.ts';

    import { useLoadingStore } from '../../../stores/loading';
    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { Project } from '../models/project';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { PAGER_DEFAULT_RESULTS_PAGE, type Pagination } from '../../../shared/types/pager.ts';
    import type { Order } from '../../../shared/types/order.ts';
    import { projectService } from '../services/project.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { SearchRequest, ProjectResponse, PatchRequest } from '../types/dto.ts';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';

    import NewProjectForm from './NewProjectForm.vue';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';

    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import ProjectTypeSelector from '../../project-types/components/ProjectTypeSelector.vue';
    import ProjectStatusSelector from '../../project-statuses/components/ProjectStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { renderColoredTag, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';

    import ProjectResumeFloatingCard from './ProjectResumeFloatingCard.vue';
    import ChangeProjectStatusDropdown from '../../../shared/components/dropdowns/ChangeProjectStatusDropdown.vue';
    import type { ProjectStatus } from '../../project-statuses/models/project-status.ts';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "ProjectsTable" });;

    const { t } = useI18n();
    const { notify } = useNotify();
    const router = useRouter();

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

    const items = shallowRef<Project[]>([]);
    const tmpItem = ref<Project>(new Project());

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

    interface ProjectsTableFilters {
        slug: string;
        typeId: string | null;
        priorityId: string | null;
        statusId: string | null;
        summary: string;
        createdAt: TimestampRange;
        createdByUserId: string | null;
    }

    const filters = reactive<ProjectsTableFilters>(
        {
            slug: "",
            typeId: null,
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
            filters.typeId,
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
    const isFilteredByType = computed<boolean>(() => filters.typeId !== null);
    const isFilteredByPriority = computed<boolean>(() => filters.priorityId !== null);
    const isFilteredByStatus = computed<boolean>(() => filters.statusId !== null);
    const isFilteredBySummary = computed<boolean>(() => filters.summary !== "");
    const isFilteredByCreationDate = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByCreator = computed<boolean>(() => filters.createdByUserId !== null);

    const onClearFilters = () => {
        filters.slug = "";
        filters.typeId = null;
        filters.priorityId = null;
        filters.statusId = null;
        filters.summary = "";
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        filters.createdByUserId = null;
    };


    const columnDefinitions = reactive<TableHeaderColumn<Project>[]>([
        {
            label: t("modules.project.components.ProjectsTable.header.columns.slug"),
            field: "slug",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySlug.value,
            render: (row: Project) => renderLabel(row.slug),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.type"),
            field: "type",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByType.value,
            render: (row: Project) => renderColoredTag(row.type.name, row.type.hexColor, true),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.priority"),
            field: "priority",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPriority.value,
            render: (row: Project) => renderColoredTag(row.priority.name, row.priority.hexColor, true),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.status"),
            field: "status",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByStatus.value,
            render: (row: Project) => renderColoredTag(row.status.name, row.status.hexColor, true),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.summary"),
            field: "summary",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySummary.value,
            render: (row: Project) => renderLabel(row.summary),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreationDate.value,
            render: (row: Project) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreator.value,
            render: (row: Project) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Project>[]>(() =>
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
                    slug: filters.slug.length > 0 ? filters.slug : undefined,
                    summary: filters.summary.length > 0 ? filters.summary : undefined,
                    typeId: filters.typeId !== null ? filters.typeId : undefined,
                    priorityId: filters.priorityId !== null ? filters.priorityId : undefined,
                    statusId: filters.statusId !== null ? filters.statusId : undefined,
                    createdAt: filters.createdAt,
                    createdByUserId: filters.createdByUserId !== null ? filters.createdByUserId : undefined,
                }
            };
            const response = await projectService.search(payload);
            currentPagination.enabled = response.pager.enabled;
            currentPagination.currentPage = response.pager.currentPage;
            currentPagination.resultsPage = response.pager.resultsPage;
            currentPagination.totalPages = response.pager.totalPages;
            currentPagination.totalResults = response.pager.totalResults;
            items.value = response.projects.map((project: ProjectResponse) => new Project(project))
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectsTable.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectsTable.vue", method: "onRefresh" }, { err: fatalError });
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
        tmpItem.value = new Project();
        showFormModal.value = true;
    };


    const onProjectAdded = (project: ProjectResponse) => {
        showFormModal.value = false;
        tmpItem.value = new Project();
        notify('success', t("modules.user.components.ProjectsTable.notifications.projectAdded", { name: project.summary }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new Project();
    };

    const showDrawer = ref<boolean>(false);

    const currentProject = ref<Project>(new Project());

    const OpenProject = (project: Project) => {
        router.push(
            {
                name: "project",
                params: {
                    projectId: project.id
                }
            }
        ).catch((e) => {
            console.error(e);
        });
    };
    const onShowProjectResume = (project: Project) => {
        showDrawer.value = true;
        currentProject.value = project;
    };

    let updatedStatus: ProjectStatus;

    const onStatusChanged = async (updatedProject: Project, status: ProjectStatus) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: PatchRequest = {
                id: updatedProject.id ?? "",
                status: {
                    id: status.id ?? ""
                },
            };
            const response: ProjectResponse = await projectService.patch(payload);
            if (response.id === updatedProject.id) {
                onRefresh();
                notify('success', t("modules.project.components.ProjectsTable.notifications.projectStatusUpdated", { summary: updatedProject.summary, status: status.name }));
            } else {
                state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.statusUpdateError", { summary: updatedProject.summary });
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            tmpItem.value = updatedProject;
                            updatedStatus = status;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectsTable.onStatusChanged" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.statusUpdateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectsTable.errors.statusUpdateError");
                    console.error("Unhandled API error", { file: "ProjectsTable.vue", method: "onStatusChanged" }, { err: fatalError });
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
            if (payload.to.includes("ProjectsTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectsTable.onStatusChanged")) {
                onStatusChanged(tmpItem.value, updatedStatus);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="showFormModal">
        <NewProjectForm class="project-form" @add="onProjectAdded" @cancel="hideFormModal" />
    </n-modal>
    <n-drawer v-model:show="showDrawer" :default-width="768" resizable placement="right">
        <n-drawer-content :native-scrollbar="false">
            <template #header>
                <n-flex align="center" justify="space-between">
                    Project {{ currentProject.slug }}
                    <n-icon :component="DONEO_ICON_ACTION_CANCEL" class="doneo-cursor-pointer"
                        @click="showDrawer = false;" />
                </n-flex>
            </template>
            <ProjectResumeFloatingCard v-if="showDrawer && currentProject.id" :project-id="currentProject.id" />
        </n-drawer-content>
    </n-drawer>

    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="items" :row-key="row => row.id"
        :columns="columns" :order="currentOrder" :pager-data="currentPagination" pager-position="both"
        :show-no-items-warning-message="showNoItemsWarningMessage"
        :no-items-warning-message="t('modules.project.components.ProjectsTable.warnings.noItemsFound')" @sort="onSort"
        @refresh="onRefresh" @add="onAdd" @pager-changed="onPagerChanged" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'slug'" clearable :disabled="state.ajaxRunning"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.slug.placeholder')"
                    v-model:value="filters.slug" @keydown-enter="onRefresh" />
                <ProjectTypeSelector v-else-if="column.field === 'type'" :disabled="state.ajaxRunning"
                    v-model:id="filters.typeId" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.type.placeholder')" />
                <ProjectPrioritySelector v-else-if="column.field === 'priority'" :disabled="state.ajaxRunning"
                    size="small" v-model:id="filters.priorityId" :hide-prefix="true" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.priority.placeholder')" />
                <ProjectStatusSelector v-else-if="column.field === 'status'" :disabled="state.ajaxRunning"
                    v-model:id="filters.statusId" :hide-prefix="true" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.status.placeholder')" />
                <TextFilterInput v-else-if="column.field === 'summary'" clearable :disabled="state.ajaxRunning"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                    v-model:value="filters.summary" @keydown-enter="onRefresh" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable ref="createdAtFilterRef"
                    :disabled="state.ajaxRunning" v-model:range="filters.createdAt" />
                <UserSelector v-else-if="column.field === 'createdBy'" clearable :disabled="state.ajaxRunning"
                    v-model:id="filters.createdByUserId"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.creator.placeholder')" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <!-- TODO: use ManageTableActionButtons -->
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button class="doneo-table-actions-button" :disabled="state.ajaxRunning"
                    @click="onShowProjectResume(row)">
                    {{ t("shared.buttons.Preview.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_PREVIEW" />
                    </template>
                </n-button>
                <n-button class="doneo-table-actions-button" :disabled="state.ajaxRunning" @click="OpenProject(row)">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
                <ChangeProjectStatusDropdown :current-status="row.status" :disabled="state.ajaxRunning"
                    @change="(newStatus: ProjectStatus) => onStatusChanged(row, newStatus)" />
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .project-form {
        width: 95%;
        max-width: 640px;
    }
</style>