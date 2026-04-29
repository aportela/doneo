import enUS from "./en-US";
import esES from "./es-ES";
import glGL from "./gl-GL";

const messages = {
  "en-US": enUS,
  "es-ES": esES,
  "gl-GL": glGL,
};

const availableSystemLocales: string[] = Object.keys(messages);

const localeSelectorOptionItems = [
  { shortLabel: "EN", label: "English", key: "en-US" },
  { shortLabel: "ES", label: "Español", key: "es-ES" },
  { shortLabel: "GL", label: "Galego", key: "gl-GL" },
];

const availableLocaleSelectorOptionItems = localeSelectorOptionItems.filter(
  (l) => availableSystemLocales.includes(l.key),
);

const getlocaleSelectorOptionItem = (locale: string) => {
  const currentIndex = availableLocaleSelectorOptionItems.findIndex(
    (l) => l.key === locale,
  );
  return availableLocaleSelectorOptionItems[
    currentIndex >= 0 ? currentIndex : 0
  ]!;
};

export {
  messages,
  availableSystemLocales,
  availableLocaleSelectorOptionItems,
  getlocaleSelectorOptionItem,
};
