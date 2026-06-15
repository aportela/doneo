<script setup lang="ts">
    import { shallowRef, reactive, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NSpace, NCard, NButtonGroup, NButton } from "naive-ui";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../../stores/loading';
    import { useNotify } from '../../../../shared/composables/notification';
    import { appBus } from '../../../../shared/composables/bus';
    import { Note } from "../../../notes/models/note.ts";
    import type { AddRequest, UpdateRequest } from "../../../notes/types/dto.ts";
    import { noteService } from "../../../notes/services/note.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';
    import type { SearchResponse } from "../../../notes/types/dto.ts";
    import NoteItem from "../../../notes/components/NoteItem.vue";
    import { useSessionStore } from "../../../../stores/session.ts";

    interface ProjectNotesProps {
        style?: string | CSSProperties;
        projectId: string;
        taskId: string;
    }

    const props = defineProps<ProjectNotesProps>();

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();
    const sesionStore = useSessionStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Note[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.taskId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onRefresh = async () => {
        if (props.taskId) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const results: SearchResponse = await noteService.getTaskNotes(props.projectId, props.taskId);
                items.value = results.notes.map((note) => new Note(note));
                itemCount.value = items.value?.length ?? 0;
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectNotes.onRefresh" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                        console.error("Unhandled API error", { file: "ProjectNotes.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("project id not set", { file: "ProjectNotes.vue", method: "onRefresh" });
        }
    };

    const onAddNote = () => {
        items.value = [new Note({
            id: "",
            user: {
                id: sesionStore.sessionUserId ?? "",
                name: sessionStorage.sessionUserName ?? "",
            },
            createdAt: new Date().getTime(),
            updatedAt: null,
            body: "",
        }),
        ...items.value
        ];
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        if (props.taskId) {
            onRefresh();
        }
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectNotes.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectNotes.onDelete")) {
                //onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

    const onSaveNote = async (note: Note) => {
        if (props.taskId) {
            try {
                if (!note.id) {
                    const payload: AddRequest = {
                        body: note.body
                    };
                    note = new Note(await noteService.addTaskNote(props.projectId, props.taskId, payload));
                    items.value = [note, ...items.value.filter((note) => !!note.id)]
                    itemCount.value = items.value?.length ?? 0;
                    notify('success', t("modules.note.components.ProjectNotesTab.notifications.projectNoteAdded"));
                } else if (note.id) {
                    const payload: UpdateRequest = {
                        id: note.id,
                        user: {
                            id: note.user.id ?? "",
                            name: note.user.name ?? "",
                        },
                        createdAt: note.createdAt?.msTimestamp ?? 0,
                        updatedAt: null,
                        body: note.body
                    };
                    note = new Note(await noteService.updateTaskote(props.projectId, props.taskId, note.id, payload));
                    items.value = items.value.map((item) => item.id === note.id ? note : item)
                    itemCount.value = items.value?.length ?? 0;
                    notify('success', t("modules.note.components.ProjectNotesTab.notifications.projectNoteUpdated"));
                    // TODO: this will remove pending notes, do not allow add more than 1 one without saving
                }
            } catch { }
        } else {
            console.error("project id not set", { file: "ProjectNotes.vue", method: "onSaveNote" });
        }
    }

    const onDeleteNote = async (id: string) => {
        if (props.taskId) {
            try {
                await noteService.deleteTaskNote(props.projectId, props.taskId, id);
                items.value = items.value.filter((note) => note.id != id)
                itemCount.value = items.value?.length ?? 0;
                notify('success', t("modules.note.components.ProjectNotesTab.notifications.projectNoteDeleted"));
            } catch { }
        } else {
            console.error("project id not set", { file: "ProjectNotes.vue", method: "onDeleteNote" });
        }
    }
</script>

<template>
    <n-card bordered :style="props.style">
        <n-button-group style="margin-bottom: 16px;">
            <n-button @click="onAddNote">Add Note</n-button>
            <n-button @click="onRefresh">Refresh notes</n-button>
        </n-button-group>
        <n-space vertical size="large" style="margin-right: 12px;">
            <NoteItem v-for="note, index in items" :key="note.id ?? index" :note="note" @save="onSaveNote"
                @delete="onDeleteNote" />
        </n-space>
    </n-card>
</template>

<style lang="css" scoped></style>