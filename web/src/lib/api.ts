async function request<T = null>(
  path: string,
  method: string,
  json: any = null
): Promise<T> {
  const init: RequestInit = { method };
  if (json) {
    init.body = JSON.stringify(json);
    init.headers = { "Content-Type": "application/json" };
  }
  const res = await fetch(`https://localhost:4000${path}`, init);
  return res.body ? res.json() : null;
}

export async function signUpInit(email: string) {
  await request("/user/token/confirmation", "POST", { email });
}

export async function signUpFinish(
  token: string,
  name: string,
  password: string
) {
  await request("/user", "POST", { token, name, password });
}

export async function signIn(email: string, password: string) {
  const res = await request<{ token: string }>("/user/token/session", "POST", {
    email,
    password,
  });
  localStorage.setItem("session", res.token);
}

export async function resetPasswordInit(email: string) {
  await request("/user/token/password-reset", "POST", { email });
}

export async function resetPasswordFinish(token: string, password: string) {
  const res = await request<{ token: string }>("/user/password-reset", "POST", {
    token,
    password,
  });
  localStorage.setItem("session", res.token);
}
