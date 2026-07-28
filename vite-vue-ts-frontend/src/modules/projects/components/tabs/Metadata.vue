<script setup lang="ts">
    import { ref, reactive, watch, computed, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NForm, NFormItem, NInput, NButton, NButtonGroup, NIcon, NFlex } from 'naive-ui';

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
    import { IconDeviceFloppy, IconCancel } from '@tabler/icons-vue';
    import ToggleDateTimePicker from '../../../../shared/components/form-blocks/ToggleDateTimePicker.vue';
    import { IDate } from '../../../../shared/types/idate.ts';
    import ToggleMarkDownEditor from '../../../../shared/components/form-blocks/ToggleMarkDownEditor.vue';

    interface Props {
        readOnly?: boolean;
        style?: string | CSSProperties;
        disabled?: boolean;
    }

    const props = defineProps<Props>();

    const project = defineModel<Project>("project", { required: true });

    const emit = defineEmits(["save", "refresh"]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const { t } = useI18n();
    const loadingStore = useLoadingStore();

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

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
                currentMode.value = "view";
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

    const onFillEmptyStartDate = () => {
        if (!project.value.startedAt.hasValue()) {
            project.value.startedAt = new IDate(Date.now())
        }
    };

    const onSetStartDate = () => {
        project.value.startedAt = new IDate(Date.now())
    };

    const onFillEmptyFinishDate = () => {
        if (!project.value.finishedAt.hasValue()) {
            project.value.finishedAt = new IDate(Date.now())
        }
    };

    const onSetFinishDate = () => {
        if (!project.value.finishedAt.hasValue()) {
            project.value.finishedAt = new IDate(Date.now())
        }
    };

    const onUnsetFinishDateOnLeave = () => {
        if (project.value.finishedAt.hasValue()) {
            project.value.finishedAt.clear();
        }
    };


    const currentMode = ref<"view" | "edit">("view");

    const readOnlyMode = computed(() => currentMode.value === 'view');

    const onRefresh = () => {
        currentMode.value = "view";
        emit("refresh");
    };
</script>

<template>
    <!-- TODO: add missing i18n labels -->
    <n-card bordered :style="props.style">
        <n-button @click=" currentMode = 'edit'" secondary>Update form</n-button>
        <n-button @click=" currentMode = 'edit'" secondary v-if="project.archivedAt.hasValue()">UnArchive</n-button>
        <n-flex align=" center" justify="space-between">
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
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly && !readOnlyMode" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.startedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Finished at">
                <ToggleDateTimePicker clearable v-model:value="project.finishedAt.msTimestamp"
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly && !readOnlyMode" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.finishedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Due at">
                <ToggleDateTimePicker clearable v-model:value="project.dueAt.msTimestamp"
                    :disabled="props.disabled || state.ajaxRunning" v-if="!props.readOnly && !readOnlyMode" />
                <span class="doneo-datetime-label-readonly" v-else>
                    {{ project.dueAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Archived at">
                <span class="doneo-datetime-label-readonly">
                    {{ project.archivedAt?.toLocaleString() }}
                </span>
            </n-form-item>
        </n-flex>
        <n-form>
            <n-flex>
                <n-form-item label="Slug">
                    <n-input v-model:value="project.slug" :show-count="!(props.readOnly || readOnlyMode)"
                        :maxlength="MAX_SLUG_LENGTH" :disabled="props.disabled || state.ajaxRunning"
                        :read-only="props.readOnly || readOnlyMode" />
                </n-form-item>
                <n-form-item label="Type">
                    <ProjectTypeSelector v-model:id="project.type.id" :disabled="props.disabled || state.ajaxRunning"
                        :read-only="props.readOnly || readOnlyMode" />
                </n-form-item>
                <n-form-item label="Priority">
                    <ProjectPrioritySelector v-model:id="project.priority.id"
                        :disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly || readOnlyMode" />
                </n-form-item>
                <n-form-item label="Status">
                    <ProjectStatusSelector v-model:id="project.status.id"
                        ::disabled="props.disabled || state.ajaxRunning" :read-only="props.readOnly || readOnlyMode"
                        @fill-empty-start-date="onFillEmptyStartDate" @set-start-date="onSetStartDate"
                        @fill-empty-finish-date="onFillEmptyFinishDate" @set-finish-date="onSetFinishDate"
                        @unset-finish-date-on-leave="onUnsetFinishDateOnLeave" />
                </n-form-item>
            </n-flex>
            <n-form-item label="Summary">
                <n-input v-model:value="project.summary" :show-count="!(props.readOnly || readOnlyMode)"
                    :maxlength="MAX_SUMMARY_LENGTH" :disabled="props.disabled || state.ajaxRunning"
                    :read-only="props.readOnly || readOnlyMode" />
            </n-form-item>
            <n-form-item label="description">
                <template #label>
                    <n-flex align="center">
                        <span>Description</span>
                    </n-flex>
                </template>
                <ToggleMarkDownEditor v-if="true" :read-only="props.readOnly || readOnlyMode"
                    v-model:value="project.description" />
            </n-form-item>
        </n-form>
        <n-button-group v-if="!props.readOnly && !readOnlyMode && project.allowedOperations.updateProject">
            <n-button @click="onUpdate" :disabled="props.disabled">
                <template #icon>
                    <n-icon :component="IconDeviceFloppy"></n-icon>
                </template>
                {{ t("shared.buttons.Save.label") }}
            </n-button>
            <n-button @click="onRefresh" :disabled="props.disabled">
                <template #icon>
                    <n-icon :component="IconCancel"></n-icon>
                </template>
                {{ t("shared.buttons.Cancel.label") }}
            </n-button>
        </n-button-group>
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