import { getUser, type User } from "./api";

export interface Store {
  location: string;
  session: string | null;
  user?: User | null;
}

export const store: Store = $state({
  location: decodeURIComponent(window.location.pathname),
  session: localStorage.getItem("session"),
});

$effect.root(() => {
  $effect(() => {
    if (store.session === localStorage.getItem("session")) {
      return;
    }
    if (store.session) {
      localStorage.setItem("session", store.session);
    } else {
      localStorage.removeItem("session");
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
