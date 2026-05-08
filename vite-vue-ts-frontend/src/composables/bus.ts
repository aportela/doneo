import { useEventBus } from "@vueuse/core";

export type AppBusEvent =
  | {
      type: "reauthRequired";
      emitter: string;
    }
  | {
      type: "reauthValidNotify";
      to: string[];
    };

const bus = useEventBus<AppBusEvent>("doneo-app-bus");

export const useAppBus = () => {
  const emit = (event: AppBusEvent) => {
    bus.emit(event);
  };

  const on = (handler: (event: AppBusEvent) => void) => {
    return bus.on(handler);
  };

  const onType = <T extends AppBusEvent["type"]>(
    type: T,
    handler: (event: Extract<AppBusEvent, { type: T }>) => void,
  ) => {
    return bus.on((event) => {
      if (event.type === type) {
        handler(event as Extract<AppBusEvent, { type: T }>);
      }
    });
  };

  const emitReauthRequired = (emitter: string) => {
    bus.emit({
      type: "reauthRequired",
      emitter: emitter,
    });
  };

  const emitReauthValidNotify = (to: string[]) => {
    bus.emit({
      type: "reauthValidNotify",
      to,
    });
  };

  const reset = () => {
    bus.reset();
  };

  return {
    emit,
    on,
    onType,
    emitReauthRequired,
    emitReauthValidNotify,
    reset,
  };
};
