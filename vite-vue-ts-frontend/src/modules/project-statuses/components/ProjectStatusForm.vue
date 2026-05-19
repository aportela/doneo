<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette } from '@tabler/icons-vue';

    import { ProjectStatus, maxNameLength } from '../models/project-status';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectStatusService } from '../services/project-status';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectStatusResponse, AddRequest, UpdateRequest } from '../types/dto';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectStatusFormProps {
        mode: FormMode;
        projectStatusId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectStatusFormProps>();

    const { t } = useI18n();

    const projectStatus = ref<ProjectStatus>(
        new ProjectStatus(
            { id: "", name: "", hexColor: generateRandomSoftHexColor(), index: 0 }
        )
    );

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectStatusFormRef = ref<FormInst | null>(null)

    const projectStatusFormRules: FormRules =
    {
        name: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (!value) {
                    return new Error(t("projectStatusFormNameFieldEmptyError"));
                }
                else if (value.length > maxNameLength) {
                    return new Error(t("projectStatusFormNameFieldTooLargeError"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
    };

    watch(() => projectStatus.value.name, () => { delete serverErrors.value.name });


    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectStatus.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
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
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectStatusResponse = await projectStatusService.get(id);
            if (response.id === id) {
                projectStatus.value = new ProjectStatus(response);
                projectStatusFormRef.value?.restoreValidation();
            } else {
                state.ajaxErrorMessage = t("There was a problem while loading the project status data");
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
                        case 404:
                            state.ajaxErrorMessage = t("projectStatusFormIdNotFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while loading the project status data");
                            break;
                    }
                    projectStatusFormRef.value?.restoreValidation();
                    projectStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while loading the project status data");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onAdd = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: projectStatus.value.name,
                HexColor: projectStatus.value.hexColor,
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
                        case 409:
                        // conflict (invalid id || name ?)
                        default:
                            state.ajaxErrorMessage = t("There was a problem while adding the project status data");
                            break;
                    }
                    projectStatusFormRef.value?.restoreValidation();
                    projectStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while adding the project status data");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: projectStatus.value.id,
                name: projectStatus.value.name,
                HexColor: projectStatus.value.hexColor,
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
                        case 409:
                        // conflict (invalid id || name ?)
                        default:
                            state.ajaxErrorMessage = t("There was a problem while updating the project status data");
                            break;
                    }
                    projectStatusFormRef.value?.restoreValidation();
                    projectStatusFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while updating the project status data");
                    console.error("Unhandled API error", { file: "ProjectStatusForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
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
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon :component="props.mode == 'add' ? IconPlus : IconEdit" />
                {{ t(props.mode == "add" ? "Add project status" : "Update project status") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectStatusFormRef" :model="projectStatus" :rules="projectStatusFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('projectStatusFormNameFieldPlaceholder')"
                    v-model:value="projectStatus.name" :maxlength="maxNameLength" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectStatus.hexColor)" style="width: 100%;">
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
        <template #footer v-if="state.ajaxErrorMessage">
            <RemoteAPIAlert type="error" :title="t('Error')" :message="state.ajaxErrorMessage" />
        </template>
        <template #action>
            <n-flex>
                <n-button @click="onSave" :disabled="isSaveDisabled">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy" />
                    </template>
                    {{ t("Save") }}
                </n-button>
                <n-button @click="onCancel" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconCancel" />
                    </template>
                    {{ t("Cancel") }}
                </n-button>
            </n-flex>
        </template>
    </n-card>

</template>

<style lang="css" scoped></style>