<script setup lang="ts">
    import { ref, reactive, watch, computed, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NForm, NFormItem, NInput, NButton, NButtonGroup, NIcon, type InputInst, NFlex, NEllipsis } from 'naive-ui';

    import { useLoadingStore } from '../../../../stores/loading';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import { projectService } from '../../services/project';
    import { handleAPIError } from '../../../../api/client/errorHandler';
    import { appBus } from '../../../../shared/composables/bus';
    import type { ProjectResponse, UpdateRequest } from '../../types/dto';

    import { Project, MAX_SLUG_LENGTH, MAX_SUMMARY_LENGTH } from "../../models/project";
    import ProjectPrioritySelector from "../../../project-priorities/components/ProjectPrioritySelector.vue";
    import ProjectStatusSelector from "../../../project-statuses/components/ProjectStatusSelector.vue";
    import ProjectTypeSelector from "../../../project-types/components/ProjectTypeSelector.vue";
    import AvatarUserName from '../../../../shared/components/AvatarUserName.vue';
    import { IconX, IconCheck, IconDeviceFloppy } from '@tabler/icons-vue';
    import { useMarkdown } from "../../../../shared/composables/useMarkdown.ts";
    import ToggleInput from '../../../../shared/components/ToggleInput.vue';
    import ToggleDateTimePicker from '../../../../shared/components/ToggleDateTimePicker.vue';
    import { IDate } from '../../../../shared/types/idate.ts';

    interface ProjectFormProps {
        readOnly?: boolean;
        style?: string | CSSProperties;
        disabled?: boolean;
    }

    const project = defineModel<Project>("project", { required: true });

    const props = defineProps<ProjectFormProps>();

    const emit = defineEmits(["save"]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const { t } = useI18n();
    const loadingStore = useLoadingStore();
    const { render, toMarkdown } = useMarkdown();

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    interface ToggleInputComponent {
        setEditMode: () => void
        setViewMode: () => void
    };

    const slugRef = ref<ToggleInputComponent | undefined>();

    const descriptionEditMode = ref<boolean>(false);

    const descriptionExpanded = ref<boolean>(true);

    const htmlMarkDownDescriptionPreview = computed(() => render(project.value.description ?? ""));

    const onUpdate = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: project.value.id ?? "",
                slug: project.value.slug ?? "",
                summary: project.value.summary ?? "",
                description: project.value.description,
                type: {
                    id: project.value.type.id ?? ""
                },
                priority: {
                    id: project.value.priority.id ?? ""
                },
                status: {
                    id: project.value.status.id ?? ""
                },
                startedAt: project.value.startedAt?.msTimestamp ?? null,
                finishedAt: project.value.finishedAt?.msTimestamp ?? null,
                dueAt: project.value.dueAt?.msTimestamp ?? null,
            };
            const response: ProjectResponse = await projectService.update(payload);
            if (response.id === project.value.id) {
                project.value = new Project(response);
            } else {
                state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPage.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
                    console.error("Unhandled API error", { file: "ProjectPage.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }

        }
    };

    const descriptionRef = ref<InputInst | null>(null);

    const onConfirmNewSlugValue = (newValue: string | null) => {
        if (project.value.slug != newValue) {
            project.value.slug = newValue;
        }
        slugRef.value?.setViewMode();
    };

    const onCancelNewSlugValue = () => {
        slugRef.value?.setViewMode();
    };

    const onToggleDescriptionMode = () => {
        if (!props.readOnly && project.value.allowedOperations.updateProject) {
            descriptionEditMode.value = !descriptionEditMode.value;
            if (descriptionEditMode.value) {
                nextTick(() => {
                    descriptionRef.value?.focus();
                });
            }
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

    const onFillEmptyStartDate = () => {
        if (!project.value.startedAt.hasValue()) {
            project.value.startedAt = new IDate(new Date().getTime())
        }
    };

    const onSetStartDate = () => {
        project.value.startedAt = new IDate(new Date().getTime())
    };

    const onFillEmptyFinishDate = () => {
        if (!project.value.finishedAt.hasValue()) {
            project.value.finishedAt = new IDate(new Date().getTime())
        }
    };

    const onSetFinishDate = () => {
        if (!project.value.finishedAt.hasValue()) {
            project.value.finishedAt = new IDate(new Date().getTime())
        }
    };

    const onUnsetFinishDateOnLeave = () => {
        if (project.value.finishedAt.hasValue()) {
            project.value.finishedAt.clear();
        }
    };

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
    <!-- TODO: add missing i18n labels -->
    <n-card bordered :style="props.style">
        <n-flex align="center" justify="space-between">
            <n-form-item label="Created by">
                <div class="note-user">
                    <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
                </div>
            </n-form-item>
            <div>
                <div>Created at: {{ project.createdAt.toLocaleString() }}</div>
                <div v-if="project.updatedAt.hasValue()">Updated at: {{ project.updatedAt?.toLocaleString() }}</div>
            </div>
        </n-flex>
        <n-flex>
            <n-form-item label="Created at">
                <span class="doneo-datetime-label-readonly">
                    {{ project.createdAt.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Updated at">
                <span class="doneo-datetime-label-readonly">
                    {{ project.updatedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Started at">
                <ToggleDateTimePicker clearable v-model:value="project.startedAt.msTimestamp"
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.startedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Finished at">
                <ToggleDateTimePicker clearable v-model:value="project.finishedAt.msTimestamp"
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.finishedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Due at">
                <ToggleDateTimePicker clearable v-model:value="project.dueAt.msTimestamp"
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.dueAt?.toLocaleString() }}
                </span>
            </n-form-item>
        </n-flex>
        <n-form>
            <n-flex>
                <n-form-item label="Slug">
                    <ToggleInput v-model:value="project.slug" show-count :max-length="MAX_SLUG_LENGTH"
                        :disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly"
                        v-on:confirm="onConfirmNewSlugValue" v-on:cancel="onCancelNewSlugValue" ref="slugRef" />
                </n-form-item>
                <n-form-item label="Type">
                    <ProjectTypeSelector v-model:id="project.type.id" :disabled="props.disabled || state.ajaxRunning"
                        :read-only="props.readOnly" />
                </n-form-item>
                <n-form-item label="Priority">
                    <ProjectPrioritySelector v-model:id="project.priority.id"
                        :disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly" />
                </n-form-item>
                <n-form-item label="Status">
                    <ProjectStatusSelector v-model:id="project.status.id"
                        ::disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly"
                        @fill-empty-start-date="onFillEmptyStartDate" @set-start-date="onSetStartDate"
                        @fill-empty-finish-date="onFillEmptyFinishDate" @set-finish-date="onSetFinishDate"
                        @unset-finish-date-on-leave="onUnsetFinishDateOnLeave" />
                </n-form-item>
            </n-flex>
            <n-form-item label="Summary">
                <ToggleInput v-model:value="project.summary" show-count :max-length="MAX_SUMMARY_LENGTH"
                    :disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly" />
            </n-form-item>
            <n-form-item label="description">
                <template #label>
                    <n-flex align="center">
                        <span>Description</span>
                    </n-flex>
                </template>
                <div v-if="descriptionEditMode" style="width: 100%;">
                    <n-input v-model:value="project.description" type="textarea" clearable
                        :disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly" @paste="onPaste"
                        ref="descriptionRef" :rows="8" />
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
                <div v-else v-html="htmlMarkDownDescriptionPreview" class="doneo-project-description-markdown-preview"
                    :class="{ 'doneo-project-description-markdown-preview-expanded': descriptionExpanded, 'doneo-cursor-pointer': !props.readOnly }"
                    @click="onToggleDescriptionMode" />
                <!-- TODO: test alternatives -->
                <n-ellipsis v-if="false" expand-trigger="click" line-clamp="4" :tooltip="false" class="ellipsis"
                    v-html="htmlMarkDownDescriptionPreview">
                </n-ellipsis>
            </n-form-item>
        </n-form>
        <n-button @click="onUpdate" :disabled="props.disabled"
            v-if="!props.readOnly && project.allowedOperations.updateProject">
            <template #icon>
                <n-icon :component="IconDeviceFloppy"></n-icon>
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
        min-height: 1.5em;
        overflow: hidden;
        max-height: 12em;
        transition: max-height 0.3s ease;
    }

    .doneo-project-description-markdown-preview-expanded {
        max-height: unset;
    }
</style>