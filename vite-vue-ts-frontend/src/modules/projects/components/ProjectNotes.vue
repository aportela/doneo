<script setup lang="ts">
    import type { CSSProperties } from "vue";

    import { NSpace, NCard } from "naive-ui";

    import { UserBase } from "../../users/models/user.ts";

    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";

    interface ProjectNotesProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectNotesProps>();

    const user = new UserBase({ id: '019e5b68-f701-7af7-a694-c8146e3366a1', name: 'John Doe', avatarUrl: 'https://i.pravatar.cc/32?u=019e5b68-f701-7af7-a694-c8146e3366a1' });

</script>

<template>
    <n-card bordered :style="props.style">
        <n-space vertical size="large" style="margin-right: 12px;">
            <n-card v-for="index in [1, 2, 3, 4, 5]" :key="index" size="small" bordered>
                <div class="note-header">
                    <div class="note-user">
                        <AvatarUserName :user="user" />
                    </div>

                    <span class="note-date">
                        12/02/2026
                    </span>
                </div>

                <div class="note-content">

                    random note text {{ index }}

                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec fringilla pellentesque
                    bibendum.
                    Nunc efficitur est eget magna consectetur euismod. Nam eget ante arcu. Vestibulum ex
                    enim,
                    consectetur sit amet fringilla hendrerit, tristique vel risus. Nunc eu magna at magna
                    dictum
                    posuere ac id lectus. Fusce sed euismod elit. Vivamus porta diam quis bibendum tempor.
                    Mauris at
                    mauris laoreet, porta nunc quis, suscipit orci. Donec eu ex nec tellus luctus tincidunt.
                    Quisque
                    malesuada sollicitudin luctus.

                    Donec rhoncus hendrerit ante id tincidunt. Etiam quis scelerisque tellus. Vivamus sit
                    amet massa
                    finibus, commodo erat sit amet, porttitor quam. Integer pellentesque urna vitae risus
                    mattis
                    elementum. Nam ornare ornare enim, quis rutrum nunc ornare et. Donec quis ex at magna
                    facilisis
                    cursus. Duis et libero sollicitudin, aliquet nunc porta, tristique nulla. Nam eget nisl
                    mauris.
                    Aliquam quis ultrices ante. Nullam ullamcorper tortor consequat metus auctor, vel
                    sagittis mi
                    maximus. Donec ut aliquet neque. Ut malesuada lorem hendrerit nulla scelerisque
                    volutpat.
                    Pellentesque commodo pellentesque tortor, at pellentesque enim feugiat lacinia. Integer
                    ac nulla
                    pretium ex eleifend feugiat id et massa. Donec ex mauris, sodales at justo posuere,
                    mattis
                    rutrum nisi.

                </div>
            </n-card>
        </n-space>
    </n-card>
</template>

<style lang="css" scoped>
    .note-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
    }

    .note-user {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 500;
    }

    .note-date {
        font-size: 12px;
        color: #999;
    }

    .note-content {
        font-size: 14px;
        white-space: pre-line;
    }


</style>