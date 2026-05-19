<script setup lang="ts">
    import { ref, reactive, computed, onMounted, type CSSProperties, watch, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, NGrid, NGi, NSwitch } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { Role, maxNameLength } from '../models/role';
    import { roleService } from '../services/role'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { RoleResponse, AddRequest, UpdateRequest } from '../types/dto';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface RoleFormProps {
        mode: FormMode;
        roleId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<RoleFormProps>();

    const { t } = useI18n();

    const role = ref<Role>(
        new Role(
            {
                "id": "",
                name: "",
                permissions: {
                    allowUpdateProject: false,
                    allowDeleteProject: false,
                    allowViewProject: false,
                    allowAddTask: false,
                    allowUpdateTask: false,
                    allowDeleteTask: false,
                    allowViewTask: false,
                }
            }
        )
    );

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const showPasswordField = ref<boolean>(true);

    const serverErrors = ref<Record<string, string>>({});

    const roleFormRef = ref<FormInst | null>(null)

    const roleFormRules: FormRules =
    {
        name: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (!value) {
                    return new Error(t("roleFormNameFieldEmptyError"));
                }
                else if (value.length > maxNameLength) {
                    return new Error(t("roleFormNameFieldTooLargeError"));
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
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: RoleResponse = await roleService.get(id);
            if (response.id === id) {
                role.value = new Role(response);
                roleFormRef.value?.restoreValidation();
            } else {
                state.ajaxErrorMessage = t("There was a problem while loading the role data");
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
                            state.ajaxErrorMessage = t("roleFormIdNotFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while loading the role data");
                            break;
                    }
                    roleFormRef.value?.restoreValidation();
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while loading the role data");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onAdd = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: role.value.name,
                permissions: role.value.permissions,
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
                                serverErrors.value.name = "roleFormNameFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while adding the role data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while adding the role data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while adding the role data");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                roleFormRef.value?.restoreValidation();
                roleFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onUpdate = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: role.value.id,
                name: role.value.name,
                permissions: role.value.permissions,
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
                                serverErrors.value.name = "roleFormNameFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while updating the role data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while updating the role data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while updating the role data");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                roleFormRef.value?.restoreValidation();
                roleFormRef.value?.validate().then(() => { }).catch(() => { });
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
                {{ t(props.mode == "add" ? "Add role" : "Update role") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="roleFormRef" :model="role" :rules="state.ajaxRunning ? {} : roleFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('roleFormNameLabel')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('roleFormNameFieldPlaceholder')" v-model:value="role.name"
                    :maxlength="maxNameLength" :show-count="true" clearable required autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <h4>{{ t("Role permissions") }}</h4>

            <n-grid :x-gap="8" :y-gap="8" :cols="2">
                <n-gi>
                    <h4 class="doneo-permission-group-header">{{ t("Project permissions") }}</h4>
                    <n-switch v-model:value="role.permissions.allowUpdateProject" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("Update project allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("Update project allowed") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowDeleteProject" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("Delete project allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("Delete project allowed") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowViewProject" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("View project allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("View project allowed") }}
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <h4 class="doneo-permission-group-header">{{ t("Task permissions") }}</h4>
                    <n-switch v-model:value="role.permissions.allowAddTask" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("Add task allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("Add task allowed") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowUpdateTask" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("Update task allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("Update task allowed") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowDeleteTask" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("Delete task allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("Delete task allowed") }}
                        </template>
                    </n-switch>
                    <n-switch v-model:value="role.permissions.allowViewTask" class="doneo-permission-switch">
                        <template #checked>
                            {{ t("View task allowed") }}
                        </template>
                        <template #unchecked>
                            {{ t("View task allowed") }}
                        </template>
                    </n-switch>
                </n-gi>
            </n-grid>
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