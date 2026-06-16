<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NTooltip } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette, IconStar, IconCalendarBolt, IconCalendarCancel, IconCalendarMinus } from '@tabler/icons-vue';

    import { TaskStatus, MAX_NAME_LENGTH } from '../models/task-status';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskStatusService } from '../services/task-status';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { TaskStatusResponse, AddRequest, UpdateRequest } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface TaskStatusFormProps {
        mode: FormMode;
        taskStatusId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<TaskStatusFormProps>();

    const { t } = useI18n();

    const taskStatus = ref<TaskStatus>(new TaskStatus());
    taskStatus.value.hexColor = generateRandomSoftHexColor();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const taskStatusFormRef = ref<FormInst | null>(null)

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
        return !taskStatus.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        taskStatusFormRef.value?.restoreValidation();
        try {
            await taskStatusFormRef.value?.validate();
            if (props.mode === "add") {
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
        taskStatusFormRef.value?.restoreValidation();
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
        taskStatusFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: taskStatus.value.name ?? "",
                hexColor: taskStatus.value.hexColor ?? "",
                index: taskStatus.value.index ?? 0,
                flags: {
                    defaultStatusOnCreation: taskStatus.value.flags?.defaultStatusOnCreation ?? false,
                    fillEmptyStartDate: taskStatus.value.flags?.fillEmptyStartDate ?? false,
                    setStartDate: taskStatus.value.flags?.setStartDate ?? false,
                    fillEmptyFinishDate: taskStatus.value.flags?.fillEmptyFinishDate ?? false,
                    setFinishDate: taskStatus.value.flags?.setFinishDate ?? false,
                    unsetFinishDateOnLeave: taskStatus.value.flags?.unsetFinishDateOnLeave ?? false,
                }
            };
            const addedRole: TaskStatusResponse = await taskStatusService.add(payload);
            emit('add', addedRole)
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
                    taskStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        taskStatusFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: taskStatus.value.id ?? "",
                name: taskStatus.value.name ?? "",
                hexColor: taskStatus.value.hexColor ?? "",
                index: taskStatus.value.index ?? 0,
                flags: {
                    defaultStatusOnCreation: taskStatus.value.flags?.defaultStatusOnCreation ?? false,
                    fillEmptyStartDate: taskStatus.value.flags?.fillEmptyStartDate ?? false,
                    setStartDate: taskStatus.value.flags?.setStartDate ?? false,
                    fillEmptyFinishDate: taskStatus.value.flags?.fillEmptyFinishDate ?? false,
                    setFinishDate: taskStatus.value.flags?.setFinishDate ?? false,
                    unsetFinishDateOnLeave: taskStatus.value.flags?.unsetFinishDateOnLeave ?? false,
                }
            };
            const updatedRole: TaskStatusResponse = await taskStatusService.update(payload);
            emit('update', updatedRole)
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
                    taskStatusFormRef.value?.validate().then(() => { }).catch(() => { });
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
                } else {
                    console.error(`TODO: missing taskStatusId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("TaskStatusForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("TaskStatusForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.taskStatusId) {
                onGet(props.taskStatusId);
            } else {
                console.error(`TODO: missing taskStatusId property for ${props.mode} action`);
            }
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

    const flagIconSize = 22;
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon :component="props.mode == 'add' ? IconPlus : IconEdit" />
                {{ t(props.mode == "add" ? "modules.taskStatus.components.TaskStatusForm.headers.addTaskStatus"
                    : "modules.taskStatus.components.TaskStatusForm.headers.updateTaskStatus") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="taskStatusFormRef" :model="taskStatus" :rules="taskStatusFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.taskStatus.components.TaskStatusForm.inputs.name.placeholder')"
                    v-model:value="taskStatus.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-flex>
                <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.index.label')" path="index"
                    show-feedback>
                    <n-input-number :min="0"
                        :placeholder="t('modules.taskStatus.components.TaskStatusForm.inputs.index.placeholder')"
                        v-model:value="taskStatus.index" required>
                    </n-input-number>
                </n-form-item>
                <n-form-item :label="t('modules.taskStatus.components.TaskStatusForm.inputs.flags.label')">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconStar" :size="flagIconSize" class="doneo-cursor-help"
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
                            <n-icon :component="IconCalendarBolt" :size="flagIconSize" class="doneo-cursor-help"
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
                            <n-icon :component="IconCalendarCancel" :size="flagIconSize" class="doneo-cursor-help"
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
                            <n-icon :component="IconCalendarBolt" :size="flagIconSize" class="doneo-cursor-help"
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
                            <n-icon :component="IconCalendarCancel" :size="flagIconSize" class="doneo-cursor-help"
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
                            <n-icon :component="IconCalendarMinus" :size="flagIconSize" class="doneo-cursor-help"
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
                    <n-tag :color="getNaiveUITagColorProperty(taskStatus.hexColor ?? '#888888')" style="width: 100%;">
                        {{ taskStatus.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="taskStatus.hexColor">
                        <template #trigger="{ onClick, ref: triggerRef }">
                            <n-button :ref="triggerRef" quaternary @click="onClick">
                                <template #icon>
                                    <IconPalette />
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
            </n-flex>
        </template>
    </n-card>
</template>

<style lang="css" scoped></style>