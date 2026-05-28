<script setup lang="ts">
    import type { CSSProperties } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NForm, NFormItem, NInput, NButton, NIcon } from 'naive-ui';

    import type { FormMode } from '../../../shared/types/form-mode';
    import { Project, MAX_KEY_LENGTH, MAX_SUMMARY_LENGTH } from "../models/project";
    import ProjectPrioritySelector from "../../project-priorities/components/ProjectPrioritySelector.vue";
    import ProjectStatusSelector from "../../project-statuses/components/ProjectStatusSelector.vue";
    import ProjectTypeSelector from "../../project-types/components/ProjectTypeSelector.vue";
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { IconDeviceFloppy } from '@tabler/icons-vue';

    interface ProjectFormProps {
        mode: FormMode;
        style?: string | CSSProperties;
    }

    const props = defineProps<ProjectFormProps>();

    const emit = defineEmits(["save"]);

    const project = defineModel<Project>("project", { required: true });

    const { t } = useI18n();

    const onSave = () => {
        emit("save");
    };

</script>

<template>
    <n-card bordered :style="props.style">
        <n-form-item label="Created by">
            <div class="note-user">
                <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
            </div>
        </n-form-item>
        <n-form-item label="Created at">
            {{ project.createdAt.toLocaleString() }}
        </n-form-item>
        <n-form>
            <n-form-item label="Key">
                <n-input v-model:value="project.key" :show-count="true" :maxlength="MAX_KEY_LENGTH" />
            </n-form-item>
            <n-form-item label="Summary">
                <n-input v-model:value="project.summary" :show-count="true" :maxlength="MAX_SUMMARY_LENGTH" />
            </n-form-item>
            <n-form-item label="Description">
                <n-input v-model:value="project.description" type="textarea" clearable />
            </n-form-item>
            <n-form-item label="Type">
                <ProjectTypeSelector v-model:id="project.type.id" />
            </n-form-item>
            <n-form-item label="Priority">
                <ProjectPrioritySelector v-model:id="project.priority.id" />
            </n-form-item>
            <n-form-item label="Status">
                <ProjectStatusSelector v-model:id="project.status.id" />
            </n-form-item>
        </n-form>
        <n-button @click="onSave">
            <template #icon>
                <n-icon :component="IconDeviceFloppy"></n-icon>
            </template>
            {{ t("shared.buttons.Save.label") }}
        </n-button>
    </n-card>
</template>

<style lang="css" scoped></style>