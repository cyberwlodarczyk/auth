export const store = $state({
  location: decodeURIComponent(window.location.pathname),
});

export const navigate = (path: string) => {
  history.pushState(null, "", path);
  store.location = path;
};

window.addEventListener("popstate", () => {
  store.location = decodeURIComponent(window.location.pathname);
});
