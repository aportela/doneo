<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NSwitch } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconPlus } from '@tabler/icons-vue';

    import { ProjectTask, MAX_SUMMARY_LENGTH } from '../models/tasks';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTaskService } from '../services/task.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { TaskResponse, AddRequest } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import TaskPrioritySelector from '../../task-priorities/components/TaskPrioritySelector.vue';
    import TaskStatusSelector from '../../task-statuses/components/TaskStatusSelector.vue';

    interface NewTaskFormProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const emit = defineEmits(['add', 'cancel'])

    const props = defineProps<NewTaskFormProps>();

    const { t } = useI18n();

    const openTaskAfterCreate = ref<boolean>(true);

    const task = ref<ProjectTask>(new ProjectTask());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const newTaskFormRef = ref<FormInst | null>(null)

    const newTaskFormRules: FormRules =
    {
        summary: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length > MAX_SUMMARY_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
        priority: {
            id: {
                required: true,
                validator: (_rule: FormItemRule, value: string) => {
                    if (state.ajaxRunning) {
                        return true;
                    }
                    if (!value?.trim()) {
                        return new Error(t("shared.warningMessages.fieldIsRequired"));
                    } else {
                        return true;
                    }
                },
                trigger: ['blur'],
            }
        },
        status: {
            id: {
                required: true,
                validator: (_rule: FormItemRule, value: string) => {
                    if (state.ajaxRunning) {
                        return true;
                    }
                    if (!value?.trim()) {
                        return new Error(t("shared.warningMessages.fieldIsRequired"));
                    } else {
                        return true;
                    }
                },
                trigger: ['blur'],
            }
        }

    };

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !task.value.summary || !task.value.priority.id || !task.value.status.id;
    });

    const onSave = async () => {
        serverErrors.value = {};
        newTaskFormRef.value?.restoreValidation();
        try {
            await newTaskFormRef.value?.validate();
            await onAdd();
        }
        catch (error: any) {
            console.warn("Warning", { file: "NewTaskForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }


    const onAdd = async () => {
        serverErrors.value = {};
        newTaskFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                summary: task.value.summary ?? "",
                description: task.value.description,
                priority: { id: task.value.priority.id ?? "" },
                status: { id: task.value.status.id ?? "" },
            };
            const addedTask: TaskResponse = await projectTaskService.add(props.projectId, payload);
            emit('add', addedTask, openTaskAfterCreate.value)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "NewTaskForm.onAdd" } });
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.task.components.NewTaskForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.task.components.NewTaskForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.NewTaskForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.NewTaskForm.errors.addError");
                    console.error("Unhandled API error", { file: "NewTaskForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    newTaskFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("NewTaskForm.onAdd")) {
                onAdd();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon :component="IconPlus" />
                {{ t("modules.task.components.NewTaskForm.headers.addTask") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="taskTypeFormRef" :model="task" :rules="newTaskFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.task.components.NewTaskForm.inputs.summary.label')" path="summary"
                show-feedback>
                <n-input type="text" :placeholder="t('modules.task.components.NewTaskForm.inputs.summary.placeholder')"
                    v-model:value="task.summary" :maxlength="MAX_SUMMARY_LENGTH" :show-count="true" clearable required>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.task.components.NewTaskForm.inputs.description.label')" path="description"
                show-feedback>
                <n-input type="textarea"
                    :placeholder="t('modules.task.components.NewTaskForm.inputs.description.placeholder')"
                    v-model:value="task.description" clearable>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.task.components.NewTaskForm.selectors.taskPriority.label')"
                path="priority.id">
                <TaskPrioritySelector required v-model:id="task.priority.id"
                    :placeholder="t('modules.task.components.NewTaskForm.selectors.taskPriority.placeholder')" />
            </n-form-item>
            <n-form-item :label="t('modules.task.components.NewTaskForm.selectors.taskStatus.label')" path="status.id">
                <TaskStatusSelector required v-model:id="task.status.id"
                    :placeholder="t('modules.task.components.NewTaskForm.selectors.taskStatus.placeholder')" />
            </n-form-item>
        </n-form>
        <template #action>
            <n-flex>
                <n-button @click="onSave" :disabled="isSaveDisabled">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy" />
                    </template>
                    {{ t("shared.buttons.Save.label") }}
                </n-button>
                <n-button @click="onCancel" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconCancel" />
                    </template>
                    {{ t("shared.buttons.Cancel.label") }}
                </n-button>
                <n-switch size="large" class="doneo-open-task-after-create-switch" v-model:value="openTaskAfterCreate">
                    <template #checked>
                        {{
                            t('modules.task.components.NewTaskForm.selectors.switches.openTaskAfterCreate.label')
                        }}
                    </template>
                    <template #unchecked>
                        {{
                            t('modules.task.components.NewTaskForm.selectors.switches.openTaskAfterCreate.label')
                        }}
                    </template>
                </n-switch>
            </n-flex>
        </template>
    </n-card>

</template>

<style lang="css" scoped>
    .doneo-open-task-after-create-switch {
        margin-top: 4px;
    }
</style>