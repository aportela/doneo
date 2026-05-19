<script setup lang="ts">
    import { ref, reactive, computed, onMounted, type CSSProperties, nextTick, watch, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NRadio, NRadioGroup, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon, type InputInst, NTooltip } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconEye, IconEyeCancel, IconMail, IconUser, IconUserEdit, IconUserPlus, IconKey } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { User, maxNameLength, maxEmailLength, minPasswordLength } from '../models/user';
    import { userService } from '../services/user'
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { UserResponse, AddRequest, UpdateRequest } from '../types/dto';
    import { isValidEmail } from '../../../shared/composables/form-validators';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface UserFormProps {
        mode: FormMode;
        userId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<UserFormProps>();

    const { t } = useI18n();

    const user = ref<User>(
        new User({ "id": "", name: "", email: "", permissions: { isSuperUser: false }, createdAt: 0, updatedAt: 0, deletedAt: 0, avatarUrl: "" })
    );

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
                if (!value) {
                    return new Error(t("userFormNameFieldEmptyError"));
                }
                else if (value.length > maxNameLength) {
                    return new Error(t("userFormNameFieldTooLargeError"));
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
                if (!value) {
                    return new Error(t("userFormEmailFieldEmptyError"));
                }
                else if (!isValidEmail(value)) {
                    return new Error(t("userFormEmailFieldInvalidError"));
                }
                else if (value.length > maxEmailLength) {
                    return new Error(t("userFormEmailFieldTooLargeError"));
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
                if (!showPasswordField.value) {
                    return true;
                }
                if (!value) {
                    return new Error(t("userFormPasswordFieldEmptyError"));
                }
                else if (value.length < minPasswordLength) {
                    return new Error(t("userFormPasswordFieldTooShortError"));
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
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: UserResponse = await userService.get(id);
            if (response.id === id) {
                user.value = new User(response);
                userFormRef.value?.restoreValidation();
            } else {
                state.ajaxErrorMessage = t("There was a problem while loading the user data");
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
                        case 404:
                            state.ajaxErrorMessage = t("userFormEmailNotFoundError"); // TODO: check this label, ! email but id ????
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while loading the user data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while loading the user data");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                userFormRef.value?.restoreValidation();
                userFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onAdd = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: user.value.name,
                email: user.value.email,
                permissions: {
                    isSuperUser: user.value.permissions.isSuperUser,
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
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "userFormNameFieldAlreadyExists";
                            } else if (apiError.details?.field === "email") {
                                serverErrors.value.email = "userFormEmailFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while adding the user data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while adding the user data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while adding the user data");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                userFormRef.value?.restoreValidation();
                userFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: user.value.id,
                name: user.value.name,
                password: user.value.password || undefined,
                email: user.value.email,
                permissions: {
                    isSuperUser: user.value.permissions.isSuperUser,
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
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "userFormNameFieldAlreadyExists";
                            } else if (apiError.details?.field === "email") {
                                serverErrors.value.email = "userFormEmailFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while updating the user data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while updating the user data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while updating the user data");
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                userFormRef.value?.restoreValidation();
                userFormRef.value?.validate().then(() => { }).catch(() => { });
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
                {{ t(props.mode == "add" ? "Add user" : "Update user") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="userFormRef" :model="user" :rules="state.ajaxRunning ? {} : userFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('userFormNameLabel')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('userFormNameFieldPlaceholder')" v-model:value="user.name"
                    :maxlength="maxNameLength" :show-count="true" clearable required autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('userFormEmailLabel')" path="email" show-feedback>
                <n-input type="text" :placeholder="t('userFormEmailFieldPlaceholder')" v-model:value="user.email"
                    :maxlength="maxEmailLength" :show-count="true" clearable required autofocus>
                    <template #prefix>
                        <n-icon :component="IconMail" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('userFormPasswordLabel')" path="password" show-feedback>
                <n-input v-if="showPasswordField" type="password" :placeholder="t('userFormPasswordFieldPlaceholder')"
                    v-model:value="user.password" show-password-on="click" ref="inputPasswordRef">
                    <template #prefix>
                        <n-icon :component="IconKey" />
                    </template>
                    <template #password-visible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="IconEyeCancel" />
                            </template>
                            {{ t("hide password") }}
                        </n-tooltip>
                    </template>
                    <template #password-invisible-icon>
                        <n-tooltip trigger="hover">
                            <template #trigger>
                                <n-icon :size="16" :component="IconEye" />
                            </template>
                            {{ t("show password") }}
                        </n-tooltip>
                    </template>
                </n-input>
                <n-button v-else @click="onShowPasswordFormItem" block>{{ t("userFormChangePasswordButtonLabel")
                    }}</n-button>
            </n-form-item>
            <n-form-item :label="t('userFormPermissionsLabel')">
                <n-radio-group v-model:value="user.permissions.isSuperUser" name="radiogroup">
                    <n-radio :value="true" name="aaa" :label="t('superUserPermission')">
                    </n-radio>
                    <n-radio :value="false" name="aaa" :label="t('normalUserPermission')">
                    </n-radio>
                </n-radio-group>
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

<style lang="css" scoped>
    .doneo-flex-center-align .n-icon {
        margin-right: 4px;
    }
</style>