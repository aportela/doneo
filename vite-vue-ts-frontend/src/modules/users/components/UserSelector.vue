<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, NAvatar, type SelectOption, type SelectInst } from 'naive-ui';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { userService } from '../services/user';
    import type { UserBaseResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { DONEO_ICON_ALERT, DONEO_ICON_CIRCLE_USER } from '../../../shared/types/icons';

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
        autoFocus?: boolean;
        hideAvatar?: boolean;
    }

    const props = defineProps<Props>();

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const userId = defineModel<string | null>('id');

    const avatarURL = computed<string | null>(() => userId.value ? `/api/wc/avatars/user/${userId.value}` : null);

    const options = shallowRef<SelectOption[]>([]);

    const focus = async () => {
        await nextTick();
        selectInstRef.value?.focus();
    };

    defineExpose({ focus });

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.users.length === 0) {
                const response = await userService.searchBase();
                cacheStore.setUsersCache(response.users);
            }
            options.value = cacheStore.users.map((user: UserBaseResponse) => ({ label: user.name, value: user.id }));
        } catch (error: unknown) {
            options.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UserSelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "UserSelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "UserSelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (!state.ajaxErrors && props.autoFocus) {
                focus();
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("UserSelector.onRefresh")) {
                onRefresh();
            }
        });
        onRefresh();
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-input-group>
        <div v-if="!props.hideAvatar">
            <n-avatar v-if="avatarURL" :src="avatarURL" color="transparent" />
            <n-button disabled class="doneo-cursor-default doneo-disable-opacity" v-else>
                <template #icon>
                    <n-icon :component="DONEO_ICON_CIRCLE_USER" />
                </template>
            </n-button>
        </div>
        <n-select filterable ref="selectInstRef" auto :clearable="props.clearable" v-model:value="userId"
            :options="options" :placeholder="props.placeholder" :disabled="isDisabled" />
        <n-button disabled class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="DONEO_ICON_ALERT" />
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>