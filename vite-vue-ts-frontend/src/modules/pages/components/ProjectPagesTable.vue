<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";
    import { useRouter } from 'vue-router';

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';

    import { appBus } from '../../../shared/composables/bus';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { Page } from '../models/page.ts';


    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { pageService } from '../services/page.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';

    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import { renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';
    import type { SearchResponse } from '../types/dto.ts';
    import { DONEO_ICON_ACTION_EDIT } from '../../../shared/types/icons.ts';

    interface Props {
        id?: string;
        readOnly?: boolean;
        projectId: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "ProjectPagesTable" });;

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const { t } = useI18n();
    const router = useRouter();
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

    const items = shallowRef<Page[]>([]);

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "title", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                a.title.localeCompare(b.title)
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                b.title.localeCompare(a.title)
            );
        }
    };


    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);
    const updatedAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface ProjectPagesTableFilters {
        title: string;
        userId: string | null;
        createdAt: TimestampRange;
        updatedAt: TimestampRange
        ;
    }

    const filters = reactive<ProjectPagesTableFilters>(
        {
            title: "",
            userId: null,
            createdAt: {
                from: null,
                to: null,
            },
            updatedAt: {
                from: null,
                to: null,
            },
        }
    );


    const isFilteredByTitle = computed<boolean>(() => filters.title !== "");
    const isFilteredByUser = computed<boolean>(() => filters.userId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByUpdatedAt = computed<boolean>(() => filters.updatedAt.from != null || filters.updatedAt.to != null);

    const localFilteredItems = computed(() => {
        return items.value.filter((page: Page) => {
            const title = page.title?.toLowerCase();
            return (
                (!title || title?.includes(titleFilterLowerCase.value)) &&
                (filters.userId === null || filters.userId === page.createdBy.id) &&
                ((filters.createdAt.from === null && filters.createdAt.to === null) || (page.createdAt.msTimestamp != null && filters.createdAt.from != null && filters.createdAt.from <= page.createdAt.msTimestamp && filters.createdAt.to != null && filters.createdAt.to >= page.createdAt.msTimestamp))
            );
        });
    });

    const onClearFilters = () => {
        filters.title = "";
        filters.userId = null;
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        if (updatedAtFilterRef.value) {
            updatedAtFilterRef.value[0]?.reset();
        }
    };

    const titleFilterLowerCase = computed(() =>
        filters.title.toLowerCase()
    );

    const columnDefinitions = reactive<TableHeaderColumn<Page>[]>([
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.title"),
            field: "title",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByTitle.value,
            render: (row: Page) => renderLabel(row.title),
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreatedAt.value,
            render: (row: Page) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUser.value,
            render: (row: Page) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.updatedAt"),
            field: "updatedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUpdatedAt.value,
            render: (row: Page) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Page>[]>(() =>
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
        try {
            const results: SearchResponse = await pageService.getProjectPages(props.projectId);
            items.value = results.pages.map((page) => new Page(page));
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
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };


    const onOpen = (projectPageId: string | null) => {
        router.push(
            { name: "projectPage", params: { projectId: props.projectId, pageId: projectPageId } }
        ).catch((e) => {
            console.error(e);
        });
    };

    const onAdd = () => {

    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectAttachmentsTab.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters"
        :buttons="props.readOnly ? ['refresh', 'settings'] : ['refresh', 'add', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'title'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.projectPage.components.ProjectPagesTable.filters.title.placeholder')"
                    v-model:value="filters.title" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="state.ajaxRunning" />
                <UserSelector v-else-if="column.field === 'createdBy'" v-model:id="filters.userId"
                    :disabled="state.ajaxRunning" size="small" hide-avatar clearable
                    :placeholder="t('modules.projectPage.components.ProjectPagesTable.filters.user.placeholder')" />
                <DateFilterSelect v-else-if="column.field === 'updatedAt'" clearable v-model:range="filters.updatedAt"
                    ref="updatedAtFilterRef" :disabled="state.ajaxRunning" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onOpen(row.id)" :disabled="state.ajaxRunning" class="doneo-table-actions-button">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>