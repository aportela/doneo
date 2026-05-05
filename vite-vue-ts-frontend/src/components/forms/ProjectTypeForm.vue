<script setup lang="ts">
    import { ref, reactive, computed, onMounted, watch, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules } from 'naive-ui';
    import { v7 as uuidv7 } from 'uuid';
    import { IconCancel, IconDeviceFloppy, IconPalette, IconTrash } from '@tabler/icons-vue';
    import { type AxiosAPIError } from '../../composables/axios';
    import type { EntityAction } from '../../types/common';
    import { type GetProjectTypeResponse } from '../../types/apiResponses';
    import { type ProjectTypeInterface, ProjectTypeClass, maxNameLength as projectTypeMaxNameLength } from '../../types/models/projectType';
    import { type AjaxStateInterface, defaultAjaxState } from '../../types/ajaxState';
    import { api } from '../../composables/api';
    import { required, minLength, runValidators } from '../../composables/form-validators';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../composables/color';
    import { default as RemoteAPIAlert } from '../alerts/RemoteAPIAlert.vue';

    const emit = defineEmits(['add', 'update', 'delete', 'cancel'])

    const { t } = useI18n();

    interface ProjectTypeFormProps {
        mode: EntityAction;
        projectTypeId?: string;
        style?: string | CSSProperties;
    }

    const props = defineProps<ProjectTypeFormProps>();

    const projectType = ref<ProjectTypeInterface>(
        new ProjectTypeClass(
            { id: uuidv7(), name: "", hexColor: generateRandomSoftHexColor() }
        )
    );

    const addMode = computed<boolean>(() => props.mode === "add");
    const updateMode = computed<boolean>(() => props.mode === "update");
    const deleteMode = computed<boolean>(() => props.mode === "delete");

    const title = computed<string>(() => {
        switch (props.mode) {
            case "add":
                return t("Add project type");
            case "update":
                return t("Update project type");
            case "delete":
                return t("Delete project type");
            default:
                return "";
        }
    });

    const saveButtonDisabled = computed<boolean>(() => {
        const inactiveMode = !(addMode.value || updateMode.value);
        const noName = !projectType.value.name;
        return inactiveMode || noName;
    });

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypeFormRef = ref<FormInst | null>(null)

    const projectTypeFormValues = ref({ name: null });

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

    watch(() => projectTypeFormValues.value.name, () => {
        delete serverErrors.value.name
    });

    const onSave = async () => {
        serverErrors.value = {};
        try {
            await projectTypeFormRef.value?.validate();
            if (addMode.value) {
                onAdd();
            } else if (updateMode.value) {
                onUpdate()
            } else {
                console.error("TODO");
            }
        } catch (e) {
            //console.debug("ProjectType form validation error", e)
        }
    }

    const onGet = (id: string) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.projectTypes.get(id).then((response: GetProjectTypeResponse) => {
            if (response.data.projectType.id === id) {
                projectType.value = response.data.projectType;
            } else {
                state.ajaxErrorMessage = t("Unable to load project type: invalid API response data");
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
                        state.ajaxErrorMessage = "Unable to load project type: project type not found";
                        break;
                    default:
                        state.ajaxErrorMessage = t("Unable to load project type: invalid API response code");
                        break;
                }
            } else {
                state.ajaxErrorMessage = `Unable to load project type: ${errorResponse} `;
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onAdd = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.projectTypes.add(projectType.value).then((_response: any) => {
            emit('add')
        }).catch((_error: AxiosAPIError) => {
            // TODO:
            console.error(_error);
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onUpdate = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.projectTypes.update(projectType.value).then((_response: any) => {
            emit('update')
        }).catch((_error: AxiosAPIError) => {
            // TODO:
            console.error(_error);
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onDelete = () => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.projectTypes.delete(projectType.value.id).then((_response: any) => {
            emit('delete')
        }).catch((_error: AxiosAPIError) => {
            // TODO:
            console.error(_error);
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onCancel = () => {
        emit('cancel')
    }

    onMounted(() => {
        if (props.mode === "update" || props.mode === "delete") {
            if (props.projectTypeId) {
                onGet(props.projectTypeId);
            } else {
                console.error(`TODO: missing projectTypeId prop for ${props.mode} action`);
            }
        } else if (props.mode === "add") {
            projectTypeFormRef.value?.validate();
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
        <n-form ref="projectTypeFormRef" :model="projectType" :rules="projectTypeFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input :placeholder="t('Project type name')" v-model:value="projectType.name"
                    :maxlength="projectTypeMaxNameLength" :show-count="!deleteMode" clearable required autofocus
                    :readonly="deleteMode"></n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')" v-if="!deleteMode">
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