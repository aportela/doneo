<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy } from '@tabler/icons-vue';

    import { ProjectPermission } from '../models/project-permission.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPermissionService } from '../services/project-permission.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectPermissionResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import UserSelector from '../../users/components/UserSelector.vue';
    import RoleSelector from '../../roles/components/RoleSelector.vue';
    import { DONEO_ICON_ADD } from '../../../shared/types/icons.ts';

    interface Props {
        projectId: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'cancel'])

    const { t } = useI18n();

    const projectPermission = ref<ProjectPermission>(new ProjectPermission());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPermissionFormRef = ref<FormInst | null>(null)

    const projectPermissionFormRules: FormRules =
    {
        user: {
            id: {
                required: true,
                validator: (_rule: FormItemRule, value: string) => {
                    if (state.ajaxRunning) {
                        return true;
                    }
                    if (!value?.trim()) {
                        return new Error(t("shared.warningMessages.fieldIsRequired"));
                    } else if (serverErrors.value.userId) {
                        return new Error(t(serverErrors.value.userId));
                    } else {
                        return true;
                    }
                },
                trigger: ['blur'],
            },
        },
        role: {
            id: {
                required: true,
                validator: (_rule: FormItemRule, value: string) => {
                    if (state.ajaxRunning) {
                        return true;
                    }
                    if (!value?.trim()) {
                        return new Error(t("shared.warningMessages.fieldIsRequired"));
                    } else {
                        return true;
                    }
                },
                trigger: ['blur'],
            },
        },
    };

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectPermission.value.user.id || !projectPermission.value.role.id || state.ajaxRunning;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectPermissionFormRef.value?.restoreValidation();
        try {
            await projectPermissionFormRef.value?.validate();
            await onAdd();
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
                const addedProjectPermission: ProjectPermissionResponse = await projectPermissionService.add(props.projectId, projectPermission.value.toAddRequestPayload());
                emit('add', new ProjectPermission(addedProjectPermission));
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
                                if (apiError.details?.field === "userId") {
                                    // TODO:
                                    serverErrors.value.userId = "modules.projectPermission.components.ProjectPermissionForm.warnings.userAlreadyExists";
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
                    if (state.ajaxErrorMessage) {
                        appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                    } else {
                        await nextTick();
                        projectPermissionFormRef.value?.validate().then(() => { }).catch(() => { });
                    }
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
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon class="doneo-mr-4px" :component="DONEO_ICON_ADD" />
                {{ t("modules.projectPermission.components.ProjectPermissionForm.headers.addProjectPermission") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectPermissionFormRef" :model="projectPermission" :rules="projectPermissionFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.projectPermission.components.ProjectPermissionForm.inputs.user.label')"
                path="user.id">
                <UserSelector auto-focus required v-model:id="projectPermission.user.id"
                    :placeholder="t('modules.projectPermission.components.ProjectPermissionForm.inputs.user.placeholder')"
                    :disabled="state.ajaxRunning" />
            </n-form-item>
            <n-form-item :label="t('modules.projectPermission.components.ProjectPermissionForm.inputs.role.label')"
                path="role.id">
                <RoleSelector required v-model:id="projectPermission.role.id"
                    :placeholder="t('modules.projectPermission.components.ProjectPermissionForm.inputs.role.placeholder')"
                    :disabled="state.ajaxRunning" />
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