<script setup lang="ts">
    import { ref, reactive, computed, onMounted, watch, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NGrid, NGi, NSwitch } from 'naive-ui';
    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_SAVE, DONEO_ICON_ADD, DONEO_ICON_EDIT, DONEO_ICON_ROLE } from '../../../shared/types/icons';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { Role, MAX_NAME_LENGTH } from '../models/role';
    import { roleService } from '../services/role'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { RoleResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';

    interface Props {
        roleId?: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'update', 'cancel'])

    const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const role = ref<Role>(new Role());

    const serverErrors = ref<Record<string, string>>({});

    const formRef = ref<FormInst | null>(null);

    const roleFormRules: FormRules =
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

    watch(() => role.value.name, () => { delete serverErrors.value.name });

    const isSaveDisabled = computed<boolean>(() => {
        return !role.value.name || state.ajaxRunning;
    });

    const onSave = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        try {
            await formRef.value?.validate();
            if (!props.roleId) {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "RoleForm.vue", method: "onSave", details: "form validation error", error: error });
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
            const response: RoleResponse = await roleService.get(id);
            role.value = new Role(response);
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RoleForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.loadError");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onGet" }, { err: fatalError });
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
            const addedRole: RoleResponse = await roleService.add(role.value.toAddRoleRequestPayload());
            emit('add', addedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RoleForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.role.components.RoleForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.addError");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onAdd" }, { err: fatalError });
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
            const updatedRole: RoleResponse = await roleService.update(role.value.toUpdateRoleRequestPayload());
            emit('update', updatedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RoleForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.role.components.RoleForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.updateError");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onUpdate" }, { err: fatalError });
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
            if (payload.to.includes("RoleForm.onGet")) {
                if (props.roleId) {
                    onGet(props.roleId);
                }
            } else if (payload.to.includes("RoleForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("RoleForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.roleId) {
            onGet(props.roleId);
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
                <n-icon class="doneo-mr-4px" :component="!props.roleId ? DONEO_ICON_ADD : DONEO_ICON_EDIT" />
                {{
                    t(!props.roleId ? "modules.role.components.RoleForm.headers.addRole" :
                        "modules.role.components.RoleForm.headers.editRole")
                }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="formRef" :model="role" :rules="state.ajaxRunning ? {} : roleFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.role.components.RoleForm.inputs.name.label')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('modules.role.components.RoleForm.inputs.name.placeholder')"
                    v-model:value="role.name" :maxlength="MAX_NAME_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_ROLE" />
                    </template>
                </n-input>
            </n-form-item>
            <h4>{{ t("modules.role.components.RoleForm.headers.rolePermissions") }}</h4>

            <n-grid :x-gap="8" :y-gap="8" :cols="2">
                <n-gi>
                    <h4 class="doneo-permission-group-header">{{
                        t("modules.role.components.RoleForm.headers.projectPermissions")
                        }}</h4>
                    <n-switch v-model:value="role.permissions.allowUpdateProject" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.updateProjectAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.updateProjectDenied") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowDeleteProject" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.deleteProjectAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.deleteProjectDenied") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowViewProject" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.viewProjectAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.viewProjectDenied") }}
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <h4 class="doneo-permission-group-header">{{
                        t("modules.role.components.RoleForm.headers.taskPermissions")
                        }}
                    </h4>
                    <n-switch v-model:value="role.permissions.allowAddTask" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.addTaskAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.addTaskDenied") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowUpdateTask" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.updateTaskAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.updateTaskDenied") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowDeleteTask" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.deleteTaskAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.deleteTaskDenied") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowViewTask" class="doneo-permission-switch"
                        :disabled="state.ajaxRunning">
                        <template #checked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.viewTaskAllowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("modules.role.components.RoleForm.permissionSwitches.viewTaskDenied") }}
                        </template>
                    </n-switch>
                </n-gi>
            </n-grid>
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

<style lang="css" scoped>
    .doneo-permission-group-header {
        margin: 0px;
    }

    .doneo-permission-switch {
        display: block;
        margin: 8px 0px;
    }
</style>