<script setup lang="ts">
    import { ref } from 'vue';
    import { NTabs, NTabPane, NCard, NAvatar, NFlex, NButton, NButtonGroup, NIcon } from 'naive-ui';
    import { IconFileUpload, IconTrash, IconSun, IconMoon } from '@tabler/icons-vue';

    import { useColorSchemeStore } from '../stores/colorScheme';

    const colorSchemeStore = useColorSchemeStore();

    const currentTab = ref<string>("mySettings");
</script>

<template>
    <n-tabs placement="left" type="line" animated v-model:value="currentTab">
        <n-tab-pane name="myAccount" tab="My account">
            <n-card bordered>
                <h1>My account</h1>
                <h2>Profile details</h2>
                <n-flex style="align-items:center;">
                    <n-avatar :size="128" src="https://i.pravatar.cc/128?u=019dfcaa-aa06-755b-aca0-bb05e4ea44a7" />
                    <n-button>
                        <template #icon>
                            <n-icon>
                                <IconFileUpload />
                            </n-icon>
                        </template>
                        Change avatar
                    </n-button>
                    <n-button>
                        <template #icon>
                            <n-icon>
                                <IconTrash />
                            </n-icon>
                        </template>
                        Delete avatar
                    </n-button>
                </n-flex>
            </n-card>
        </n-tab-pane>
        <n-tab-pane name="mySettings" tab="My settings">
            <n-card bordered>
                <h1>My settings</h1>

                <h2>Locale</h2>

                <h2>Theme</h2>

                <n-button-group>
                    <n-button :disabled="colorSchemeStore.light" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon>
                                <IconSun />
                            </n-icon>
                        </template>
                        light
                    </n-button>
                    <n-button :disabled="colorSchemeStore.dark" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon>
                                <IconMoon />
                            </n-icon>
                        </template>
                        dark
                    </n-button>
                </n-button-group>
                <h2>Layout</h2>
            </n-card>
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped></style>