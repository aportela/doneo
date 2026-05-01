<script setup lang="ts">
    import { ref, computed, defineComponent, h, onMounted, shallowRef } from 'vue'
    import { NLayout, NLayoutSider, NLayoutContent, NMenu, NCard, NTable, NColorPicker, NTag } from 'naive-ui'
    import type { MenuOption } from 'naive-ui'
    import { api } from '../composables/api';


    interface ProjectTypeInterface {
        id: string;
        name: string;
    }

    class ProjectType implements ProjectTypeInterface {
        id: string;
        name: string;
        constructor(item: ProjectTypeInterface) {
            this.id = item.id;
            this.name = item.name;
        }
    }

    interface ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    class ProjectStatus implements ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
        constructor(item: ProjectStatusInterface) {
            this.id = item.id;
            this.name = item.name;
            this.index = item.index;
            this.hexColor = item.hexColor;
        }
    }

    interface ProjectPriorityInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }


    const ProjectTypes = defineComponent({
        setup() {
            return () =>
                h('div', null, 'Project types')
        }
    })

    const ProjectPriorities = defineComponent({
        setup() {
            return () =>
                h('div', null, 'Project priorities')
        }
    })

    const ProjectStatuses = defineComponent({
        setup() {
            return () =>
                h('div', null, 'Project statuses')
        }
    })

    const componentMap = {
        projectTypes: ProjectTypes,
        projectPriorities: ProjectPriorities,
        projectStatuses: ProjectStatuses
    }

    type SettingsKey = keyof typeof componentMap

    const selectedKey = ref<SettingsKey>('projectTypes')

    const menuOptions: MenuOption[] = [
        { label: 'Types', key: 'projectTypes' },
        { label: 'Priorities', key: 'projectPriorities' },
        { label: 'Statuses', key: 'projectStatuses' }
    ]

    const currentComponent = computed(() => componentMap[selectedKey.value])


    const projectTypes = shallowRef<ProjectType[]>([]);
    const projectStatuses = shallowRef<ProjectStatus[]>([]);
    const projectPriorities = shallowRef<ProjectPriorityInterface[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.projectTypes.search().then((successResponse: any) => {
            projectTypes.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
        api.projectStatuses.search().then((successResponse: any) => {
            projectStatuses.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
        api.projectPriorities.search().then((successResponse: any) => {
            projectPriorities.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

    const hexToRgba = (hex: string, alphaOverride?: number) => {
        if (!hex) return `rgba(0,0,0,1)`

        let h = hex.replace('#', '')

        let r, g, b, a = 1

        if (h.length === 8) {
            // RRGGBBAA
            r = parseInt(h.slice(0, 2), 16)
            g = parseInt(h.slice(2, 4), 16)
            b = parseInt(h.slice(4, 6), 16)
            a = parseInt(h.slice(6, 8), 16) / 255
        } else {
            // RRGGBB
            r = parseInt(h.slice(0, 2), 16)
            g = parseInt(h.slice(2, 4), 16)
            b = parseInt(h.slice(4, 6), 16)
        }

        const alpha = alphaOverride ?? a

        return `rgba(${r}, ${g}, ${b}, ${alpha})`
    }

    const tagColor = (base: string) => {
        console.log(base);
        console.log(hexToRgba(base, 1));
        console.log({
            color: hexToRgba(base, 0.2),
            textColor: hexToRgba(base, 1),
            borderColor: hexToRgba(base, 0.5)
        });
        return {
            color: hexToRgba(base, 0.2),
            textColor: hexToRgba(base, 1),
            borderColor: hexToRgba(base, 0.5)
        }
    }

    const color = ref<string>("#000000");
</script>

<template>

    <n-card class="settings-card" content-style="padding: 0;">
        <n-layout has-sider class="inner-layout">
            <n-layout-sider width="220" class="settings-sider" content-style="padding: 12px;">
                <h2>
                    Project settings
                </h2>
                <n-menu v-model:value="selectedKey" :options="menuOptions" />
            </n-layout-sider>
            <n-layout-content class="settings-content">
                <div class="content-inner">
                    <n-table>
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Color</th>
                                <th>Preview</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="projectType in projectTypes" :key="projectType.id">
                                <td>{{ projectType.name }}</td>
                                <td><n-color-picker v-model:value="color" :modes="['hex']"
                                        :show-alpha="false"></n-color-picker></td>
                                <td><n-tag :color="tagColor(color)">{{ projectType.name }}</n-tag></td>
                            </tr>
                        </tbody>
                    </n-table>
                    <component :is="currentComponent" />
                </div>
            </n-layout-content>
        </n-layout>
    </n-card>
</template>

<style scoped>
    .settings-card {
        border-radius: 8px;
    }

    .inner-layout {
        background: transparent;
    }

    .settings-sider {
        background: transparent;
        border-right: 1px solid #eee;
    }

    .settings-content {
        background: transparent;
    }

    .content-inner {
        padding: 24px;
    }
</style>