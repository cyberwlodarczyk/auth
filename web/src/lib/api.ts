import { store } from "./store.svelte";

function toBase64(text: string) {
  return btoa(String.fromCharCode(...new TextEncoder().encode(text)));
}

export interface User {
  id: number;
  email: string;
  name: string;
  createdAt: string;
}

async function request<T = null>(
  path: string,
  method: string,
  json: any = null,
  session: string | null = null
): Promise<T> {
  const headers: Record<string, string> = {};
  const init: RequestInit = { method, headers };
  if (json) {
    init.body = JSON.stringify(json);
    headers["Content-Type"] = "application/json";
  }
  if (session) {
    headers["Authorization"] = `Bearer ${session}`;
  }
  const res = await fetch(`/api${path}`, init);
  if (res.status === 401) {
    return null as T;
  }
  return res.headers.get("Content-Type") === "application/json"
    ? res.json()
    : null;
}

export async function signUpInit(email: string) {
  await request("/user/token/confirmation", "POST", { email });
}

export async function signUpFinish(
  token: string,
  name: string,
  password: string
) {
  const { session, user } = await request<{ session: string; user: User }>(
    "/user",
    "POST",
    {
      token,
      name,
      password: toBase64(password),
    }
  );
  store.session = session;
  store.user = user;
}

export async function signIn(email: string, password: string) {
  const { token, user } = await request<{ token: string; user: User }>(
    "/user/token/session",
    "POST",
    {
      email,
      password: toBase64(password),
    }
  );
  store.session = token;
  store.user = user;
}

export async function resetPasswordInit(email: string) {
  await request("/user/token/password-reset", "POST", { email });
}

export async function resetPasswordFinish(token: string, password: string) {
  const { session, user } = await request<{ session: string; user: User }>(
    "/user/password-reset",
    "POST",
    {
      token,
      password: toBase64(password),
    }
  );
  store.session = session;
  store.user = user;
}

export async function getUser() {
  const res = await request<{ user: User } | null>(
    "/user",
    "GET",
    null,
    store.session
  );
  if (!res) {
    store.user = null;
  } else {
    store.user = res.user;
  }
}

export async function changePassword(oldPassword: string, newPassword: string) {
  await request(
    "/user/password",
    "PUT",
    {
      password: toBase64(oldPassword),
      newPassword: toBase64(newPassword),
    },
    store.session
  );
}

export async function enterSudoMode(password: string) {
  await request(
    "/user/token/sudo",
    "POST",
    {
      password: toBase64(password),
    },
    store.session
  );
}

export async function changeEmail(token: string) {
  await request("/user/email", "PUT", { token }, store.sudo);
}
