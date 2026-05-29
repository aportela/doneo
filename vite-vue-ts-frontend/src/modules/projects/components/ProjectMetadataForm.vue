<script setup lang="ts">
    import { ref, computed, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NForm, NFormItem, NInput, NInputGroup, NButton, NButtonGroup, NIcon, type InputInst, NFlex } from 'naive-ui';

    import type { FormMode } from '../../../shared/types/form-mode';
    import { Project, MAX_KEY_LENGTH, MAX_SUMMARY_LENGTH } from "../models/project";
    import ProjectPrioritySelector from "../../project-priorities/components/ProjectPrioritySelector.vue";
    import ProjectStatusSelector from "../../project-statuses/components/ProjectStatusSelector.vue";
    import ProjectTypeSelector from "../../project-types/components/ProjectTypeSelector.vue";
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { IconX, IconCheck, IconChevronsDown, IconChevronsUp, IconDeviceFloppy } from '@tabler/icons-vue';
    import { useMarkdown } from "../../../shared/composables/useMarkdown.ts";

    interface ProjectFormProps {
        mode: FormMode;
        style?: string | CSSProperties;
        disabled?: boolean;
    }

    const props = defineProps<ProjectFormProps>();

    const emit = defineEmits(["save"]);


    const project = defineModel<Project>("project", { required: true });

    const { t } = useI18n();
    const { render, toMarkdown } = useMarkdown();


    const keyEditMode = ref<boolean>(false);

    const summaryEditMode = ref<boolean>(false);

    const descriptionEditMode = ref<boolean>(false);

    const descriptionExpanded = ref<boolean>(false);

    const htmlMarkDownDescriptionPreview = computed(() => render(project.value.description ?? ""));

    const onSave = () => {
        emit("save");
    };

    const keyRef = ref<InputInst | null>(null);
    const summaryRef = ref<InputInst | null>(null);
    const descriptionRef = ref<InputInst | null>(null);


    const onToggleKeyMode = () => {
        keyEditMode.value = !keyEditMode.value;
        if (keyEditMode) {
            nextTick(() => {
                keyRef.value?.focus();
            });
        }
    };

    const onToggleSummaryMode = () => {
        summaryEditMode.value = !summaryEditMode.value;
        if (summaryEditMode.value) {
            nextTick(() => {
                summaryRef.value?.focus();
            });
        }
    };

    const onToggleDescriptionMode = () => {
        descriptionEditMode.value = !descriptionEditMode.value;
        if (descriptionEditMode.value) {
            nextTick(() => {
                descriptionRef.value?.focus();
            });
        }
    };

    const insertAtCursor = (value: string) => {
        const el = document.activeElement as HTMLTextAreaElement
        if (!el) {
            project.value.description += value
            return
        }

        const start = el.selectionStart ?? project.value.description?.length
        const end = el.selectionEnd ?? project.value.description?.length

        project.value.description =
            project.value.description?.slice(0, start) +
            value +
            project.value.description?.slice(end)

        // restore cursor
        requestAnimationFrame(() => {
            el.selectionStart = el.selectionEnd = start + value.length
        })
    }

    const onPaste = (e: ClipboardEvent) => {
        const clipboard = e.clipboardData
        if (!clipboard) return

        const html = clipboard.getData('text/html')
        const plain = clipboard.getData('text/plain')

        let markdown = plain

        if (html) {
            markdown = toMarkdown(html)
        }

        e.preventDefault()

        insertAtCursor(markdown)
    };

</script>

