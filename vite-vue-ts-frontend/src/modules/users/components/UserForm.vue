<script setup lang="ts">
    import { ref, reactive, computed, onMounted, nextTick, watch, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NRadio, NRadioGroup, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, type InputInst, NTooltip } from 'naive-ui';
    import { DONEO_ICON_ACTION_CANCEL, DONEO_ICON_ACTION_HIDE, DONEO_ICON_ACTION_SAVE, DONEO_ICON_ACTION_SHOW, DONEO_ICON_ADD, DONEO_ICON_EDIT, DONEO_ICON_EMAIL, DONEO_ICON_PASSWORD, DONEO_ICON_USER } from '../../../shared/types/icons';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { User, MAX_NAME_LENGTH, MAX_EMAIL_LENGTH, MIN_PASSWORD_LENGTH } from '../models/user';
    import { userService } from '../services/user'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { UserResponse } from '../types/dto';
    import { isValidEmail } from '../../../shared/composables/form-validators';
    import { appBus } from '../../../shared/composables/bus';

    interface Props {
        userId?: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'update', 'cancel'])

    const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const user = ref<User>(new User());

    const showPasswordField = ref<boolean>(true);

    const serverErrors = ref<Record<string, string>>({});

    const formRef = ref<FormInst | null>(null)

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
        return !user.value.name || state.ajaxRunning;
    });

    const onShowPasswordFormItem = async () => {
        showPasswordField.value = true;
        await nextTick();
        inputPasswordRef.value?.focus();
    };

    const onSave = async () => {
        serverErrors.value = {};
        formRef.value?.restoreValidation();
        try {
            await formRef.value?.validate();
            if (!props.userId) {
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
        formRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: UserResponse = await userService.get(id);
            user.value = new User(response);
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
            const addedUser: UserResponse = await userService.add(user.value.toAddUserRequestPayload());
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
            const updatedUser: UserResponse = await userService.update(user.value.toUpdateUserRequestPayload());
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
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.notFoundError");
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
                    formRef.value?.validate().then(() => { }).catch(() => { });
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
                }
            } else if (payload.to.includes("UserForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("UserForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.userId) {
            showPasswordField.value = false;
            onGet(props.userId);
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
                <n-icon class="doneo-mr-4px" :component="!props.userId ? DONEO_ICON_ADD : DONEO_ICON_EDIT" />
                {{
                    t(!props.userId ? "modules.user.components.UserForm.headers.addUser" :
                        "modules.user.components.UserForm.headers.editUser")
                }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="formRef" :model="user" :rules="state.ajaxRunning ? {} : userFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.user.components.UserForm.inputs.name.label')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('modules.user.components.UserForm.inputs.name.placeholder')"
                    v-model:value="user.name" :maxlength="MAX_NAME_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_USER" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.user.components.UserForm.inputs.email.label')" path="email" show-feedback>
                <n-input type="text" :placeholder="t('modules.user.components.UserForm.inputs.email.placeholder')"
                    v-model:value="user.email" :maxlength="MAX_EMAIL_LENGTH" :show-count="true"
                    :disabled="state.ajaxRunning" clearable autofocus>
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_EMAIL" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.user.components.UserForm.inputs.password.label')" path="password"
                show-feedback>
                <n-input v-if="showPasswordField" type="password"
                    :placeholder="t('modules.user.components.UserForm.inputs.password.placeholder')"
                    v-model:value="user.password" show-password-on="click" ref="inputPasswordRef"
                    :disabled="state.ajaxRunning">
                    <template #prefix>
                        <n-icon :component="DONEO_ICON_PASSWORD" />
                    </template>
                    <template #password-visible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="DONEO_ICON_ACTION_HIDE" />
                            </template>
                            {{ t("modules.user.components.UserForm.inputs.password.hidePasswordTooltipIcon") }}
                        </n-tooltip>
                    </template>
                    <template #password-invisible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="DONEO_ICON_ACTION_SHOW" />
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
                    <n-radio :disabled="state.ajaxRunning" :value="true" name="isSuperUser"
                        :label="t('modules.user.components.UserForm.radios.permissions.superUser.label')">
                    </n-radio>
                    <n-radio :disabled="state.ajaxRunning" :value="false" name="isSuperUser"
                        :label="t('modules.user.components.UserForm.radios.permissions.normalUser.label')">
                    </n-radio>
                </n-radio-group>
            </n-form-item>
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

<style lang="css" scoped></style>