<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette } from '@tabler/icons-vue';

    import { TaskPriority, MAX_NAME_LENGTH } from '../models/task-priority';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskPriorityService } from '../services/task-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { TaskPriorityResponse, AddRequest, UpdateRequest } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectStatusFormProps {
        mode: FormMode;
        taskPriorityId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectStatusFormProps>();

    const { t } = useI18n();

    const taskPriority = ref<TaskPriority>(new TaskPriority());
    taskPriority.value.hexColor = generateRandomSoftHexColor();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const taskPriorityFormRef = ref<FormInst | null>(null)

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
        return !taskPriority.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        taskPriorityFormRef.value?.restoreValidation();
        try {
            await taskPriorityFormRef.value?.validate();
            if (props.mode === "add") {
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
        taskPriorityFormRef.value?.restoreValidation();
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
                    taskPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        taskPriorityFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: taskPriority.value.name ?? "",
                hexColor: taskPriority.value.hexColor ?? "",
                index: taskPriority.value.index ?? 0,
            };
            const addedRole: TaskPriorityResponse = await taskPriorityService.add(payload);
            emit('add', addedRole)
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
                    taskPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        taskPriorityFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: taskPriority.value.id ?? "",
                name: taskPriority.value.name ?? "",
                hexColor: taskPriority.value.hexColor ?? "",
                index: taskPriority.value.index ?? 0,
            };
            const updatedRole: TaskPriorityResponse = await taskPriorityService.update(payload);
            emit('update', updatedRole)
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
                    taskPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
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
                } else {
                    console.error(`TODO: missing taskPriorityId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("ProjectPriorityForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectPriorityForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.taskPriorityId) {
                onGet(props.taskPriorityId);
            } else {
                console.error(`TODO: missing taskPriorityId property for ${props.mode} action`);
            }
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon :component="props.mode == 'add' ? IconPlus : IconEdit" />
                {{ t(props.mode == "add" ?
                    "modules.taskPriority.components.TaskPriorityForm.headers.addTaskPriority" :
                    "modules.taskPriority.components.TaskPriorityForm.headers.updateTaskPriority") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="taskPriorityFormRef" :model="taskPriority" :rules="taskPriorityFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.taskPriority.components.TaskPriorityForm.inputs.name.placeholder')"
                    v-model:value="taskPriority.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.index.label')" path="index"
                show-feedback>
                <n-input-number :min="0"
                    :placeholder="t('modules.taskPriority.components.TaskPriorityForm.inputs.index.placeholder')"
                    v-model:value="taskPriority.index" required>
                </n-input-number>
            </n-form-item>
            <n-form-item :label="t('modules.taskPriority.components.TaskPriorityForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(taskPriority.hexColor ?? '#888888')" style="width: 100%;">
                        {{ taskPriority.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="taskPriority.hexColor">
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