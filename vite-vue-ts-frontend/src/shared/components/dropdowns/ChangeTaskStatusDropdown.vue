<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NDropdown, NButton, NIcon, type DropdownOption, type ButtonSize } from 'naive-ui';
    import { IconStatusChange } from '@tabler/icons-vue';

    import type { TaskStatus } from '../../../modules/task-statuses/models/task-status';
    import { useCacheStore } from '../../../stores/cache';

    import { DEFAULT_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface IProps {
        currentStatus: TaskStatus,
        size?: ButtonSize;
        iconSize?: number,
        disabled?: boolean,
        readOnly?: boolean,
    };

    const props = withDefaults(defineProps<IProps>(), {
        size: DEFAULT_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
        readOnly: false,
    });

    const emit = defineEmits(['change']);

    const { t } = useI18n();
    const cacheStore = useCacheStore();

    const options = computed<DropdownOption[]>(() =>
        cacheStore.taskStatuses.map(item => ({
            label: item.name ?? '',
            key: item.id ?? '',
            disabled: item.id === props.currentStatus.id
        }))
    );

    const onChange = (key: string | number) => {
        const status = cacheStore.taskStatuses.find(
            item => item.id === key
        );
        emit("change", status);
    };
</script>

<template>
    <n-dropdown trigger="click" :options="options" @select="onChange" v-if="!props.readOnly">
        <n-button :size="props.size" :disabled="props.disabled">
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconStatusChange" />
            </template>
            {{ t("shared.components.dropDowns.ChangeTaskStatusDropdown.label") }}
        </n-button>
    </n-dropdown>
    <n-button v-else disabled :size="props.size">
        <template #icon>
            <n-icon :size="props.iconSize" :component="IconStatusChange" />
        </template>
        {{ t("shared.components.dropDowns.ChangeTaskStatusDropdown.label") }}
    </n-button>
</template>

<style lang="css" scoped></style>