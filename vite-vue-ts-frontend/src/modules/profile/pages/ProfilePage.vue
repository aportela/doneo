<script setup lang="ts">
    import { computed, ref, reactive, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import dayjs from 'dayjs'

    import { NForm, NAvatar, NFlex, NFormItem, NInput, NButton, NButtonGroup, NPopover, NIcon, NGrid, NGridItem, NDivider, NTag, type FormInst, type FormRules, type FormItemRule } from 'naive-ui';
    import { IconTrash, IconSun, IconMoon, IconInfoCircle, IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand, IconDeviceFloppy, IconImageGeneration } from '@tabler/icons-vue';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { profileService } from '../services/profile';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { useLoadingStore } from '../../../stores/loading';
    import { useColorSchemeStore } from '../../../stores/colorScheme';
    import { useUserSettingsStore } from '../../../stores/userSettings';
    import { Profile } from '../models/profile';
    import type { ProfileResponse, UpdateRequest } from '../types/dto';
    import { MAX_EMAIL_LENGTH, MAX_NAME_LENGTH, MIN_PASSWORD_LENGTH } from '../../users/models/user';
    import { defaultDateTimeMask } from '../../../shared/composables/datetime.ts';

    import GenerateAvatarModal from '../../../shared/components/modals/GenerateAvatarModal.vue';
    import { isValidEmail } from '../../../shared/composables/form-validators.ts';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

    const profile = ref<Profile>(new Profile());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const profileFormRef = ref<FormInst | null>(null);

    const newPassword = ref<string | null>(null);
    const confirmedPassword = ref<string | null>(null);


    const profileFormRules: FormRules =
    {
        name: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length > MAX_NAME_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
        email: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (!isValidEmail(value)) {
                    return new Error(t("shared.warningMessages.fieldHasInvalidFormat"));
                }
                else if (value.length > MAX_EMAIL_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.email) {
                    return new Error(t(serverErrors.value.email));
                } else {
                    return true;
                }
            }, trigger: ['blur'],
        },
    };
    const lastAvatarTimestamp = ref<number>(new Date().getTime());

    const currentAvatarURL = computed(() => profile.value.id ? `/api/wc/avatars/user/${profile.value.id}?t=${lastAvatarTimestamp.value}` : undefined);

    const showAvatarGeneratorModal = ref<boolean>(false);

    const currentDatetimeMask = ref<string | null>(userSettingsStore.currentDatetimeMask);

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => currentDatetimeMask.value, (newValue) => {
        userSettingsStore.setDatetimeMask(newValue || defaultDateTimeMask)
    });

    const validNewPassword = computed<boolean>(() => {
        if (newPassword.value) {
            return newPassword.value.length >= MIN_PASSWORD_LENGTH;
        } else {
            return true;
        }
    });

    const matchedPasswords = computed<boolean>(() => newPassword.value === confirmedPassword.value);

    const currentDatetimeMaskPreview = computed<string | null>({
        get() {
            return currentDatetimeMask.value ? dayjs().format(currentDatetimeMask.value) : null;
        },
        set(_value) { }
    });

    const allowSubmit = computed<boolean>(() => !!profile.value.name && !!profile.value.email && matchedPasswords.value && validNewPassword.value);

    const onGet = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProfileResponse = await profileService.get();
            profile.value = new Profile(response);
            lastAvatarTimestamp.value = new Date().getTime();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProfilePage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProfilePage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        profileFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                name: profile.value.name ?? "",
                password: newPassword.value || undefined,
                email: profile.value.email ?? "",
            };
            const response: ProfileResponse = await profileService.update(payload);
            profile.value = new Profile(response);
            notify('success', t("modules.profile.components.ProfilePage.notifications.profileUpdated"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.profile.components.ProfilePage.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                serverErrors.value.email = "modules.profile.components.ProfilePage.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    profileFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onSaveAvatar = async (svg: string) => {
        showAvatarGeneratorModal.value = false;
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload = { svg: svg };
            await profileService.saveAvatar(payload);
            notify('success', t("modules.profile.components.ProfilePage.notifications.avatarUpdated"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onSaveAvatar" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateAvatarError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateAvatarError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onSaveAvatar" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            lastAvatarTimestamp.value = new Date().getTime();
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    //userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onDeleteAvatar = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await profileService.deleteAvatar();
            notify('success', t("modules.profile.components.ProfilePage.notifications.avatarDeleted"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onDeleteAvatar" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                //serverErrors.value.name = "modules.user.components.UserForm.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                //serverErrors.value.email = "modules.user.components.UserForm.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onDeleteAvatar" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            lastAvatarTimestamp.value = new Date().getTime();
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    //userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onGet();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProfilePage.onGet")) {
                onGet();
            } else if (payload.to.includes("ProfilePage.onUpdate")) {
                onUpdate();
            } else if (payload.to.includes("ProfilePage.onDeleteAvatar")) {
                onDeleteAvatar();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <GenerateAvatarModal v-if="showAvatarGeneratorModal" v-model:show="showAvatarGeneratorModal"
        @confirm="(svg: string) => onSaveAvatar(svg)" @cancel="showAvatarGeneratorModal = false;" />
    <h1>Profile</h1>

    <n-flex align="center" justify="space-between">
        <h2>My account</h2>
        <div>
            <n-tag type="success">Created on {{
                profile.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</n-tag>
            <n-tag type="info" v-if="profile.updatedAt?.hasValue()">Last update on {{
                profile.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</n-tag>
        </div>

    </n-flex>
    <n-flex style="align-items:center;">
        <n-avatar :size="128" :src="currentAvatarURL" :key="lastAvatarTimestamp" color="transparent" />


        <n-button @click="showAvatarGeneratorModal = true">
            <template #icon>
                <n-icon :component="IconImageGeneration" />
            </template>
            Change avatar
        </n-button>
        <n-button tertiary type="error" @click="onDeleteAvatar">
            <template #icon>
                <n-icon :component="IconTrash" />
            </template>
            Delete avatar
        </n-button>
    </n-flex>

    <n-divider />

    <n-form ref="profileFormRef" :model="profile" :rules="state.ajaxRunning ? {} : profileFormRules"
        :disabled="state.ajaxRunning">
        <n-flex align="center">
            <n-form-item label="Name" path="name" show-feedback>
                <n-input v-model:value="profile.name" />
            </n-form-item>
            <n-form-item label="Email" path="email" show-feedback>
                <n-input v-model:value="profile.email" />
            </n-form-item>
            <n-form-item label="New password" :feedback="!validNewPassword ? 'min password length 4 chars' : undefined"
                :validation-status="!validNewPassword ? 'error' : undefined">
                <n-input type="password" :min="MIN_PASSWORD_LENGTH" v-model:value="newPassword"
                    placeholder="type new password" clearable />
            </n-form-item>
            <n-form-item label="Confirm new password"
                :feedback="!matchedPasswords ? 'passwords do not match' : undefined"
                :validation-status="!matchedPasswords ? 'error' : undefined">
                <n-input type="password" :min="MIN_PASSWORD_LENGTH" v-model:value="confirmedPassword"
                    placeholder="confirm new password" clearable />
            </n-form-item>
            <n-button @click="onUpdate" :disabled="!allowSubmit">
                <template #icon>
                    <n-icon :component="IconDeviceFloppy" />
                </template>
                {{ t("shared.buttons.Save.label") }}
            </n-button>
        </n-flex>

        <n-divider />

        <h2>Settings</h2>

        <n-flex align="center">
            <div>
                <h3>Theme</h3>
                <n-button-group>
                    <n-button :disabled="colorSchemeStore.light" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon :component="IconSun" />
                        </template>
                        light
                    </n-button>
                    <n-button :disabled="colorSchemeStore.dark" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon :component="IconMoon" />
                        </template>
                        dark
                    </n-button>
                </n-button-group>
            </div>
            <div>
                <h3>Layout</h3>
                <n-button-group>
                    <n-button :disabled="userSettingsStore.sideNavigationMode"
                        @click="userSettingsStore.toggleNavigationMode">
                        <template #icon>
                            <n-icon :component="IconLayoutSidebarLeftExpand" />
                        </template>
                        side
                    </n-button>
                    <n-button :disabled="userSettingsStore.topNavigationMode"
                        @click="userSettingsStore.toggleNavigationMode">
                        <template #icon>
                            <n-icon :component="IconLayoutNavbarExpand" />
                        </template>
                        top
                    </n-button>
                </n-button-group>
            </div>

            <div>
                <h3 class="doneo-flex-center-align">Locale
                    <n-popover>
                        <template #trigger>
                            <n-icon style="margin-left: 8px; margin-top: 2px;" class="doneo-cursor-help" :size="20"
                                :component="IconInfoCircle" />
                        </template>
                        <template #header>
                            List of all available formats
                        </template>
                        <n-grid :cols="2" :x-gap="128">
                            <n-grid-item>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Format</th>
                                            <th>Output</th>
                                            <th>Description</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td><code>YY</code></td>
                                            <td>18</td>
                                            <td>Two-digit year</td>
                                        </tr>
                                        <tr>
                                            <td><code>YYYY</code></td>
                                            <td>2018</td>
                                            <td>Four-digit year</td>
                                        </tr>
                                        <tr>
                                            <td><code>M</code></td>
                                            <td>1-12</td>
                                            <td>The month, beginning at 1</td>
                                        </tr>
                                        <tr>
                                            <td><code>MM</code></td>
                                            <td>01-12</td>
                                            <td>The month, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>MMM</code></td>
                                            <td>Jan-Dec</td>
                                            <td>The abbreviated month name</td>
                                        </tr>
                                        <tr>
                                            <td><code>MMMM</code></td>
                                            <td>January-December</td>
                                            <td>The full month name</td>
                                        </tr>
                                        <tr>
                                            <td><code>D</code></td>
                                            <td>1-31</td>
                                            <td>The day of the month</td>
                                        </tr>
                                        <tr>
                                            <td><code>DD</code></td>
                                            <td>01-31</td>
                                            <td>The day of the month, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>d</code></td>
                                            <td>0-6</td>
                                            <td>The day of the week, with Sunday as 0</td>
                                        </tr>
                                        <tr>
                                            <td><code>dd</code></td>
                                            <td>Su-Sa</td>
                                            <td>The min name of the day of the week</td>
                                        </tr>
                                        <tr>
                                            <td><code>ddd</code></td>
                                            <td>Sun-Sat</td>
                                            <td>The short name of the day of the week</td>
                                        </tr>
                                        <tr>
                                            <td><code>dddd</code></td>
                                            <td>Sunday-Saturday</td>
                                            <td>The name of the day of the week</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </n-grid-item>
                            <n-grid-item>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Format</th>
                                            <th>Output</th>
                                            <th>Description</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td><code>H</code></td>
                                            <td>0-23</td>
                                            <td>The hour</td>
                                        </tr>
                                        <tr>
                                            <td><code>HH</code></td>
                                            <td>00-23</td>
                                            <td>The hour, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>h</code></td>
                                            <td>1-12</td>
                                            <td>The hour, 12-hour clock</td>
                                        </tr>
                                        <tr>
                                            <td><code>hh</code></td>
                                            <td>01-12</td>
                                            <td>The hour, 12-hour clock, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>m</code></td>
                                            <td>0-59</td>
                                            <td>The minute</td>
                                        </tr>
                                        <tr>
                                            <td><code>mm</code></td>
                                            <td>00-59</td>
                                            <td>The minute, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>s</code></td>
                                            <td>0-59</td>
                                            <td>The second</td>
                                        </tr>
                                        <tr>
                                            <td><code>ss</code></td>
                                            <td>00-59</td>
                                            <td>The second, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>SSS</code></td>
                                            <td>000-999</td>
                                            <td>The millisecond, 3-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>Z</code></td>
                                            <td>+05:00</td>
                                            <td>The offset from UTC, ±HH:mm</td>
                                        </tr>
                                        <tr>
                                            <td><code>ZZ</code></td>
                                            <td>+0500</td>
                                            <td>The offset from UTC, ±HHmm</td>
                                        </tr>
                                        <tr>
                                            <td><code>A</code></td>
                                            <td>AM PM</td>
                                            <td></td>
                                        </tr>
                                        <tr>
                                            <td><code>a</code></td>
                                            <td>am pm</td>
                                            <td></td>
                                        </tr>
                                    </tbody>
                                </table>
                            </n-grid-item>
                        </n-grid>
                    </n-popover>
                </h3>
                <n-button-group>
                    <n-tag size="large">Datetime format mask</n-tag>
                    <n-input placeholder="Type datetime format current mask" v-model:value="currentDatetimeMask" />
                    <n-tag size="large">Preview</n-tag>
                    <n-input placeholder="mask preview (current datetime)" v-model:value="currentDatetimeMaskPreview"
                        readonly />
                </n-button-group>
            </div>
        </n-flex>
    </n-form>

</template>

<style lang="css" scoped></style>