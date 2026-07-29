<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_SAVE, DONEO_ICON_ADD, DONEO_ICON_EDIT, DONEO_ICON_NAME, DONEO_ICON_PALETTE } from '../../../shared/types/icons';

    import { ProjectPriority, MAX_NAME_LENGTH } from '../models/project-priority';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPriorityService } from '../services/project-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectPriorityResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';

    interface Props {
        projectPriorityId?: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'update', 'cancel'])

    const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPriority = ref<ProjectPriority>(new ProjectPriority());

    projectPriority.value.hexColor = generateRandomSoftHexColor();

    const serverErrors = ref<Record<string, string>>({});

    const formRef = ref<FormInst | null>(null)

    const projectPriorityFormRules: FormRules =
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

    watch(() => projectPriority.value.name, () => { delete serverErrors.value.name });
    watch(() => projectPriority.value.index, () => { delete serverErrors.value.index });

    const isSaveDisabled = computed<boolean>(() => {
        return !projectPriority.value.name || state.ajaxRunning;
    });

    const onSave = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        try {
            await formRef.value?.validate();
            if (!props.projectPriorityId) {
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
            const response: ProjectPriorityResponse = await projectPriorityService.get(id);
            if (response.id === id) {
                projectPriority.value = new ProjectPriority(response);
            } else {
                state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
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
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectPriorityForm.vue", method: "onGet" }, { err: fatalError });
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
            const addedProjectPriority: ProjectPriorityResponse = await projectPriorityService.add(projectPriority.value.toAddProjectPriorityRequestPayload());
            emit('add', addedProjectPriority)
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
                        case 404:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectPriority.components.ProjectPriorityForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.projectPriority.components.ProjectPriorityForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.addError");
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
            const updatedProjectPriority: ProjectPriorityResponse = await projectPriorityService.update(projectPriority.value.toUpdateProjectPriorityRequestPayload());
            emit('update', updatedProjectPriority)
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
                                serverErrors.value.name = "modules.projectPriority.components.ProjectPriorityForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.projectPriority.components.ProjectPriorityForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.updateError");
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
                if (props.projectPriorityId) {
                    onGet(props.projectPriorityId);
                }
            } else if (payload.to.includes("ProjectPriorityForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectPriorityForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.projectPriorityId) {
            onGet(props.projectPriorityId);
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon class="doneo-mr-4px" :component="!props.projectPriorityId ? DONEO_ICON_ADD : DONEO_ICON_EDIT" />
                {{
                    t(!props.projectPriorityId ?
                        "modules.projectPriority.components.ProjectPriorityForm.headers.addProjectPriority" :
                        "modules.projectPriority.components.ProjectPriorityForm.headers.updateProjectPriority")
                }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="formRef" :model="projectPriority" :rules="projectPriorityFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectPriority.components.ProjectPriorityForm.inputs.name.label')"
                path="name" show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.projectPriority.components.ProjectPriorityForm.inputs.name.placeholder')"
                    v-model:value="projectPriority.name" :maxlength="MAX_NAME_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_NAME" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.projectPriority.components.ProjectPriorityForm.inputs.index.label')"
                path="index" show-feedback>
                <n-input-number :min="0"
                    :placeholder="t('modules.projectPriority.components.ProjectPriorityForm.inputs.index.placeholder')"
                    v-model:value="projectPriority.index" :disabled="state.ajaxRunning">
                </n-input-number>
            </n-form-item>
            <n-form-item :label="t('modules.projectPriority.components.ProjectPriorityForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)" style="width: 100%;">
                        {{ projectPriority.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectPriority.hexColor"
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