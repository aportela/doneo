export type Validator = (value: string) => true | Error;

export const runValidators = (value: string, validators: Validator[]) => {
  for (const v of validators) {
    const result = v(value);
    if (result !== true) return result;
  }
  return true;
};

// TODO: i18n
export const required = (label: string): Validator => {
  return (value: string) => {
    if (!value) return new Error(`${label} is required`);
    return true;
  };
};

// TODO: i18n
export const minLength = (min: number): Validator => {
  return (value: string) => {
    if (!value) return true; // required se encarga de esto
    if (value.length < min) {
      return new Error(`min length is ${min}`);
    }
    return true;
  };
};

// TODO: i18n
export const validEmail: Validator = (value: string) => {
  if (!value) return true;
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(value)) {
    return new Error("invalid email format");
  }
  return true;
};
