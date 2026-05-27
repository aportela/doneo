<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette } from '@tabler/icons-vue';

    import { ProjectPriority, MAX_NAME_LENGTH } from '../models/project-priority';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPriorityService } from '../services/project-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectPriorityResponse, AddRequest, UpdateRequest } from '../types/dto';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectStatusFormProps {
        mode: FormMode;
        projectPriorityId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectStatusFormProps>();

    const { t } = useI18n();

    const projectPriority = ref<ProjectPriority>(
        new ProjectPriority(
            { id: "", name: "", hexColor: generateRandomSoftHexColor() }
        )
    );

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPriorityFormRef = ref<FormInst | null>(null)

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
    };

    watch(() => projectPriority.value.name, () => { delete serverErrors.value.name });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectPriority.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectPriorityFormRef.value?.restoreValidation();
        try {
            await projectPriorityFormRef.value?.validate();
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
        projectPriorityFormRef.value?.restoreValidation();
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
            if (state.ajaxErrors) {
                await nextTick();
                projectPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        projectPriorityFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: projectPriority.value.name,
                HexColor: projectPriority.value.hexColor,
            };
            const addedRole: ProjectPriorityResponse = await projectPriorityService.add(payload);
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
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectPriority.components.ProjectPriorityForm.warnings.nameAlreadyExists";
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
                await nextTick();
                projectPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        projectPriorityFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: projectPriority.value.id,
                name: projectPriority.value.name,
                HexColor: projectPriority.value.hexColor,
            };
            const updatedRole: ProjectPriorityResponse = await projectPriorityService.update(payload);
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
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.projectPriority.components.ProjectPriorityForm.warnings.nameAlreadyExists";
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
                await nextTick();
                projectPriorityFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPriorityForm.onGet")) {
                if (props.projectPriorityId) {
                    onGet(props.projectPriorityId);
                } else {
                    console.error(`TODO: missing projectPriorityId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("ProjectPriorityForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectPriorityForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.projectPriorityId) {
                onGet(props.projectPriorityId);
            } else {
                console.error(`TODO: missing projectPriorityId property for ${props.mode} action`);
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
                    "modules.projectPriority.components.ProjectPriorityForm.headers.addProjectPriority" :
                    "modules.projectPriority.components.ProjectPriorityForm.headers.updateProjectPriority") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectStatusFormRef" :model="projectPriority" :rules="projectPriorityFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectPriority.components.ProjectPriorityForm.inputs.name.label')"
                path="name" show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.projectPriority.components.ProjectPriorityForm.inputs.name.placeholder')"
                    v-model:value="projectPriority.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable
                    required autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.projectPriority.components.ProjectPriorityForm.inputs.preview.label')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)" style="width: 100%;">
                        {{ projectPriority.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectPriority.hexColor">
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
        <template #footer v-if="state.ajaxErrorMessage">
            <RemoteAPIAlert type="error" :title="t('shared.errorMessages.Error')" :message="state.ajaxErrorMessage" />
        </template>
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