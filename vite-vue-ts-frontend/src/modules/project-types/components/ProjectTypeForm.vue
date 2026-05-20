<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconUser, IconEdit, IconPlus, IconPalette } from '@tabler/icons-vue';

    import { ProjectType, maxNameLength } from '../models/project-type';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTypeService } from '../services/project-type';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { ProjectTypeResponse, AddRequest, UpdateRequest } from '../types/dto';
    import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';

    interface ProjectTypeFormProps {
        mode: FormMode;
        projectTypeId?: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<ProjectTypeFormProps>();

    const { t } = useI18n();

    const projectType = ref<ProjectType>(
        new ProjectType(
            { id: "", name: "", hexColor: generateRandomSoftHexColor() }
        )
    );

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypeFormRef = ref<FormInst | null>(null)

    const projectTypeFormRules: FormRules =
    {
        name: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("projectTypeFormNameFieldEmptyError"));
                }
                else if (value.length > maxNameLength) {
                    return new Error(t("projectTypeFormNameFieldTooLargeError"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
    };

    watch(() => projectType.value.name, () => { delete serverErrors.value.name });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !projectType.value.name;
    });

    const onSave = async () => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        try {
            await projectTypeFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                await onUpdate()
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "ProjectTypeForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onGet = async (id: string) => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectTypeResponse = await projectTypeService.get(id);
            if (response.id === id) {
                projectType.value = new ProjectType(response);
            } else {
                state.ajaxErrorMessage = t("There was a problem while loading the project type data");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onGet" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified project type");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while loading the project type data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while loading the project type data");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                name: projectType.value.name,
                HexColor: projectType.value.hexColor,
            };
            const addedRole: ProjectTypeResponse = await projectTypeService.add(payload);
            emit('add', addedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onAdd" } });
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "projectTypeFormNameFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while adding the project type data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while adding the project type data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while adding the project type data");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                projectTypeFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        projectTypeFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: projectType.value.id,
                name: projectType.value.name,
                HexColor: projectType.value.hexColor,
            };
            const updatedRole: ProjectTypeResponse = await projectTypeService.update(payload);
            emit('update', updatedRole)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeForm.onUpdate" } });
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "projectTypeFormNameFieldAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("There was a problem while updating the project type data");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while updating the project type data");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while updating the project type data");
                    console.error("Unhandled API error", { file: "ProjectTypeForm.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
                projectTypeFormRef.value?.validate().then(() => { }).catch(() => { });
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectTypeForm.onGet")) {
                if (props.projectTypeId) {
                    onGet(props.projectTypeId);
                } else {
                    console.error(`TODO: missing projectTypeId property for ${props.mode} action`);
                }
            } else if (payload.to.includes("ProjectTypeForm.onAdd")) {
                onAdd();
            } else if (payload.to.includes("ProjectTypeForm.onUpdate")) {
                onUpdate()
            }
        });
        if (props.mode === "update") {
            if (props.projectTypeId) {
                onGet(props.projectTypeId);
            } else {
                console.error(`TODO: missing projectTypeId property for ${props.mode} action`);
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
                {{ t(props.mode == "add" ? "Add project type" : "Update project type") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectTypeFormRef" :model="projectType" :rules="projectTypeFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input type="text" :placeholder="t('projectTypeFormNameFieldPlaceholder')"
                    v-model:value="projectType.name" :maxlength="maxNameLength" :show-count="true" clearable required
                    autofocus>
                    <template #prefix>
                        <n-icon :component="IconUser" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')">
                <n-flex style="width: 100%" align="center" :wrap="false">
                    <n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)" style="width: 100%;">
                        {{ projectType.name }}
                    </n-tag>
                    <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectType.hexColor">
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