export interface FieldState {
  value: string;
  error: string | null;
}

export const newFieldState = (): FieldState => {
  return { value: "", error: null };
};

export const isFieldEmpty = (state: FieldState) => {
  if (state.value === "") {
    state.error = "This field is required";
    return true;
  }
  return false;
};

export const isEmailFieldInvalid = (state: FieldState) => {
  if (!/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$/.test(state.value)) {
    state.error = "Email is not in the correct format";
    return true;
  }
  return false;
};

export const isPasswordFieldInvalid = (state: FieldState) => {
  if (state.value.length < 12) {
    state.error = "Password must be at least 12 characters long";
    return true;
  }
  if (state.value.length > 64) {
    state.error = "Password must be at most 64 characters long";
    return true;
  }
  return false;
};

export const arePasswordFieldsDifferent = (
  state1: FieldState,
  state2: FieldState
) => {
  if (state1.value !== state2.value) {
    state2.error = "Passwords do not match";
    return true;
  }
  return false;
};
