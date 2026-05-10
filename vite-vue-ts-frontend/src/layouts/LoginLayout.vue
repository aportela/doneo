<script setup lang="ts">
    import { ref, onMounted } from "vue";
    import { useRoute, useRouter } from 'vue-router';
    import { useSessionStore } from "../stores/session";
    import { TokenManager } from '../api/services/tokenManager';

    const router = useRouter();
    const route = useRoute();
    const sessionStore = useSessionStore();

    const showRouterView = ref<boolean>(false);
    onMounted(async () => {
        if (!sessionStore.hasAccessToken) {
            try {
                const success = await TokenManager.refreshAccessToken(sessionStore);
                if (success) {
                    const redirectPath = (route.query.redirect as string) || '/home';
                    await router.push(redirectPath);
                }
            } catch (error: unknown) {
                console.error("An unhandled exception occurred during access token refresh", error);
            } finally {
                showRouterView.value = true;
            }
        }
    });
</script>

<template>
    <router-view v-if="showRouterView" />
</template>