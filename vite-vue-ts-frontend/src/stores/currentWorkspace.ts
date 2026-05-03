import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../composables/localStorage";

const localStorageCurrentWorkspace = createStorageEntry<string>(
  "currentWorkspace.workspaceId",
  "",
);

interface State {
  currentWorkspaceId: string;
}

export const useCurrentWorkspaceStore = defineStore("workspaceStore", {
  state: (): State => ({
    currentWorkspaceId: localStorageCurrentWorkspace.get(),
  }),
  getters: {
    workspaceId: (state): string => state.currentWorkspaceId,
  },
  actions: {
    set(workspaceId: string): void {
      this.currentWorkspaceId = workspaceId;
      localStorageCurrentWorkspace.set(this.currentWorkspaceId);
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useCurrentWorkspaceStore, import.meta.hot),
  );
}
