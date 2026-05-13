import { useEventBus } from "@vueuse/core";

type AppBusEvents = {
  reauthRequired: { emitter: string };
  reauthValidNotify: { to: string[] };
};

type BusEventMap = {
  [K in keyof AppBusEvents]: {
    type: K;
    payload: AppBusEvents[K];
  };
};

type AppBusEvent = BusEventMap[keyof BusEventMap];

const bus = useEventBus<AppBusEvent>("doneo-app-bus");

const isEventType = <T extends keyof AppBusEvents>(
  event: AppBusEvent,
  type: T,
): event is BusEventMap[T] => {
  return event.type === type;
};

const emit = <T extends keyof AppBusEvents>(event: BusEventMap[T]) => {
  bus.emit(event);
};

const on = <T extends keyof AppBusEvents>(
  type: T,
  handler: (payload: AppBusEvents[T]) => void,
) => {
  return bus.on((event) => {
    if (!isEventType(event, type)) return;
    handler(event.payload);
  });
};

export const appBus = {
  emit,
  on,
};
