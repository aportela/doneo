<script setup lang="ts">
    import { ref, h } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NFlex, NButtonGroup, NButton, NIcon, NFormItem, NInput, useDialog } from 'naive-ui';
    import { IconDeviceFloppy, IconCancel, IconEdit, IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from "../../../shared/composables/naive-ui-icon.ts";
    import { Note } from '../models/note';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    import type { NoteItemMode } from "../types/item-mode.ts";
    //import { noteService } from "../services/note.ts";

    interface NoteItemProps {
        note: Note;
    };

    const props = defineProps<NoteItemProps>();

    const emit = defineEmits(['save', 'delete'])

    const { t } = useI18n();
    const dialog = useDialog();

    const currentMode = ref<NoteItemMode>(!!props.note.id ? "view" : "add");
    const body = ref<string>("");

    const onConfirmDelete = () => {
        dialog.warning({
            title: t("modules.projectPriority.components.ProjectPrioritiesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectPriority.components.ProjectPrioritiesTable.dialogs.deleteConfirmation.message", { name: "AAAAAAA" }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", props.note.id);
            },
        });
    };

    const onSave = () => {
        let n = { ...props.note };
        n.body = body.value;
        emit("save", n);
    };

    const onUpdate = () => {
        currentMode.value = "update";
        body.value = props.note.body;
    }

    const onCancel = () => {
        currentMode.value = "view";
    }

</script>

<template>
    <n-card size="small" bordered>
        <div class="note-header">
            <div class="note-user">
                <AvatarUserName :user-id="props.note.user.id" :user-name="props.note.user.name" />
            </div>
            <span class="note-date">
                <div>
                    Created at: {{ props.note.createdAt?.toLocaleString() }}
                </div>
                <div v-if="props.note.updatedAt">
                    Updated at: {{ props.note.updatedAt?.toLocaleString() }}
                </div>
            </span>
        </div>
        <div class="note-content" v-if="currentMode === 'view'">
            {{ props.note.body }}
        </div>
        <n-form-item v-else>
            <n-input placeholder="Type note body" v-model:value="body" type="textarea" rows="6" />
        </n-form-item>

        <n-flex justify="end">
            <n-button-group size="small" class="doneo-note-bottom-action-buttons"
                v-if="['add', 'update'].includes(currentMode)">
                <n-button @click="onSave" :disabled="!body">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy"></n-icon>
                    </template>
                    {{ t("shared.buttons.Save.label") }}
                </n-button>
                <n-button @click="onCancel" v-if="note.id">
                    <template #icon>
                        <n-icon :component="IconCancel"></n-icon>
                    </template>
                    {{ t("shared.buttons.Cancel.label") }}
                </n-button>
            </n-button-group>
            <n-button-group size="small" class="doneo-note-bottom-action-buttons" v-else>
                <n-button @click="onUpdate">
                    <template #icon>
                        <n-icon :component="IconEdit"></n-icon>
                    </template>
                    {{ t("shared.buttons.Update.label") }}
                </n-button>
                <n-button @click="onConfirmDelete">
                    <template #icon>
                        <n-icon :component="IconTrash"></n-icon>
                    </template>
                    {{ t("shared.buttons.Delete.label") }}
                </n-button>
            </n-button-group>
        </n-flex>
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

    .doneo-note-bottom-action-buttons {
        margin-top: 16px;
    }
</style>