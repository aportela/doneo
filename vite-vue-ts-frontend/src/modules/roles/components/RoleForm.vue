<script setup lang="ts">
    import { ref, reactive, computed, onMounted, type CSSProperties, watch, onBeforeUnmount } from 'vue';
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
                    allowCreate: false,
                    allowUpdate: false,
                    allowDelete: false,
                    allowView: false,
                    allowList: false,
                    allowExecute: false,
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
                        // conflict (invalid id || name ?)
                        default:
                            state.ajaxErrorMessage = t("There was a problem while adding the rpñe data");
                            break;
                    }
                    roleFormRef.value?.restoreValidation();
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while adding the role data");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onAdd" }, { err: fatalError });
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
                        // conflict (invalid id || name ?)
                        default:
                            state.ajaxErrorMessage = t("There was a problem while updating the role data");
                            break;
                    }
                    roleFormRef.value?.restoreValidation();
                    roleFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while updating the role data");
                    console.error("Unhandled API error", { file: "RoleForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
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
            <h4>Permissions</h4>

            <n-grid :y-gap="8" :cols="2">
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowCreate">
                        <template #checked>
                            Create allowed
                        </template>
                        <template #unchecked>
                            Create not allowed
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowUpdate">
                        <template #checked>
                            Update allowed
                        </template>
                        <template #unchecked>
                            Update not allowed
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowDelete">
                        <template #checked>
                            Delete allowed
                        </template>
                        <template #unchecked>
                            Delete not allowed
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowView">
                        <template #checked>
                            View allowed
                        </template>
                        <template #unchecked>
                            View not allowed
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowList">
                        <template #checked>
                            List allowed
                        </template>
                        <template #unchecked>
                            List not allowed
                        </template>
                    </n-switch>
                </n-gi>
                <n-gi>
                    <n-switch v-model:value="role.permissions.allowExecute">
                        <template #checked>
                            Execute allowed
                        </template>
                        <template #unchecked>
                            Execute not allowed
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
</style>