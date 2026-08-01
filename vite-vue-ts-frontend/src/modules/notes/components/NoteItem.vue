<script setup lang="ts">
    import { ref, h, watch, nextTick, computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NFlex, NButtonGroup, NButton, NIcon, NFormItem, useDialog, NDropdown, type InputInst, type DropdownOption } from 'naive-ui';
    import { IconDeviceFloppy, IconCancel, IconEdit, IconTrash, IconDots } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { renderIcon } from "../../../shared/composables/naive-ui-icon.ts";
    import { Note } from '../models/note';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    import type { NoteItemMode } from "../types/item-mode.ts";

    import ToggleMarkDownEditor from "../../../shared/components/form-blocks/ToggleMarkDownEditor.vue";

    interface NoteItemProps {
        readOnly?: boolean;
        note: Note;
    };

    const props = defineProps<NoteItemProps>();

    const emit = defineEmits(['save', 'delete'])

    const { t } = useI18n();
    const dialog = useDialog();
    const userSettingsStore = useUserSettingsStore();

    const currentMode = ref<NoteItemMode>(!!props.note.id ? "view" : "add");

    const body = ref<string>(props.note.body);

    const onConfirmDelete = () => {
        dialog.warning({
            title: t("modules.note.components.NoteItem.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.note.components.NoteItem.dialogs.deleteConfirmation.message"),
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

    const bodyRef = ref<InputInst | null>(null);

    watch(currentMode, (newValue: NoteItemMode) => {
        if (newValue === "update") {
            nextTick(() => {
                bodyRef.value?.focus();
            });
        }
    });

    watch(() => props.note.updatedAt, (newValue) => {
        if (newValue && props.note.id) {
            currentMode.value = "view";
        }
    });


    const noteOpts = computed(() => [
        {
            label: t("shared.buttons.Update.label"),
            key: "update",
            icon: renderIcon(IconEdit)(18),
        },
        {
            label: t("shared.buttons.Delete.label"),
            key: "delete",
            icon: renderIcon(IconTrash)(18),
        }
    ]);

    const onSelectNoteAction = (key: string | number, _option: DropdownOption) => {
        switch (key) {
            case "update":
                onUpdate();
                break;
            case "delete":
                onConfirmDelete();
                break;
        }
    };
</script>

<template>
    <n-card size="small" bordered>
        <div class="note-header">
            <n-flex justify="space-between" align="center">
                <AvatarUserName :user-id="props.note.createdBy.id" :user-name="props.note.createdBy.name">
                    <template #after>
                        <span>commented on: {{
                            props.note.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask)
                        }}</span>
                        <span v-if="props.note.updatedAt">(comment updated on: {{
                            props.note.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }})
                        </span>
                    </template>
                </AvatarUserName>
                <div v-if="!props.readOnly">
                    <n-dropdown :options="noteOpts" trigger="click" @select="onSelectNoteAction"
                        v-if="currentMode === 'view'">
                        <n-button size="tiny">
                            <template #icon>
                                <n-icon :component="IconDots" />
                            </template>
                        </n-button>
                    </n-dropdown>
                    <n-button-group size="small" class="doneo-note-bottom-action-buttons" v-else>
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
                </div>
            </n-flex>
        </div>
        <n-form-item>
            <ToggleMarkDownEditor v-model:value="body" hide-preview simple-toolbars placeholder="Add note"
                :read-only="currentMode === 'view'" auto-focus />
        </n-form-item>
    </n-card>
</template>


<style lang="css" scoped></style>