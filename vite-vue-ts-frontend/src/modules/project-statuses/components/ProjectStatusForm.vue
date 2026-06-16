<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NInputNumber, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NTooltip } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette, IconStar, IconCalendarBolt, IconCalendarCancel, IconCalendarMinus } from '@tabler/icons-vue';

    import { ProjectStatus, MAX_NAME_LENGTH } from '../models/project-status';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectStatusService } from '../services/project-status';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectStatusResponse, AddRequest, UpdateRequest } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectStatusFormProps {
        mode: FormMode;
        projectStatusId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectStatusFormProps>();

    const { t } = useI18n();

    const projectStatus = ref<ProjectStatus>(new ProjectStatus());
    projectStatus.value.hexColor = generateRandomSoftHexColor();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectStatusFormRef = ref<FormInst | null>(null)

    const projectStatusFormRules: FormRules =
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

    watch(() => projectStatus.value.name, () => { delete serverErrors.value.name });
    watch(() => projectStatus.value.index, () => { delete serverErrors.value.index });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectStatus.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectStatusFormRef.value?.restoreValidation();
        try {
            await projectStatusFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "ProjectStatusForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onGet = async (id: string) => {
        serverErrors.value = {};
        projectStatusFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectStatusResponse = await projectStatusService.get(id);
            if (response.id === id) {
                projectStatus.value = new ProjectStatus(response);
            } else {
                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onGet" }, { err: fatalError });
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
        projectStatusFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: projectStatus.value.name ?? "",
                hexColor: projectStatus.value.hexColor ?? "",
                index: projectStatus.value.index ?? 0,
                flags: {
                    defaultStatusOnCreation: projectStatus.value.flags?.defaultStatusOnCreation ?? false,
                    fillEmptyStartDate: projectStatus.value.flags?.fillEmptyStartDate ?? false,
                    setStartDate: projectStatus.value.flags?.setStartDate ?? false,
                    fillEmptyFinishDate: projectStatus.value.flags?.fillEmptyFinishDate ?? false,
                    setFinishDate: projectStatus.value.flags?.setFinishDate ?? false,
                    unsetFinishDateOnLeave: projectStatus.value.flags?.unsetFinishDateOnLeave ?? false,
                }
            };
            const addedRole: ProjectStatusResponse = await projectStatusService.add(payload);
            emit('add', addedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectStatus.components.ProjectStatusForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.projectStatus.components.ProjectStatusForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.addError");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    projectStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        projectStatusFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: projectStatus.value.id ?? "",
                name: projectStatus.value.name ?? "",
                hexColor: projectStatus.value.hexColor ?? "",
                index: projectStatus.value.index ?? 0,
                flags: {
                    defaultStatusOnCreation: projectStatus.value.flags?.defaultStatusOnCreation ?? false,
                    fillEmptyStartDate: projectStatus.value.flags?.fillEmptyStartDate ?? false,
                    setStartDate: projectStatus.value.flags?.setStartDate ?? false,
                    fillEmptyFinishDate: projectStatus.value.flags?.fillEmptyFinishDate ?? false,
                    setFinishDate: projectStatus.value.flags?.setFinishDate ?? false,
                    unsetFinishDateOnLeave: projectStatus.value.flags?.unsetFinishDateOnLeave ?? false,
                }
            };
            const updatedRole: ProjectStatusResponse = await projectStatusService.update(payload);
            emit('update', updatedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectStatus.components.ProjectStatusForm.warnings.nameAlreadyExists";
                            } else if (apiError.details?.field === "index") {
                                serverErrors.value.index = "modules.projectStatus.components.ProjectStatusForm.warnings.indexAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectStatus.components.ProjectStatusForm.errors.updateError");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    projectStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectStatusForm.onGet")) {
                if (props.projectStatusId) {
                    onGet(props.projectStatusId);
                } else {
                    console.error(`TODO: missing projectStatusId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("ProjectStatusForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectStatusForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.projectStatusId) {
                onGet(props.projectStatusId);
            } else {
                console.error(`TODO: missing projectStatusId property for ${props.mode} action`);
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
                {{ t(props.mode == "add" ? "modules.projectStatus.components.ProjectStatusForm.headers.addProjectStatus"
                    : "modules.projectStatus.components.ProjectStatusForm.headers.updateProjectStatus") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectStatusFormRef" :model="projectStatus" :rules="projectStatusFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectStatus.components.ProjectStatusForm.inputs.name.label')" path="name"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.projectStatus.components.ProjectStatusForm.inputs.name.placeholder')"
                    v-model:value="projectStatus.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable
                    required autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-flex>
                <n-form-item :label="t('modules.projectStatus.components.ProjectStatusForm.inputs.index.label')"
                    path="index" show-feedback>
                    <n-input-number :min="0"
                        :placeholder="t('modules.projectStatus.components.ProjectStatusForm.inputs.index.placeholder')"
                        v-model:value="projectStatus.index" required>
                    </n-input-number>
                </n-form-item>
                <n-form-item :label="t('modules.projectStatus.components.ProjectStatusForm.inputs.flags.label')">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconStar" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.defaultStatusOnCreation }"
                                @click="projectStatus.flags.defaultStatusOnCreation = !projectStatus.flags.defaultStatusOnCreation" />
                        </template>
                        {{ t(projectStatus.flags.defaultStatusOnCreation ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasDefaultStatusOnCreation"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotdefaultStatusOnCreation")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconCalendarBolt" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.fillEmptyStartDate }"
                                @click="projectStatus.flags.fillEmptyStartDate = !projectStatus.flags.fillEmptyStartDate" />
                        </template>
                        {{ t(projectStatus.flags.fillEmptyStartDate ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasFillEmptyStartDate"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotFillEmptyStartDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconCalendarCancel" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.setStartDate }"
                                @click="projectStatus.flags.setStartDate = !projectStatus.flags.setStartDate" />
                        </template>
                        {{ t(projectStatus.flags.setStartDate ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasSetStartDate"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotSetStartDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconCalendarBolt" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.fillEmptyFinishDate }"
                                @click="projectStatus.flags.fillEmptyFinishDate = !projectStatus.flags.fillEmptyFinishDate" />
                        </template>
                        {{ t(projectStatus.flags.fillEmptyFinishDate ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasFillEmptyFinishDate"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotFillEmptyFinishDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconCalendarCancel" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.setFinishDate }"
                                @click="projectStatus.flags.setFinishDate = !projectStatus.flags.setFinishDate" />
                        </template>
                        {{ t(projectStatus.flags.setFinishDate ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasSetFinishDate"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotSetFinishDate")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="IconCalendarMinus" :size="flagIconSize" class="doneo-cursor-help"
                                :class="{ 'doneo-disabled-icon': !projectStatus.flags.unsetFinishDateOnLeave }"
                                @click="projectStatus.flags.unsetFinishDateOnLeave = !projectStatus.flags.unsetFinishDateOnLeave" />
                        </template>
                        {{ t(projectStatus.flags.unsetFinishDateOnLeave ?
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasUnsetFinishDateOnLeave"
                            :
                            "modules.projectStatus.components.ProjectStatusesTable.body.columns.permissionsHints.hasNotUnsetFinishDateOnLeave")
                        }}
                    </n-tooltip>
                </n-form-item>
            </n-flex>
            <n-form-item :label="t('modules.projectStatus.components.ProjectStatusForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectStatus.hexColor ?? '#888888')"
                        style="width: 100%;">
                        {{ projectStatus.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectStatus.hexColor">
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