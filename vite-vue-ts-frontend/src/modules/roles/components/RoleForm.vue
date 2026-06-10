<script setup lang="ts">
    import { ref, reactive, computed, onMounted, type CSSProperties, watch, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NGrid, NGi, NSwitch } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { Role, MAX_NAME_LENGTH } from '../models/role';
    import { roleService } from '../services/role'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { RoleResponse, AddRequest, UpdateRequest } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';
    import { getDefaultPermissions } from '../types/dto';

    interface RoleFormProps {
        mode: FormMode;
        roleId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<RoleFormProps>();

    const { t } = useI18n();

    const role = ref<Role>(new Role());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const showPasswordField = ref<boolean>(true);

    const serverErrors = ref<Record<string, string>>({});

    const roleFormRef = ref<FormInst | null>(null)

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
        return !role.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        roleFormRef.value?.restoreValidation();
        try {
            await roleFormRef.value?.validate();
            if (props.mode === "add") {
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
        roleFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: RoleResponse = await roleService.get(id);
            if (response.id === id) {
                role.value = new Role(response);
            } else {
                state.ajaxErrorMessage = t("modules.role.components.RoleForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RoleForm.onGet" } });
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
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        roleFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: role.value.name ?? "",
                permissions: role.value.permissions ?? getDefaultPermissions(),
            };
            const addedRole: RoleResponse = await roleService.add(payload);
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
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        roleFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: role.value.id ?? "",
                name: role.value.name ?? "",
                permissions: role.value.permissions ?? getDefaultPermissions(),
            };
            const updatedRole: RoleResponse = await roleService.update(payload);
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
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
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
                } else {
                    console.error(`TODO: missing roleId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("RoleForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("RoleForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            showPasswordField.value = false;
            if (props.roleId) {
                onGet(props.roleId);
            } else {
                console.error(`TODO: missing roleId property for ${props.mode} action`);
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
                {{ t(props.mode == "add" ? "modules.role.components.RoleForm.headers.addRole" :
                    "modules.role.components.RoleForm.headers.updateRole") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="roleFormRef" :model="role" :rules="state.ajaxRunning ? {} : roleFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.role.components.RoleForm.inputs.name.label')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('modules.role.components.RoleForm.inputs.name.placeholder')"
                    v-model:value="role.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
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

<style lang="css" scoped>
    .doneo-flex-center-align .n-icon {
        margin-right: 4px;
    }

    .doneo-permission-group-header {
        margin: 0px;
    }

    .doneo-permission-switch {
        display: block;
        margin: 8px 0px;
    }
</style>