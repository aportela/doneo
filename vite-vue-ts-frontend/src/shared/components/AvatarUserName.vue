<script setup lang="ts">
    import { computed } from 'vue';

    import { NAvatar } from 'naive-ui';

    interface AvatarUserNameProps {
        userId: string | null;
        userName: string | null;
        avatarSize?: number;
    };

    const props = withDefaults(defineProps<AvatarUserNameProps>(), {
        avatarSize: 32,
    });

    const avatarURL = computed(() => `/api/avatars/${props.avatarSize}/user/${props.userId}`);
</script>

<template>
    <div class="doneo-flex-center-align" style="gap: 8px;">
        <n-avatar v-if="props.userId" :src="avatarURL" class="doneo-avatar-username" :size="props.avatarSize" />
        {{ props.userName }}
    </div>
</template>

<style lang="css" scoped>
    .doneo-avatar-username {
        margin-right: 4px;
        flex-shrink: 0;
    }
</style>