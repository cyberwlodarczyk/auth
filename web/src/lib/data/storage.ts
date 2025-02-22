const STORAGE = localStorage;

export enum StorageKey {
  Session = "session",
  Sudo = "sudo",
}

export const storage = {
  get: (key: StorageKey) => {
    return STORAGE.getItem(key);
  },
  set: (key: StorageKey, value: string) => {
    STORAGE.setItem(key, value);
  },
  remove: (key: StorageKey) => {
    STORAGE.removeItem(key);
  },
};
