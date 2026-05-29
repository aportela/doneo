<script setup lang="ts">
    import { shallowRef, reactive, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NSpace, NCard, NButton } from "naive-ui";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    //    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { Note } from "../../notes/models/note.ts";
    import type { AddRequest, UpdateRequest } from "../../notes/types/dto.ts";
    import { noteService } from "../../notes/services/note.ts";
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { SearchResponse } from "../../notes/types/dto.ts";
    import NoteItem from "../../notes/components/NoteItem.vue";

    interface ProjectNotesProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectNotesProps>();

    const { t } = useI18n();
    //const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Note[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onRefresh = async () => {
        if (props.projectId) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const results: SearchResponse = await noteService.getProjectNotes(props.projectId);
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
                id: "019e6db9-6f4a-702f-94cd-7dd645a4d03a",
                name: "John doe",
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
        if (props.projectId) {
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
        if (props.projectId) {
            try {
                if (!note.id) {
                    const payload: AddRequest = {
                        body: note.body
                    };
                    await noteService.addProjectNote(props.projectId, payload);
                    onRefresh();
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
                    await noteService.updateProjectNote(props.projectId, note.id, payload);
                    onRefresh();
                }
            } catch { }
        } else {
            console.error("project id not set", { file: "ProjectNotes.vue", method: "onSaveNote" });
        }
    }

    const onDeleteNote = async (id: string) => {
        if (props.projectId) {
            try {
                const response = await noteService.deleteProjectNote(props.projectId, id);
                console.log(response);
                onRefresh();
            } catch { }
        } else {
            console.error("project id not set", { file: "ProjectNotes.vue", method: "onDeleteNote" });
        }
    }
</script>

<template>
    <n-card bordered :style="props.style">
        <n-button @click="onAddNote" style="margin-bottom: 16px;">Add Note</n-button>

        <n-space vertical size="large" style="margin-right: 12px;">
            <NoteItem v-for="note, index in items" :key="note.id ?? index" :note="note" @save="onSaveNote"
                @delete="onDeleteNote" />
        </n-space>
    </n-card>
</template>

<style lang="css" scoped>
    .note-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
    }

    .note-user {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 500;
    }

    .note-date {
        font-size: 12px;
        color: #999;
    }

    .note-content {
        font-size: 14px;
        white-space: pre-line;
    }
</style>