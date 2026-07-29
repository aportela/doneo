<script setup lang="ts">
    import { shallowRef, reactive, computed, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, type SelectOption } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { roleService } from '../services/role';
    import type { RoleBaseResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { DONEO_ICON_ALERT } from '../../../shared/types/icons';

    // TODO: add cache ???

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
        autoFocus?: boolean;
    }

    const props = defineProps<Props>();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const roleId = defineModel<string | null>('id');

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response = await roleService.searchBase();
            options.value = response.roles.map((role: RoleBaseResponse) => ({ label: role.name, value: role.id }));
        } catch (error: unknown) {
            options.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RoleSelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "RoleSelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "RoleSelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("RoleSelector.onRefresh")) {
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
        <n-select filterable :clearable="props.clearable" v-model:value="roleId" :options="options"
            :placeholder="props.placeholder" :disabled="isDisabled" />
        <n-button disabled class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="DONEO_ICON_ALERT" />
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>