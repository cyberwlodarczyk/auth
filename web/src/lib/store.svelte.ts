import { getUser, type User } from "./api";

const SESSION_KEY = "session";
const SUDO_KEY = "sudo";

export interface Store {
  location: string;
  session: string | null;
  sudo: string | null;
  user?: User | null;
}

export const store: Store = $state({
  location: decodeURIComponent(window.location.pathname),
  session: localStorage.getItem(SESSION_KEY),
  sudo: localStorage.getItem(SUDO_KEY),
});

$effect.root(() => {
  $effect(() => {
    if (store.session === localStorage.getItem(SESSION_KEY)) {
      return;
    }
    if (store.session) {
      localStorage.setItem(SESSION_KEY, store.session);
    } else {
      localStorage.removeItem(SESSION_KEY);
    }
  });

  $effect(() => {
    if (store.sudo === localStorage.getItem(SUDO_KEY)) {
      return;
    }
    if (store.sudo) {
      localStorage.setItem(SUDO_KEY, store.sudo);
    } else {
      localStorage.removeItem(SUDO_KEY);
    }
  });

  $effect(() => {
    if (store.location === decodeURIComponent(window.location.pathname)) {
      return;
    }
    history.pushState(null, "", store.location);
  });

  window.addEventListener("popstate", () => {
    store.location = decodeURIComponent(window.location.pathname);
  });

  $effect(() => {
    if (store.user === undefined) {
      if (store.session === null) {
        store.user = null;
      } else {
        getUser();
      }
    }
  });
});
