export const store = $state({
  location: window.location.pathname,
});

export const navigate = (path: string) => {
  history.pushState(null, "", path);
  store.location = path;
};

window.addEventListener("popstate", () => {
  store.location = window.location.pathname;
});
