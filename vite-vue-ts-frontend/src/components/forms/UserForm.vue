<script setup lang="ts">
    import { ref, reactive, computed, onMounted, watch, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconPalette, IconTrash } from '@tabler/icons-vue';
    import { type AxiosAPIError } from '../../composables/axios';
    import type { EntityAction } from '../../types/common';
    import { type GetProjectTypeResponse } from '../../types/apiResponses';
    import { type UserInterface, UserClass } from '../../types/models/user';
    import { type AjaxStateInterface, defaultAjaxState } from '../../types/ajaxState';
    import { api } from '../../composables/api';
    import { required, minLength, runValidators } from '../../composables/form-validators';
    import { default as RemoteAPIAlert } from '../alerts/RemoteAPIAlert.vue';

    const emit = defineEmits(['add', 'update', 'delete', 'cancel'])

    const { t } = useI18n();

    interface UserFormProps {
        mode: EntityAction;
        userId?: string;
        style?: string | CSSProperties;
    }

    const props = defineProps<UserFormProps>();

    const user = ref<UserInterface>(
        new UserClass({ "id": "", name: "", email: "", isSuperUser: false, createdAt: 0, updatedAt: 0, deletedAt: 0, avatar: "" })
    );

    const addMode = computed<boolean>(() => props.mode === "add");
    const updateMode = computed<boolean>(() => props.mode === "update");
    const deleteMode = computed<boolean>(() => props.mode === "delete");

    const title = computed<string>(() => {
        switch (props.mode) {
            case "add":
                return t("Add user");
            case "update":
                return t("Update user");
            case "delete":
                return t("Delete user");
            default:
                return "";
        }
    });

    const saveButtonDisabled = computed<boolean>(() => {
        const inactiveMode = !(addMode.value || updateMode.value);
        const noName = !user.value.name;
        return inactiveMode || noName;
    });

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const userFormRef = ref<FormInst | null>(null)

    const userFormValues = ref({ name: null });

    const serverErrors = ref<Record<string, string>>({});

    const projectTypeFormRules: FormRules = {
        name: {
            validator: (_rule: FormItemRule, value) => {
                if (state.ajaxRunning) return true;
                const localResult = runValidators(value, [required('name'), minLength(3)])
                if (localResult !== true) return localResult
                if (serverErrors.value.name) return new Error(serverErrors.value.name)
                return true
            },
            trigger: ['blur', 'input'],
        }
    };

    watch(() => userFormValues.value.name, () => {
        delete serverErrors.value.name
    });

    const onSave = async () => {
        serverErrors.value = {};
        try {
            await userFormRef.value?.validate();
            if (addMode.value) {
                onAdd();
            } else if (updateMode.value) {
                onUpdate()
            } else {
                console.error("onSave invalid mode");
            }
        } catch (e) {
            //console.debug("User form validation error", e)
        }
    }

    const onGet = (id: string) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.get(id).then((response: GetProjectTypeResponse) => {
            if (response.data.user.id === id) {
                user.value = response.data.user;
            } else {
                state.ajaxErrorMessage = t("There was a problem loading the user data");
            }
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                        break;
                    case 404:
                        state.ajaxErrorMessage = t("We couldn’t find the specified user");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem loading the user data");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem loading the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onAdd = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.add(user.value).then((_response: any) => {
            emit('add')
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
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
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem adding the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onUpdate = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.update(user.value).then((_response: any) => {
            emit('update')
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                        break;
                    case 409:
                        // TODO: conflict (invalid id || name ?)
                        state.ajaxErrorMessage = t("There was a problem updating the user data");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem updating the user data");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem updating the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onDelete = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.delete(user.value.id).then((_response: any) => {
            emit('delete')
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        //bus.emit("reAuthRequired", { emitter: "DocumentPage.onRefresh" });
                        break;
                    case 404:
                        state.ajaxErrorMessage = t("We couldn’t find the specified user");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem deleting the user data");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem deleting the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onCancel = () => {
        emit('cancel')
    }

    onMounted(() => {
        if (props.mode === "update" || props.mode === "delete") {
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
        <n-form ref="projectTypeFormRef" :model="user" :rules="projectTypeFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input :placeholder="t('Project type name')" v-model:value="user.name"
                    :maxlength="projectTypeMaxNameLength" :show-count="!deleteMode" clearable required autofocus
                    :readonly="deleteMode"></n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')" v-if="!deleteMode">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(user.hexColor)" style="width: 100%;">
                        {{ user.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="user.hexColor">
                        <template #trigger="{ onClick, ref: triggerRef }">
                            <n-button :ref="triggerRef" quaternary @click="onClick">
                                <template #icon>
                                    <IconPalette />
                                </template>
                            </n-button>
                        </template>
                    </n-color-picker>
                </n-flex>
            </n-form-item>
        </n-form>
        <template #footer>
            <RemoteAPIAlert v-if="state.ajaxErrorMessage" type="error" :title="t('Error')"
                :message="state.ajaxErrorMessage" />
        </template>
        <template #action>
            <n-flex>
                <n-button @click="onSave" v-if="addMode || updateMode" :disabled="saveButtonDisabled">
                    <template #icon>
                        <IconDeviceFloppy />
                    </template>
                    {{ t("Save") }}
                </n-button>
                <n-button @click="onDelete" v-else-if="deleteMode" :disabled="state.ajaxRunning">
                    <template #icon>
                        <IconTrash />
                    </template>
                    {{ t("Delete") }}
                </n-button>
                <n-button @click="onCancel" :disabled="state.ajaxRunning">
                    <template #icon>
                        <IconCancel />
                    </template>
                    {{ t("Cancel") }}
                </n-button>
            </n-flex>
        </template>
    </n-card>

</template>

<style lang="css" scoped></style>