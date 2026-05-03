<script setup lang="ts">
    import { ref, reactive, computed, onMounted, watch, type CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NSpin, NCard, NInput, NFlex, NButton, NColorPicker, NTag, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules } from 'naive-ui';
    import { v7 as uuidv7 } from 'uuid';
    import { IconCancel, IconDeviceFloppy, IconTrash } from '@tabler/icons-vue';
    import { type AxiosAPIError } from '../../composables/axios';
    import type { EntityAction } from '../../types/common';
    import { type ProjectTypeInterface, ProjectTypeClass, maxNameLength as projectTypeMaxNameLength } from '../../types/models/projectType';
    import { type AjaxStateInterface, defaultAjaxState } from '../../types/ajaxState';
    import { api } from '../../composables/api';
    import { required, minLength, runValidators } from '../../composables/form-validators';
    import { generateRandomSoftHexColor, getNaiveUITagColorProperty } from '../../composables/color';

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
        state.ajaxRunning = true;
        api.projectTypes.get(id).then((_response: any) => {
            projectType.value = _response.data.projectType;
        }).catch((_error: AxiosAPIError) => {
            // TODO:
            console.error(_error);
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onAdd = () => {
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
        if (props.mode == "update" || props.mode == "delete") {
            if (props.projectTypeId) {
                onGet(props.projectTypeId);
            } else {
                console.error(`TODO: missing projectTypeId prop for ${props.mode} action`);
            }
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
        <n-form ref="projectTypeFormRef" :model="projectType" :rules="projectTypeFormRules">
            <n-form-item :label="t('Name')" path="name" show-feedback>
                <n-input :placeholder="t('Project type name')" v-model:value="projectType.name"
                    :disabled="state.ajaxRunning || deleteMode" :maxlength="projectTypeMaxNameLength"
                    :show-count="!deleteMode" clearable required autofocus></n-input>
            </n-form-item>
            <n-form-item :label="t('Preview')" v-if="!deleteMode">
                <n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)">
                    {{ projectType.name }}
                </n-tag>
            </n-form-item>
            <n-form-item :label="t('Color')" v-if="!deleteMode">
                <n-color-picker :modes="['hex']" :show-alpha="false" v-model:value="projectType.hexColor" />
            </n-form-item>
        </n-form>
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