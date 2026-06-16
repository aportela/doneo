<script setup lang="ts">
    import { ref, reactive, computed, onMounted, type CSSProperties, nextTick, watch, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NRadio, NRadioGroup, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, type InputInst, NTooltip } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconEye, IconEyeCancel, IconMail, IconUser, IconUserEdit, IconUserPlus, IconKey } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { User, MAX_NAME_LENGTH, MAX_EMAIL_LENGTH, MIN_PASSWORD_LENGTH } from '../models/user';
    import { userService } from '../services/user'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { UserResponse, AddRequest, UpdateRequest } from '../types/dto';
    import { isValidEmail } from '../../../shared/composables/form-validators';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface UserFormProps {
        mode: FormMode;
        userId?: string | null;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<UserFormProps>();

    const { t } = useI18n();

    const user = ref<User>(new User());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const showPasswordField = ref<boolean>(true);

    const serverErrors = ref<Record<string, string>>({});

    const userFormRef = ref<FormInst | null>(null)

    const inputPasswordRef = ref<InputInst | null>(null);

    const userFormRules: FormRules =
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
        email: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (!isValidEmail(value)) {
                    return new Error(t("shared.warningMessages.fieldHasInvalidFormat"));
                }
                else if (value.length > MAX_EMAIL_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.email) {
                    return new Error(t(serverErrors.value.email));
                } else {
                    return true;
                }
            }, trigger: ['blur'],
        },
        password: {
            required: showPasswordField.value,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!showPasswordField.value) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length < MIN_PASSWORD_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldIsBelowMinimumLength"));
                } else if (serverErrors.value.password) {
                    return new Error(t(serverErrors.value.password));
                } else {
                    return true;
                }
            },
            trigger: ['blur']
        }
    };

    watch(() => user.value.name, () => { delete serverErrors.value.name });
    watch(() => user.value.email, () => { delete serverErrors.value.email });
    watch(() => user.value.password, () => { delete serverErrors.value.password });

    const isSaveDisabled = computed<boolean>(() => {
        return !user.value.name;
    });

    const onShowPasswordFormItem = async () => {
        showPasswordField.value = true;
        await nextTick();
        inputPasswordRef.value?.focus();
    };

    const onSave = async () => {
        serverErrors.value = {};
        userFormRef.value?.restoreValidation();
        try {
            await userFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "UserForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onGet = async (id: string) => {
        serverErrors.value = {};
        userFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: UserResponse = await userService.get(id);
            if (response.id === id) {
                user.value = new User(response);
            } else {
                state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UserForm.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.loadError");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        userFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: user.value.name ?? "",
                email: user.value.email ?? "",
                permissions: {
                    isSuperUser: user.value.permissions?.isSuperUser ?? false,
                }
            };
            const addedUser: UserResponse = await userService.add(payload);
            emit('add', addedUser)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UserForm.onAdd" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.user.components.UserForm.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                serverErrors.value.email = "modules.user.components.UserForm.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.addError");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        userFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: user.value.id ?? "",
                name: user.value.name ?? "",
                password: user.value.password || undefined,
                email: user.value.email ?? "",
                permissions: {
                    isSuperUser: user.value.permissions?.isSuperUser ?? false,
                }
            };
            const updatedUser: UserResponse = await userService.update(payload);
            emit('update', updatedUser)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UserForm.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.user.components.UserForm.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                serverErrors.value.email = "modules.user.components.UserForm.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("UserForm.onGet")) {
                if (props.userId) {
                    onGet(props.userId);
                } else {
                    console.error(`TODO: missing userId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("UserForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("UserForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            showPasswordField.value = false;
            if (props.userId) {
                onGet(props.userId);
            } else {
                console.error(`TODO: missing userId property for ${props.mode} action`);
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
                <n-icon :component="props.mode == 'add' ? IconUserPlus : IconUserEdit" />
                {{ t(props.mode == "add" ? "modules.user.components.UserForm.headers.addUser" :
                    "modules.user.components.UserForm.headers.updateUser") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="userFormRef" :model="user" :rules="state.ajaxRunning ? {} : userFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.user.components.UserForm.inputs.name.label')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('modules.user.components.UserForm.inputs.name.placeholder')"
                    v-model:value="user.name" :maxlength="MAX_NAME_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.user.components.UserForm.inputs.email.label')" path="email" show-feedback>
                <n-input type="text" :placeholder="t('modules.user.components.UserForm.inputs.email.placeholder')"
                    v-model:value="user.email" :maxlength="MAX_EMAIL_LENGTH" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconMail" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.user.components.UserForm.inputs.password.label')" path="password"
                show-feedback>
                <n-input v-if="showPasswordField" type="password"
                    :placeholder="t('modules.user.components.UserForm.inputs.password.placeholder')"
                    v-model:value="user.password" show-password-on="click" ref="inputPasswordRef">
                    <template #prefix>
                        <n-icon :component="IconKey" />
                    </template>
                    <template #password-visible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="IconEyeCancel" />
                            </template>
                            {{ t("modules.user.components.UserForm.inputs.password.hidePasswordTooltipIcon") }}
                        </n-tooltip>
                    </template>
                    <template #password-invisible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="IconEye" />
                            </template>
                            {{ t("modules.user.components.UserForm.inputs.password.showPasswordTooltipIcon") }}
                        </n-tooltip>
                    </template>
                </n-input>
                <n-button v-else @click="onShowPasswordFormItem" block :disabled="state.ajaxRunning">{{
                    t("modules.user.components.UserForm.buttons.changePassword.label")
                }}</n-button>
            </n-form-item>
            <n-form-item :label="t('modules.user.components.UserForm.radios.permissions.label')">
                <n-radio-group v-model:value="user.permissions.isSuperUser" name="radiogroup">
                    <n-radio :value="true" name="isSuperUser"
                        :label="t('modules.user.components.UserForm.radios.permissions.superUser.label')">
                    </n-radio>
                    <n-radio :value="false" name="isSuperUser"
                        :label="t('modules.user.components.UserForm.radios.permissions.normalUser.label')">
                    </n-radio>
                </n-radio-group>
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

<style lang="css" scoped>
    .doneo-flex-center-align .n-icon {
        margin-right: 4px;
    }
</style>