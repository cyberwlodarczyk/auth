import {
  getUser,
  type User,
  newToken,
  type Token,
  storage,
  StorageKey,
} from "$lib/data";

export interface Store {
  location: string;
  session: Token | null;
  sudo: Token | null;
  user?: User | null;
}

const getLocation = () => {
  return decodeURIComponent(window.location.pathname);
};

const getToken = (key: StorageKey): Token | null => {
  const raw = storage.get(key);
  if (!raw) {
    return null;
  }
  const token = newToken(raw);
  if (!token) {
    storage.remove(key);
    return null;
  }
  return token;
};

const newStore = (): Store => {
  const location = getLocation();
  const session = getToken(StorageKey.Session);
  const sudo = getToken(StorageKey.Sudo);
  if (!session) {
    if (sudo) {
      storage.remove(StorageKey.Sudo);
    }
    return { location, session: null, sudo: null, user: null };
  }
  return { location, session, sudo };
};

export const store = $state(newStore());

if (store.user === undefined) {
  getUser();
}

window.addEventListener("popstate", () => {
  store.location = getLocation();
});

const createTokenEffect = (key: StorageKey) => () => {
  const raw = storage.get(key);
  const token = store[key];
  if (token) {
    if (token.raw !== raw) {
      storage.set(key, token.raw);
    }
  } else {
    if (raw) {
      storage.remove(key);
    }
  }
};

$effect.root(() => {
  $effect(createTokenEffect(StorageKey.Session));

  $effect(createTokenEffect(StorageKey.Sudo));

  $effect(() => {
    if (store.location !== getLocation()) {
      history.pushState(null, "", store.location);
    }
  });
});
