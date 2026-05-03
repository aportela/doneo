import { type NotificationType, useNotification } from "naive-ui";

import { useUserSettingsStore } from "../stores/userSettings";

export const useNotify = () => {
  const userSettingsStore = useUserSettingsStore();
  const notification = useNotification();

  const notify = (type: NotificationType, content: string, meta: string) => {
    if (userSettingsStore.hasNotificationsEnabled) {
      notification[type]({
        content,
        meta,
        duration: 2500,
        keepAliveOnHover: true,
      });
    }
  };

  return { notify };
};
