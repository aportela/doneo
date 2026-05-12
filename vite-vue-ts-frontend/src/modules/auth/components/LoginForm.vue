<script setup lang="ts">
    import { ref, reactive, watch, nextTick, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NIcon, NSpin, NForm, NFormItem, NInput, NButton, type FormItemRule, type FormInst, type FormRules, type InputInst, NTooltip } from 'naive-ui'
    import { IconEye, IconEyeCancel, IconMail, IconKey } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from "../../../shared/types/ajaxState";
    import { isValidEmail } from '../../../shared/composables/form-validators';
    import { createStorageEntry } from '../../../shared/composables/localStorage';
    import { useSessionStore } from "../../../stores/session";
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import { authService } from '../../../modules/auth/services/auth';
    import type { SignInRequest, SignInResponse } from '../../../modules/auth/types/dto';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { maxEmailLength, minPasswordLength, User } from "../../../modules/users/models/user";

    type signInFormValuesInterface = {
        email: string;
        password: string;
    };

    const emit = defineEmits(['success'])

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const localStorageLastUsedEmail = createStorageEntry<string | null>("lastUsedEmail", null);

    const lastUsedEmail = localStorageLastUsedEmail.get();

    const inputEmailRef = ref<InputInst | null>(null);
    const inputPasswordRef = ref<InputInst | null>(null);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const signInFormRef = ref<FormInst | null>(null)

    const signinFormValues = ref<signInFormValuesInterface>({
        email: lastUsedEmail || "",
        password: "",
    });

    const signInFormRules: FormRules = {
        email: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (!value) {
                    return new Error(t("loginFormEmailFieldEmptyError"));
                }
                else if (!isValidEmail(value)) {
                    return new Error(t("loginFormEmailFieldInvalidError"));
                }
                else if (value.length > maxEmailLength) {
                    return new Error(t("loginFormEmailFieldTooLargeError"));
                } else if (serverErrors.value.email) {
                    return new Error(t(serverErrors.value.email));
                } else {
                    return true
                }
            },
            trigger: ['blur'],
        },
        password: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (!value) {
                    return new Error(t("loginFormPasswordFieldEmptyError"));
                }
                else if (value.length < minPasswordLength) {
                    return new Error(t("loginFormPasswordFieldTooShortError"));
                } else if (serverErrors.value.password) {
                    return new Error(t(serverErrors.value.password));
                } else {
                    return true
                }
            },
            trigger: ['blur'],
        },
    };

    watch(() => signinFormValues.value.email, () => { delete serverErrors.value.email });
    watch(() => signinFormValues.value.password, () => { delete serverErrors.value.password });

    const onSubmit = async () => {
        if (signinFormValues.value.email && signinFormValues.value.password) {
            serverErrors.value = {}
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const payload: SignInRequest = {
                    email: signinFormValues.value.email,
                    password: signinFormValues.value.password
                };
                const response: SignInResponse = await authService.signIn(payload);
                sessionStore.setAccessToken(response.accessToken.token, response.accessToken.expiresAtTimestamp);
                sessionStore.setUser(new User(response.user));
                localStorageLastUsedEmail.set(signinFormValues.value.email);
                emit('success');
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 404:
                                serverErrors.value.email = "loginFormEmailFieldAPINotFound";
                                break;
                            case 401:
                                serverErrors.value.password = "loginFormPasswordFieldIncorrect";
                                break;
                            default:
                                state.ajaxErrorMessage = t("Invalid API response code");
                                break;
                        }
                        signInFormRef.value?.restoreValidation();
                        signInFormRef.value?.validate().then(() => { }).catch(() => { });
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("Uncaught exception");
                        console.error("Fatal error", { file: "LoginForm.vue", method: "onSubmit", details: "uncaught exception", error: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("Fatal error", { file: "LoginForm.vue", method: "onSubmit", details: "missing email/password values" });
        }
    }

    const onSignIn = async (e: MouseEvent | KeyboardEvent) => {
        e.preventDefault()
        serverErrors.value = {};
        try {
            await signInFormRef.value?.validate();
            onSubmit();
        }
        catch (error: any) {
            console.warn("Warning", { file: "LoginForm.vue", method: "onSignIn", details: "form validation error", error: error });
        }
    }

    onMounted(async () => {
        await nextTick();
        if (signinFormValues.value.email == "") {
            inputEmailRef.value?.focus();
        } else {
            inputPasswordRef.value?.focus();
        }
    });
</script>

<template>
    <h2 class="title">{{ t("Login to your account") }}</h2>
    <n-spin :show="state.ajaxRunning" stroke="pink">
        <n-form ref="signInFormRef" :model="signinFormValues" label-width="100px" :rules="signInFormRules">
            <n-form-item :label="t('loginFormEmailLabel')" path="email" show-feedback>
                <n-input type="text" v-model:value="signinFormValues.email"
                    :placeholder="t('loginFormEmailFieldPlaceholder')" :disabled="state.ajaxRunning"
                    ref="inputEmailRef">
                    <template #prefix>
                        <n-icon :component="IconMail" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('loginFormPasswordLabel')" path="password" show-feedback>
                <n-input v-model:value="signinFormValues.password" type="password"
                    :placeholder="t('loginFormPasswordFieldPlaceholder')" show-password-on="click"
                    :disabled="state.ajaxRunning" @keydown.enter="onSignIn" ref="inputPasswordRef">
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
            </n-form-item>
            <n-form-item>
                <n-button secondary @click="onSignIn" block :disabled="state.ajaxRunning">{{
                    t("Sign in")
                    }}</n-button>
            </n-form-item>
        </n-form>
    </n-spin>
    <RemoteAPIAlert v-if="state.ajaxErrorMessage" type="error" :title="t('Sign in error')"
        :message="state.ajaxErrorMessage" />
</template>

<style lang="css" scoped></style>