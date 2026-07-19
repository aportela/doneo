import { acceptHMRUpdate, defineStore } from "pinia";
import { type SortOrder } from "../shared/types/common";

interface TableColumnConfig {
  field: string;
  visible: boolean;
}

export interface TableConfig {
  sorting: {
    field: string;
    order: SortOrder;
  };
  columns: TableColumnConfig[];
  filters?: Record<string, unknown>;
  pagination: {
    active: boolean;
    page: number;
    pageSize: number;
  };
}

interface State {
  tables: Record<string, TableConfig>;
}

const storePersistenceKey = "doneo.tableSettings";

function createDefaultTableConfig(): TableConfig {
  return {
    sorting: {
      field: "",
      order: "ASC",
    },
    columns: [],
    filters: {},
    pagination: {
      active: true,
      page: 1,
      pageSize: 16,
    },
  };
}

export const useTableSettingsStore = defineStore("tableSettingsStore", {
  persist: {
    key: storePersistenceKey,
  },

  state: (): State => ({
    tables: {},
  }),

  actions: {
    ensure(tableId: string): TableConfig {
      if (!this.tables[tableId]) {
        this.tables[tableId] = createDefaultTableConfig();
      }

      return this.tables[tableId];
    },

    register(tableId: string, defaults: Partial<TableConfig>) {
      if (!this.tables[tableId]) {
        this.tables[tableId] = {
          ...createDefaultTableConfig(),
          ...defaults,
          sorting: {
            ...createDefaultTableConfig().sorting,
            ...defaults.sorting,
          },
          columns: defaults.columns ?? createDefaultTableConfig().columns,
          pagination: {
            ...createDefaultTableConfig().pagination,
            ...defaults.pagination,
          },
          filters: {
            ...defaults.filters,
          },
        };
      }
    },

    get(tableId: string): TableConfig {
      return this.tables[tableId] ?? createDefaultTableConfig();
    },

    reset(tableId: string) {
      delete this.tables[tableId];
    },

    setSorting(tableId: string, sorting: Partial<TableConfig["sorting"]>) {
      const table = this.ensure(tableId);

      table.sorting = {
        ...table.sorting,
        ...sorting,
      };
    },

    addVisibleColumn(tableId: string, field: string) {
      const table = this.ensure(tableId);

      const column = table.columns.find((c) => c.field === field);

      if (column) {
        column.visible = true;
      }
    },

    removeVisibleColumn(tableId: string, field: string) {
      const table = this.ensure(tableId);

      const column = table.columns.find((c) => c.field === field);

      if (column) {
        column.visible = false;
      }
    },

    showAllColumns(tableId: string) {
      const table = this.ensure(tableId);
      table.columns.forEach((column) => (column.visible = true));
    },
    hideAllColumns(tableId: string) {
      const table = this.ensure(tableId);
      table.columns.forEach((column) => (column.visible = false));
    },
    toggleAllColumns(tableId: string) {
      const table = this.ensure(tableId);
      table.columns.forEach((column) => (column.visible = !column.visible));
    },

    toggleVisibleColumn(tableId: string, field: string) {
      const table = this.ensure(tableId);

      const column = table.columns.find((c) => c.field === field);

      if (column) {
        column.visible = !column.visible;
      }
    },

    moveColumn(tableId: string, field: string, direction: "up" | "down") {
      const table = this.ensure(tableId);

      const index = table.columns.findIndex((c) => c.field === field);

      if (index === -1) return;

      const target = direction === "up" ? index - 1 : index + 1;

      if (target < 0 || target >= table.columns.length) return;

      [table.columns[index], table.columns[target]] = [
        table.columns[target],
        table.columns[index],
      ];
    },

    setFilters(tableId: string, filters: Record<string, unknown>) {
      const table = this.ensure(tableId);

      table.filters = filters;
    },

    updateFilters(tableId: string, filters: Record<string, unknown>) {
      const table = this.ensure(tableId);

      table.filters = {
        ...table.filters,
        ...filters,
      };
    },

    clearFilters(tableId: string) {
      const table = this.ensure(tableId);

      table.filters = {};
    },

    setPagination(
      tableId: string,
      pagination: Partial<TableConfig["pagination"]>,
    ) {
      const table = this.ensure(tableId);

      table.pagination = {
        ...table.pagination,
        ...pagination,
      };
    },

    setPage(tableId: string, page: number) {
      this.ensure(tableId).pagination.page = page;
    },

    setPageSize(tableId: string, pageSize: number) {
      this.ensure(tableId).pagination.pageSize = pageSize;
    },

    replace(tableId: string, config: Partial<TableConfig>) {
      const table = this.ensure(tableId);

      this.tables[tableId] = {
        ...table,
        ...config,
        sorting: {
          ...table.sorting,
          ...config.sorting,
        },
        columns: config.columns ?? table.columns,
        pagination: {
          ...table.pagination,
          ...config.pagination,
        },
        filters: {
          ...table.filters,
          ...config.filters,
        },
      };
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useTableSettingsStore, import.meta.hot),
  );
}
