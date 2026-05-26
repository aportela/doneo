<script setup lang="ts">
    import type { CSSProperties } from 'vue';

    import { NCard, NForm, NFormItem, NInput } from 'naive-ui';

    import type { FormMode } from '../../../shared/types/form-mode';
    import { Project, maxKeyLength, maxSummaryLength } from "../models/project";
    import ProjectPrioritySelector from "../../project-priorities/components/ProjectPrioritySelector.vue";
    import ProjectStatusSelector from "../../project-statuses/components/ProjectStatusSelector.vue";
    import ProjectTypeSelector from "../../project-types/components/ProjectTypeSelector.vue";
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface ProjectFormProps {
        mode: FormMode;
        style?: string | CSSProperties;
    }

    const props = defineProps<ProjectFormProps>();

    const project = defineModel<Project>("project", { required: true });

</script>

<template>
    <n-card bordered :style="props.style">
        <n-form-item label="Created by">
            <div class="note-user">
                <AvatarUserName :user="project.createdBy" />
            </div>
        </n-form-item>
        <n-form-item label="Created at">
            {{ project.createdAt.toLocaleString() }}
        </n-form-item>
        <n-form>
            <n-form-item label="Key">
                <n-input v-model:value="project.key" :show-count="true" :maxlength="maxKeyLength" />
            </n-form-item>
            <n-form-item label="Summary">
                <n-input v-model:value="project.summary" :show-count="true" :maxlength="maxSummaryLength" />
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
    </n-card>
</template>

<style lang="css" scoped></style>