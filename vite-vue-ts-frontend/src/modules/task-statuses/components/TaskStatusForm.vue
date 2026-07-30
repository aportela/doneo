<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NTooltip } from 'naive-ui';
    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_SAVE, DONEO_ICON_ADD, DONEO_ICON_CLEAR_DATE, DONEO_ICON_EDIT, DONEO_ICON_FILL_DATE, DONEO_ICON_FILL_EMTPY_DATE, DONEO_ICON_NAME, DONEO_ICON_PALETTE, DONEO_ICON_STAR } from '../../../shared/types/icons';

    import { TaskStatus, MAX_NAME_LENGTH } from '../models/task-status';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskStatusService } from '../services/task-status';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/naive-ui-helpers';
    import { generateRandomSoftHexColor } from '../../../shared/composables/color';
    import type { TaskStatusResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';

    interface Props {
        taskStatusId?: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'update', 'cancel'])

    const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const taskStatus = ref<TaskStatus>(new TaskStatus());

    taskStatus.value.hexColor = generateRandomSoftHexColor();

    const formRef = ref<FormInst | null>(null);

    const taskStatusFormRules: FormRules =
    {
        name: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length > MAX_NAME_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
        index: {
            required: true,
            validator: (_rule: FormItemRule, _value: number) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (serverErrors.value.index) {
                    return new Error(t(serverErrors.value.index));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
    };

    watch(() => taskStatus.value.name, () => { delete serverErrors.value.name });
    watch(() => taskStatus.value.index, () => { delete serverErrors.value.index });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !taskStatus.value.name || state.ajaxRunning;
    });

    const onSave = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        try {
            await formRef.value?.validate();
            if (!props.taskStatusId) {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "TaskStatusForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onGet = async (id: string) => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: TaskStatusResponse = await taskStatusService.get(id);
            if (response.id === id) {
                taskStatus.value = new TaskStatus(response);
            } else {
                state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.loadError");
                    console.error("Unhandled API error", { file: "TaskStatusForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const addedTaskStatus: TaskStatusResponse = await taskStatusService.add(taskStatus.value.toAddTaskStatusRequestPayload());
            emit('add', addedTaskStatus)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.taskStatus.components.TaskStatusForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.taskStatus.components.TaskStatusForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.addError");
                    console.error("Unhandled API error", { file: "TaskStatusForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    formRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const updatedTaskStatus: TaskStatusResponse = await taskStatusService.update(taskStatus.value.toUpdateTaskStatusRequestPayload());
            emit('update', updatedTaskStatus)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.taskStatus.components.TaskStatusForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.taskStatus.components.TaskStatusForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusForm.errors.updateError");
                    console.error("Unhandled API error", { file: "TaskStatusForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    formRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskStatusForm.onGet")) {
                if (props.taskStatusId) {
                    onGet(props.taskStatusId);
                }
            } else if (payload.to.includes("TaskStatusForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("TaskStatusForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.taskStatusId) {
            onGet(props.taskStatusId);
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

    const FLAG_ICON_SIZE = 22;
</script>

<template>
    <n-card bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon class="doneo-mr-4px" :component="!props.taskStatusId ? DONEO_ICON_ADD : DONEO_ICON_EDIT" />
                {{
                    t(!props.taskStatusId ? "modules.taskStatus.components.TaskStatusForm.headers.addTaskStatus"
                        : "modules.taskStatus.components.TaskStatusForm.headers.updateTaskStatus")
                }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="formRef" :model="taskStatus" :rules="taskStatusFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.taskStatus.components.TaskStatusForm.inputs.name.placeholder')"
                    v-model:value="taskStatus.name" :maxlength="MAX_NAME_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_NAME" />
                    </template>
                </n-input>
            </n-form-item>
            <n-flex>
                <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.index.label')" path="index"
                    show-feedback>
                    <n-input-number :min="0"
                        :placeholder="t('modules.taskStatus.components.TaskStatusForm.inputs.index.placeholder')"
                        v-model:value="taskStatus.index" :disabled="state.ajaxRunning" required>
                    </n-input-number>
                </n-form-item>
                <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.flags.label')">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_STAR" :size="FLAG_ICON_SIZE" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.defaultStatusOnCreation }"
                                @click="taskStatus.flags.defaultStatusOnCreation = !taskStatus.flags.defaultStatusOnCreation" />
                        </template>
                        {{ t(taskStatus.flags.defaultStatusOnCreation ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasDefaultStatusOnCreation"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotdefaultStatusOnCreation")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_FILL_EMTPY_DATE" :size="FLAG_ICON_SIZE"
                                class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.fillEmptyStartDate }"
                                @click="taskStatus.flags.fillEmptyStartDate = !taskStatus.flags.fillEmptyStartDate" />
                        </template>
                        {{ t(taskStatus.flags.fillEmptyStartDate ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasFillEmptyStartDate"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotFillEmptyStartDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_FILL_DATE" :size="FLAG_ICON_SIZE" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.setStartDate }"
                                @click="taskStatus.flags.setStartDate = !taskStatus.flags.setStartDate" />
                        </template>
                        {{ t(taskStatus.flags.setStartDate ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasSetStartDate"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotSetStartDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_FILL_EMTPY_DATE" :size="FLAG_ICON_SIZE"
                                class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.fillEmptyFinishDate }"
                                @click="taskStatus.flags.fillEmptyFinishDate = !taskStatus.flags.fillEmptyFinishDate" />
                        </template>
                        {{ t(taskStatus.flags.fillEmptyFinishDate ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasFillEmptyFinishDate"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotFillEmptyFinishDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_FILL_DATE" :size="FLAG_ICON_SIZE" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.setFinishDate }"
                                @click="taskStatus.flags.setFinishDate = !taskStatus.flags.setFinishDate" />
                        </template>
                        {{ t(taskStatus.flags.setFinishDate ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasSetFinishDate"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotSetFinishDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="DONEO_ICON_CLEAR_DATE" :size="FLAG_ICON_SIZE" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !taskStatus.flags.unsetFinishDateOnLeave }"
                                @click="taskStatus.flags.unsetFinishDateOnLeave = !taskStatus.flags.unsetFinishDateOnLeave" />
                        </template>
                        {{ t(taskStatus.flags.unsetFinishDateOnLeave ?
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasUnsetFinishDateOnLeave"
                            :
                            "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotUnsetFinishDateOnLeave")
                        }}
                    </n-tooltip>
                </n-form-item>
            </n-flex>
            <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(taskStatus.hexColor)" style="width: 100%;">
                        {{ taskStatus.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="taskStatus.hexColor"
                        :disabled="state.ajaxRunning">
                        <template #trigger="{ onClick, ref: triggerRef }">
                            <n-button :ref="triggerRef" quaternary @click="onClick">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_PALETTE" />
                                </template>
                            </n-button>
                        </template>
                    </n-color-picker>
                </n-flex>
            </n-form-item>
        </n-form>
        <template #action>
            <n-flex>
                <n-button @click="onSave" :disabled="isSaveDisabled">
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_SAVE" />
                    </template>
                    {{ t("shared.buttons.Save.label") }}
                </n-button>
                <n-button @click="onCancel" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_CANCEL" />
                    </template>
                    {{ t("shared.buttons.Cancel.label") }}
                </n-button>
            </n-flex>
        </template>
    </n-card>
</template>

<style lang="css" scoped></style>