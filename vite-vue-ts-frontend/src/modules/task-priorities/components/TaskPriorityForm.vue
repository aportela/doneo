<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_SAVE, DONEO_ICON_ADD, DONEO_ICON_EDIT, DONEO_ICON_NAME, DONEO_ICON_PALETTE } from '../../../shared/types/icons';

    import { TaskPriority, MAX_NAME_LENGTH } from '../models/task-priority';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskPriorityService } from '../services/task-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/naive-ui-helpers';
    import { generateRandomSoftHexColor } from '../../../shared/composables/color';
    import type { TaskPriorityResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';

    interface Props {
        taskPriorityId?: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'update', 'cancel'])

    const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const taskPriority = ref<TaskPriority>(new TaskPriority());

    taskPriority.value.hexColor = generateRandomSoftHexColor();

    const formRef = ref<FormInst | null>(null);

    const taskPriorityFormRules: FormRules =
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

    watch(() => taskPriority.value.name, () => { delete serverErrors.value.name });
    watch(() => taskPriority.value.index, () => { delete serverErrors.value.index });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !taskPriority.value.name || state.ajaxRunning;
    });

    const onSave = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        try {
            await formRef.value?.validate();
            if (!props.taskPriorityId) {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "ProjectPriorityForm.vue", method: "onSave", details: "form validation error", error: error });
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
            const response: TaskPriorityResponse = await taskPriorityService.get(id);
            if (response.id === id) {
                taskPriority.value = new TaskPriority(response);
            } else {
                state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPriorityForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectPriorityForm.vue", method: "onGet" }, { err: fatalError });
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

    const onAdd = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const addedTaskPriority: TaskPriorityResponse = await taskPriorityService.add(taskPriority.value.toAddTaskPriorityRequestPayload());
            emit('add', addedTaskPriority)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPriorityForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.taskPriority.components.TaskPriorityForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.taskPriority.components.TaskPriorityForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.addError");
                    console.error("Unhandled API error", { file: "ProjectPriorityForm.vue", method: "onAdd" }, { err: fatalError });
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
            const updatedTaskPriority: TaskPriorityResponse = await taskPriorityService.update(taskPriority.value.toUpdateTaskPriorityRequestPayload());
            emit('update', updatedTaskPriority)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPriorityForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.taskPriority.components.TaskPriorityForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.taskPriority.components.TaskPriorityForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPriorityForm.errors.updateError");
                    console.error("Unhandled API error", { file: "ProjectPriorityForm.vue", method: "onUpdate" }, { err: fatalError });
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
            if (payload.to.includes("ProjectPriorityForm.onGet")) {
                if (props.taskPriorityId) {
                    onGet(props.taskPriorityId);
                }
            } else if (payload.to.includes("ProjectPriorityForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectPriorityForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.taskPriorityId) {
            onGet(props.taskPriorityId);
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon class="doneo-mr-4px" :component="!props.taskPriorityId ? DONEO_ICON_ADD : DONEO_ICON_EDIT" />
                {{
                    t(!props.taskPriorityId ?
                        "modules.taskPriority.components.TaskPriorityForm.headers.addTaskPriority" :
                        "modules.taskPriority.components.TaskPriorityForm.headers.updateTaskPriority")
                }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="formRef" :model="taskPriority" :rules="taskPriorityFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.taskPriority.components.TaskPriorityForm.inputs.name.placeholder')"
                    v-model:value="taskPriority.name" :maxlength="MAX_NAME_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_NAME" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.index.label')" path="index"
                show-feedback>
                <n-input-number :min="0"
                    :placeholder="t('modules.taskPriority.components.TaskPriorityForm.inputs.index.placeholder')"
                    v-model:value="taskPriority.index" :disabled="state.ajaxRunning">
                </n-input-number>
            </n-form-item>
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(taskPriority.hexColor)" style="width: 100%;">
                        {{ taskPriority.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="taskPriority.hexColor"
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