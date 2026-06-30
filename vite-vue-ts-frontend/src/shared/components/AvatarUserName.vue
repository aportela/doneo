<script setup lang="ts">
    import { computed } from 'vue';

    import { NAvatar } from 'naive-ui';

    type AvatarSize = "tiny" | "small" | "normal";

    interface AvatarUserNameProps {
        userId: string | null;
        userName: string | null;
        avatarSize?: AvatarSize;
    };

    const props = withDefaults(defineProps<AvatarUserNameProps>(), {
        avatarSize: "tiny",
    });

    const size = computed(() => {
        switch (props.avatarSize) {
            case "tiny":
                return 32;
            case "small":
                return 64;
            case "normal":
                return 128;
        }
    });

    const avatarURL = computed(() => `/api/wc/avatars/user/${props.userId}`);
</script>

<template>
    <div class="doneo-flex-center-align" style="gap: 8px;">
        <n-avatar v-if="props.userId" :src="avatarURL" class="doneo-avatar-username" :size="size" />
        {{ props.userName }}
    </div>
</template>

<style lang="css" scoped>
    .doneo-avatar-username {
        margin-right: 4px;
        flex-shrink: 0;
    }
</style>