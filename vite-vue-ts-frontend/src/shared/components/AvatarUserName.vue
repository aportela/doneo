<script setup lang="ts">
    import { NAvatar } from 'naive-ui';
    import type { UserBase } from '../../modules/users/models/user';

    interface AvatarUserNameProps {
        user: UserBase;
        avatarSize?: number;
    };

    const props = withDefaults(defineProps<AvatarUserNameProps>(), {
        avatarSize: 32,
    });
</script>

<template>
    <div class="doneo-flex-center-align" style="gap: 8px;">
        <n-avatar v-if="user.avatarUrl" :src="user.avatarUrl" class="doneo-avatar-username" :size="props.avatarSize" />
        {{ user.name }}
    </div>
</template>

<style lang="css" scoped>
    .doneo-avatar-username {
        margin-right: 4px;
    }
</style>