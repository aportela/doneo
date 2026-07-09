<script setup lang="ts">
    import { computed } from 'vue';

    import { NAvatar } from 'naive-ui';

    interface IProps {
        userId: string | null;
        userName: string | null;
        avatarSize?: "tiny" | "small" | "normal";
    };

    const props = withDefaults(defineProps<IProps>(), {
        avatarSize: "tiny",
    });

    const avatarPixelSize = computed(() => {
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
        <slot name="before"></slot>
        <n-avatar v-if="props.userId" :src="avatarURL" class="doneo-avatar-username" :size="avatarPixelSize"
            color="transparent" />
        {{ props.userName }}
        <slot name="after"></slot>
    </div>
</template>

<style lang="css" scoped>
    .doneo-avatar-username {
        margin-right: 4px;
        flex-shrink: 0;
    }
</style>