<template>
    <n-card bordered :style="props.style">
        <n-form-item label="Created by">
            <div class="note-user">
                <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
            </div>
        </n-form-item>
        <n-flex>
            <n-form-item label="Created at">
                {{ project.createdAt.toLocaleString() }}
            </n-form-item>
            <n-form-item label="Updated at">
                {{ project.updatedAt?.toLocaleString() }}
            </n-form-item>
            <n-form-item label="Started at">
                {{ project.updatedAt?.toLocaleString() }}
            </n-form-item>
            <n-form-item label="Finished at">
                {{ project.updatedAt?.toLocaleString() }}
            </n-form-item>
            <n-form-item label="Due at">
                {{ project.updatedAt?.toLocaleString() }}
            </n-form-item>
        </n-flex>
        <n-form>
            <n-flex>
                <n-form-item label="Type">
                    <ProjectTypeSelector v-model:id="project.type.id" :disabled="props.disabled" />
                </n-form-item>
                <n-form-item label="Priority">
                    <ProjectPrioritySelector v-model:id="project.priority.id" :disabled="props.disabled" />
                </n-form-item>
                <n-form-item label="Status">
                    <ProjectStatusSelector v-model:id="project.status.id" :disabled="props.disabled" />
                </n-form-item>
            </n-flex>

            <n-form-item label="Key">
                <n-input-group>
                    <n-input :readonly="!keyEditMode" v-model:value="project.key" :show-count="keyEditMode"
                        :maxlength="MAX_KEY_LENGTH" ref="keyRef" :disabled="props.disabled"
                        @click="() => { if (!keyEditMode) { onToggleKeyMode(); } }" />
                    <n-button v-if="keyEditMode" @click="onToggleKeyMode" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :component="IconCheck" />
                        </template>
                    </n-button>
                </n-input-group>
            </n-form-item>
            <n-form-item label="Summary">
                <n-input-group>
                    <n-input :readonly="!summaryEditMode" v-model:value="project.summary" :show-count="summaryEditMode"
                        :maxlength="MAX_SUMMARY_LENGTH" ref="summaryRef" :disabled="props.disabled"
                        @click="() => { if (!summaryEditMode) { onToggleSummaryMode(); } }" />
                    <n-button v-if="summaryEditMode" @click="onToggleSummaryMode" :disabled="props.disabled">
                        <template #icon>
                            <n-icon :component="IconCheck" />
                        </template>
                    </n-button>
                </n-input-group>
            </n-form-item>
            <n-form-item label="description">
                <template #label>
                    <n-flex align="center">
                        <span>Description</span>
                        <div v-if="!descriptionEditMode">
                            <n-button size="tiny" v-if="!descriptionExpanded"
                                @click="descriptionExpanded = !descriptionExpanded">
                                <template #icon>
                                    <n-icon :component="IconChevronsDown" />
                                </template>
                                expand</n-button>
                            <n-button size="tiny" v-else @click="descriptionExpanded = !descriptionExpanded">
                                <template #icon>
                                    <n-icon :component="IconChevronsUp" />
                                </template>
                                collapse</n-button>
                        </div>
                    </n-flex>
                </template>
                <div v-if="descriptionEditMode" style="width: 100%;">
                    <n-input v-model:value="project.description" type="textarea" clearable :disabled="props.disabled"
                        @paste="onPaste" ref="descriptionRef" :rows="8" />
                    <n-flex justify="end">
                        <n-button-group>
                            <n-button @click="onToggleDescriptionMode" :disabled="props.disabled">
                                <template #icon>
                                    <n-icon :component="IconCheck" />
                                </template>
                            </n-button>
                            <n-button @click="onToggleDescriptionMode" :disabled="props.disabled">
                                <template #icon>
                                    <n-icon :component="IconX" />
                                </template>
                            </n-button>
                        </n-button-group>
                    </n-flex>
                </div>
                <div v-else v-html="htmlMarkDownDescriptionPreview"
                    class="doneo-project-description-markdown-preview doneo-cursor-pointer"
                    :class="{ 'doneo-project-description-markdown-preview-expanded': descriptionExpanded }"
                    @click="onToggleDescriptionMode" />
            </n-form-item>
        </n-form>
        <n-button @click="onSave" :disabled="props.disabled">
            <template #icon>
                <n-icon :component="IconDeviceFloppy" color="red"></n-icon>
            </template>
            {{ t("shared.buttons.Save.label") }}
        </n-button>
    </n-card>
</template>

<style lang="css" scoped>
    .doneo-project-description-markdown-preview {
        width: 100%;
        border: 1px solid #e0e0e6;
        border-radius: var(--n-border-radius);
        padding: 4px 12px;
        color: var(--n-text-color);
        min-height: 1em;

        overflow: hidden;
        max-height: 12em;
        /* aprox 6 líneas */
        transition: max-height 0.3s ease;
    }

    .doneo-project-description-markdown-preview-expanded {
        max-height: unset;
    }
</style>