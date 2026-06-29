<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, NAvatar, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';
    import { IconAlertCircle, IconUserCircle } from '@tabler/icons-vue';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { userService } from '../services/user';
    import type { UserBaseResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    interface UserSelectorProps {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        hideAvatar?: boolean;
        disabled?: boolean;
    }

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const userId = defineModel<string | null>('id');

    const avatarURL = computed(() => userId.value ? `/api/wc/avatars/user/${userId.value}/tiny` : null);

    const props = defineProps<UserSelectorProps>();

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.users.length > 0) {
                options.value = cacheStore.users.map((user: UserBaseResponse) => ({ label: user.name, value: user.id }));
            } else {
                const response = await userService.searchBase();
                options.value = response.users.map((user: UserBaseResponse) => ({ label: user.name, value: user.id }));
                cacheStore.setUsersCache(response.users)
            }
            if (props.autoFocus) {
                focus();
            }
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
        }
    };

    const focus = () => {
        nextTick(() => {
            selectInstRef.value?.focus();
        });
    };

    const reset = () => {
        userId.value = null;
    }

    defineExpose({ reset });

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
            <n-avatar v-if="avatarURL" :src="avatarURL" />
            <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity" v-else>
                <template #icon>
                    <n-icon :component="IconUserCircle">
                    </n-icon>
                </template>
            </n-button>
        </div>
        <n-select filterable ref="selectInstRef" :required="props.required" auto :clearable="props.clearable"
            v-model:value="userId" :options="options" :placeholder="props.placeholder" :size="props.size"
            :disabled="isDisabled" />
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="IconAlertCircle">
                </n-icon>
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>