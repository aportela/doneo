<script setup lang="ts">
    import { ref, onMounted } from 'vue';
    import { useRoute } from 'vue-router';

    import { NIcon, NButton, NInputGroup, NInput, NFlex } from 'naive-ui';

    import ToggleMarkDownEditor from '../../../shared/components/form-blocks/ToggleMarkDownEditor.vue';

    import { pageService } from '../services/page.ts';
    import { Page } from '../models/page.ts';
    import { IconDeviceFloppy, IconLock, IconLockOpen2 } from '@tabler/icons-vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import type { UpdateRequest } from '../types/dto.ts';

    const route = useRoute();

    const projectId = route.params.projectId as string
    const pageId = route.params.pageId as string

    const page = ref<Page>(new Page());

    const locked = ref<boolean>(true);

    const onRefresh = async () => {
        const response = await pageService.getProjectPage(projectId, pageId);
        page.value = new Page(response);
    };

    const onSave = async () => {
        const payload: UpdateRequest = {
            id: page.value.id ?? "",
            title: page.value.title,
            body: page.value.body,
        };
        const response = await pageService.updateProjectPage(projectId, pageId, payload);
        page.value = new Page(response);
        locked.value = true;
    };

    onMounted(() => {
        onRefresh();
    });
</script>

<template>

    <n-input-group>
        <n-input v-model:value="page.title" :readonly="locked" size="small" />
        <n-button @click="locked = !locked" size="small">
            <template #icon>
                <n-icon :component="locked ? IconLock : IconLockOpen2" />
            </template>
        </n-button>
        <n-button v-if="!locked" size="small" @click="onSave">
            <template #icon>
                <n-icon :component="IconDeviceFloppy" />
            </template>
        </n-button>
    </n-input-group>
    <n-flex align="center" style="margin-top: 16px;">
        <AvatarUserName :user-id="page.createdBy.id" :user-name="page.createdBy.name" />
        on {{ page.createdAt?.toLocaleString() }}
    </n-flex>
    <div class="body-container" :class="{ 'body-container-readonly': locked }">
        <ToggleMarkDownEditor :read-only="locked" v-model:value="page.body" />
    </div>
</template>

<style lang="css" scoped>

    .body-container {
        margin-top: 24px;
    }

    .body-container-readonly {
        border: 1px solid #ccc;
        border-radius: 4px;
        padding: 4px 16px;
    }
</style>