import { defineStore, acceptHMRUpdate } from "pinia";
import type { UserBaseResponse } from "../modules/users/types/dto";
import type { TaskStatusResponse } from "../modules/task-statuses/types/dto";
import type { TaskPriorityResponse } from "../modules/task-priorities/types/dto";
import type { ProjectStatusResponse } from "../modules/project-statuses/types/dto";
import type { ProjectPriorityResponse } from "../modules/project-priorities/types/dto";
import type { ProjectTypeResponse } from "../modules/project-types/types/dto";

interface State {
  usersCache: UserBaseResponse[];
  projectStatusesCache: ProjectStatusResponse[];
  projectPrioritiesCache: ProjectPriorityResponse[];
  projectTypesCache: ProjectTypeResponse[];
  taskStatusesCache: TaskStatusResponse[];
  taskPrioritiesCache: TaskPriorityResponse[];
}

export const useCacheStore = defineStore("cacheStore", {
  state: (): State => ({
    usersCache: [],
    projectStatusesCache: [],
    projectPrioritiesCache: [],
    projectTypesCache: [],
    taskStatusesCache: [],
    taskPrioritiesCache: [],
  }),

  getters: {
    users: (state) => state.usersCache,
    projectStatuses: (state) => state.projectStatusesCache,
    projectPriorities: (state) => state.projectPrioritiesCache,
    projectTypes: (state) => state.projectTypesCache,
    taskStatuses: (state) => state.taskStatusesCache,
    taskPriorities: (state) => state.taskPrioritiesCache,
  },

  actions: {
    clearAllCaches() {
      this.clearUsersCache;
      this.clearProjectStatusesCache();
      this.clearProjectPrioritiesCache();
      this.clearProjectTypesCache();
      this.clearTaskStatusesCache();
      this.clearTaskPrioritiesCache();
    },
    setUsersCache(users: UserBaseResponse[]) {
      this.usersCache = users;
    },
    clearUsersCache() {
      this.usersCache.length = 0;
    },
    setProjectStatusesCache(projectStatuses: ProjectStatusResponse[]) {
      this.projectStatusesCache = projectStatuses;
    },
    clearProjectStatusesCache() {
      this.projectStatuses.length = 0;
    },
    setProjectPrioritiesCache(projectPriorities: ProjectPriorityResponse[]) {
      this.projectPrioritiesCache = projectPriorities;
    },
    clearProjectPrioritiesCache() {
      this.projectPrioritiesCache.length = 0;
    },
    setProjectTypesCache(projectTypes: ProjectTypeResponse[]) {
      this.projectTypesCache = projectTypes;
    },
    clearProjectTypesCache() {
      this.projectTypesCache.length = 0;
    },
    setTaskStatusesCache(taskStatuses: TaskStatusResponse[]) {
      this.taskStatusesCache = taskStatuses;
    },
    clearTaskStatusesCache() {
      this.taskStatusesCache.length = 0;
    },
    setTaskPrioritiesCache(taskPriorities: TaskPriorityResponse[]) {
      this.taskPrioritiesCache = taskPriorities;
    },
    clearTaskPrioritiesCache() {
      this.taskPrioritiesCache.length = 0;
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCacheStore, import.meta.hot));
}
