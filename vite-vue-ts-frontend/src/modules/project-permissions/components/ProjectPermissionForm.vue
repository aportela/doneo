<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NFlex, NButton, NForm, NFormItem, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconEdit, IconPlus } from '@tabler/icons-vue';

    import { ProjectPermission } from '../models/project-permission.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPermissionService } from '../services/project-permission.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectPermissionResponse, AddRequest, } from '../types/dto';
    import { getDefaultPermissions } from '../../roles/types/dto.ts';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';
    import UserSelector from '../../users/components/UserSelector.vue';
    import RoleSelector from '../../roles/components/RoleSelector.vue';

    interface ProjectPermissionFormProps {
        mode: FormMode;
        projectId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectPermissionFormProps>();

    const { t } = useI18n();

    const projectPermission = ref<ProjectPermission>(new ProjectPermission());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPermissionFormRef = ref<FormInst | null>(null)

    const projectPermissionFormRules: FormRules =
    {
    };

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectPermission.value.user.id || !projectPermission.value.role.id;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectPermissionFormRef.value?.restoreValidation();
        try {
            await projectPermissionFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                console.error("invalid form mode", { file: "ProjectPermissionForm.vue", method: "onSave" });
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "ProjectPermissionForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onAdd = async () => {
        if (props.projectId) {
            serverErrors.value = {};
            projectPermissionFormRef.value?.restoreValidation();
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const payload: AddRequest = {
                    user: {
                        id: projectPermission.value.user.id ?? "",
                        name: "",
                    },
                    role: {
                        id: projectPermission.value.role.id ?? "",
                        name: "",
                        permissions: getDefaultPermissions(),
                    }
                };
                const addedPermission: ProjectPermissionResponse = await projectPermissionService.add(props.projectId, payload);
                emit('add', addedPermission)
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPermissionForm.onAdd" } });
                                break;
                            case 409:
                                // TODO
                                if (apiError.details?.field === "userId") {
                                    serverErrors.value.name = "modules.projectPermission.components.ProjectPermissionForm.warnings.userAlreadyExists";
                                } else if (apiError.details?.field === "roleId") {
                                    serverErrors.value.name = "modules.projectPermission.components.ProjectPermissionForm.warnings.roleAlreadyExists";
                                } else {
                                    state.ajaxErrorMessage = t("modules.projectPermission.components.ProjectPermissionForm.errors.addError");
                                }
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.ProjectPermissionForm.errors.addError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.ProjectPermissionForm.errors.addError");
                        console.error("Unhandled API error", { file: "ProjectPermissionForm.vue", method: "onAdd" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrors) {
                    await nextTick();
                    projectPermissionFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        } else {
            console.error("project id not set", { file: "ProjectPermissionForm.vue", method: "onAdd" });
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPermissionForm.onAdd")) {
                onAdd();
            }
        });
        if (props.mode !== "add") {
            console.error("invalid form mode", { file: "ProjectPermissionForm.vue", method: "onSave" });
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
                <!-- TOOD: icon alignment ??? -->
                <n-icon :component="props.mode == 'add' ? IconPlus : IconEdit" />
                {{ t("modules.projectPermission.components.ProjectPermissionForm.headers.addProjectPermission") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectPermissionFormRef" :model="projectPermission" :rules="projectPermissionFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectPermission.components.ProjectPermissionForm.inputs.user.label')">
                <UserSelector v-model:id="projectPermission.user.id"
                    :placeholder="t('modules.projectPermission.components.ProjectPermissionForm.inputs.user.placeholder')" />
            </n-form-item>
            <n-form-item :label="t('modules.projectPermission.components.ProjectPermissionForm.inputs.role.label')">
                <RoleSelector v-model:id="projectPermission.role.id"
                    :placeholder="t('modules.projectPermission.components.ProjectPermissionForm.inputs.role.placeholder')" />
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