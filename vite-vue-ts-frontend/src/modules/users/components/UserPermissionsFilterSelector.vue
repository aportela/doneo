<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NSelect } from 'naive-ui';
    import type { SelectMixedOption } from "naive-ui/es/select/src/interface";

    import { UserPermissionFilterValue, type UserPermissionFilter } from "../types/user-admin-permission-filter";

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
    }

    const props = defineProps<Props>();

    const { t } = useI18n();

    const options = computed(() => [
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.any"), value: UserPermissionFilterValue.Any },
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.onlyAdministrators"), value: UserPermissionFilterValue.OnlyAdministrators },
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.onlyUsers"), value: UserPermissionFilterValue.OnlyUsers }
    ] satisfies SelectMixedOption[]);

    const model = defineModel<UserPermissionFilter>('value', { default: UserPermissionFilterValue.Any });

    const selectedPermission = computed({
        get() {
            return model.value;
        },
        set(value: UserPermissionFilter | null) {
            model.value = value ?? UserPermissionFilterValue.Any;
        }
    });

</script>

<template>
    <n-select :clearable="props.clearable" :disabled="props.disabled" :options="options"
        :placeholder="props.placeholder" v-model:value="selectedPermission" />
</template>

<style lang="css" scoped></style>