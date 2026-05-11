<script setup lang="ts">
    import { ref, reactive, computed, onMounted, watch, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconEye, IconEyeCancel } from '@tabler/icons-vue';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../types/ajaxState';
    import { User, maxNameLength, maxEmailLength } from '../../api/models/user';
    import { userService } from '../../api/services/user';
    import { handleAPIError } from '../../api/client/errorHandler';
    import type { AddRequest, GetResponse, UpdateRequest } from '../../api/types/dto/user';
    import { required, minLength, validEmail, runValidators, maxLength } from '../../composables/form-validators';
    import RemoteAPIAlert from '../alerts/RemoteAPIAlert.vue';

    type Mode = "add" | "update";

    interface UserFormProps {
        mode: Mode;
        userId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<UserFormProps>();

    const { t } = useI18n();

    const user = ref<User>(
        new User({ "id": "", name: "", email: "", isSuperUser: false, createdAt: 0, updatedAt: 0, deletedAt: 0, avatarUrl: "" })
    );

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const userFormRef = ref<FormInst | null>(null)

    const userFormValues = ref({ name: null });

    const userFormRules: FormRules =
    {
        name: {
            validator: (_rule: FormItemRule, value) => {
                if (state.ajaxRunning) return true;
                const localResult = runValidators(value, [required('name'), minLength(3), maxLength(maxNameLength)])
                if (localResult !== true) return localResult
                if (serverErrors.value.name) return new Error(serverErrors.value.name)
                return true
            },
            trigger: ['blur'],
        },
        email: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Email'), validEmail, maxLength(maxEmailLength)])
                if (localResult !== true) return localResult
                if (serverErrors.value.email) return new Error(serverErrors.value.email)
                return true
            },
            trigger: ['blur'],
        },
        password: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Password'), minLength(4)])
                if (localResult !== true) return localResult
                if (serverErrors.value.password) return new Error(serverErrors.value.password)
                return true
            },
            trigger: ['blur']
        }
    };

    const title = computed<string>(() => {
        switch (props.mode) {
            case "add":
                return t("Add user");
            case "update":
                return t("Update user");
            default:
                return "";
        }
    });

    const isSaveDisabled = computed<boolean>(() => {
        return !user.value.name;
    });

    watch(() => userFormValues.value.name, () => {
        delete serverErrors.value.name
    });

    const onSave = async () => {
        serverErrors.value = {};
        try {
            await userFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                await onUpdate()
            }
        } catch (error) {
            console.debug("User form validation error", error)
        }
    };

    const onCancel = () => {
        emit('cancel')
    }


    const onGet = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: GetResponse = await userService.get(id);
            if (response.user.id === id) {
                user.value = new User(response.user);
            } else {
                state.ajaxErrorMessage = t("There was a problem loading the user data");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 404:
                            serverErrors.value.email = t("Email not found");
                            break;
                        case 401:
                            state.ajaxErrors = false;
                            //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                            break;
                        default:
                            state.ajaxErrorMessage = t("API Error: fatal error");
                            break;
                    }
                    userFormRef.value?.restoreValidation();
                    userFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onAdd = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                user: {
                    name: user.value.name,
                    email: user.value.email,
                    isSuperUser: user.value.isSuperUser,
                }
            };
            await userService.add(payload);
            emit('add')
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        /*
                        case 401:
                        state.ajaxErrors = false;
                        //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                        break;
                    case 409:
                        // TODO: conflict (invalid id || name ?)
                        state.ajaxErrorMessage = t("There was a problem adding the user data");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem adding the user data");
                        break;
                            */
                    }
                    //signInFormRef.value?.restoreValidation();
                    //signInFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onUpdate = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                user: {
                    id: user.value.id,
                    name: user.value.name,
                    email: user.value.email,
                    isSuperUser: user.value.isSuperUser,
                }
            };
            await userService.update(payload);
            emit('update')
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        /*
                        case 401:
                        state.ajaxErrors = false;
                        //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                        break;
                    case 409:
                        // TODO: conflict (invalid id || name ?)
                        state.ajaxErrorMessage = t("There was a problem adding the user data");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem adding the user data");
                        break;
                            */
                    }
                    //signInFormRef.value?.restoreValidation();
                    //signInFormRef.value?.validate().then(() => { }).catch(() => { });
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "UserForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };


    onMounted(() => {
        if (props.mode === "update") {
            if (props.userId) {
                onGet(props.userId);
            } else {
                console.error(`TODO: missing userId property for ${props.mode} action`);
            }
        } else if (props.mode === "add") {
            userFormRef.value?.validate();
        }
    });
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            {{ title }}
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectTypeFormRef" :model="user" :rules="userFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input :placeholder="t('Name')" v-model:value="user.name" :maxlength="maxNameLength"
                    :show-count="true" clearable required autofocus></n-input>
            </n-form-item>
            <n-form-item :label="t('Email')" path="email" show-feedback>
                <n-input :placeholder="t('Password')" v-model:value="user.email" :maxlength="maxEmailLength"
                    :show-count="true" clearable required autofocus></n-input>
            </n-form-item>
            <n-form-item :label="t('Password')" path="password" show-feedback>
                <n-input type="password" placeholder="Enter your password" show-password-on="click"
                    :disabled="state.ajaxRunning" ref="inputPasswordRef">
                    <template #password-visible-icon>
                        <n-icon :size="16" :component="IconEyeCancel" />
                    </template>
                    <template #password-invisible-icon>
                        <n-icon :size="16" :component="IconEye" />
                    </template>
                </n-input>
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

<style lang="css" scoped></style>