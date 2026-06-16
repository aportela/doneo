<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette } from '@tabler/icons-vue';

    import { ProjectType, MAX_NAME_LENGTH } from '../models/project-type';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTypeService } from '../services/project-type';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectTypeResponse, AddRequest, UpdateRequest } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectTypeFormProps {
        mode: FormMode;
        projectTypeId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectTypeFormProps>();

    const { t } = useI18n();

    const projectType = ref<ProjectType>(new ProjectType());
    projectType.value.hexColor = generateRandomSoftHexColor();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypeFormRef = ref<FormInst | null>(null)

    const projectTypeFormRules: FormRules =
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
    };

    watch(() => projectType.value.name, () => { delete serverErrors.value.name });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectType.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        try {
            await projectTypeFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "ProjectTypeForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onGet = async (id: string) => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectTypeResponse = await projectTypeService.get(id);
            if (response.id === id) {
                projectType.value = new ProjectType(response);
            } else {
                state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onGet" }, { err: fatalError });
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
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: projectType.value.name ?? "",
                HexColor: projectType.value.hexColor ?? "",
            };
            const addedRole: ProjectTypeResponse = await projectTypeService.add(payload);
            emit('add', addedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectType.components.ProjectTypeForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.addError");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    projectTypeFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: projectType.value.id ?? "",
                name: projectType.value.name ?? "",
                HexColor: projectType.value.hexColor ?? "",
            };
            const updatedRole: ProjectTypeResponse = await projectTypeService.update(payload);
            emit('update', updatedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectType.components.ProjectTypeForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProjectTypeForm.errors.updateError");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    projectTypeFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectTypeForm.onGet")) {
                if (props.projectTypeId) {
                    onGet(props.projectTypeId);
                } else {
                    console.error(`TODO: missing projectTypeId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("ProjectTypeForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectTypeForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.projectTypeId) {
                onGet(props.projectTypeId);
            } else {
                console.error(`TODO: missing projectTypeId property for ${props.mode} action`);
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
                {{ t(props.mode == "add" ? "modules.projectType.components.ProjectTypeForm.headers.addProjectType" :
                    "modules.projectType.components.ProjectTypeForm.headers.updateProjectType") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectTypeFormRef" :model="projectType" :rules="projectTypeFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectType.components.ProjectTypeForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.projectType.components.ProjectTypeForm.inputs.name.placeholder')"
                    v-model:value="projectType.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectType.hexColor ?? '#888888')" style="width: 100%;">
                        {{ projectType.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectType.hexColor">
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