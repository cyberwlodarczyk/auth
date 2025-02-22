import { store, newTrustedToken, type Token } from "$lib/data";

export interface User {
  id: number;
  email: string;
  name: string;
  createdAt: string;
}

const toBase64 = (text: string) => {
  return btoa(String.fromCharCode(...new TextEncoder().encode(text)));
};

interface RequestWithAuthOptions {
  bearer?: Token | null;
}

interface RequestWithoutBodyOptions {
  method: "GET" | "DELETE";
}

interface RequestWithBodyOptions {
  method: "POST" | "PUT";
  json: Record<string, string>;
}

type RequestOptions = RequestWithAuthOptions &
  (RequestWithBodyOptions | RequestWithoutBodyOptions);

const request = (path: string, options: RequestOptions) => {
  const { method, bearer, ...other } = options;
  const headers: Record<string, string> = {};
  const pathWithPrefix = `/api/${path}`;
  if (bearer) {
    headers["Authorization"] = `Bearer ${bearer.raw}`;
  }
  if ("json" in other) {
    headers["Content-Type"] = "application/json";
    const { json } = other;
    return fetch(pathWithPrefix, {
      method,
      headers,
      body: JSON.stringify(json),
    });
  } else {
    return fetch(pathWithPrefix, { method, headers });
  }
};

const getBody = async <T>(req: Promise<Response>) => {
  const res = await req;
  const body: T = await res.json();
  return body;
};

export const createUserConfirmationToken = async (params: {
  email: string;
}) => {
  await request("/user/token/confirmation", { method: "POST", json: params });
};

export const createUserSessionToken = async ({
  email,
  password,
}: {
  email: string;
  password: string;
}) => {
  const { token, user } = await getBody<{ token: string; user: User }>(
    request("/user/token/session", {
      method: "POST",
      json: {
        email,
        password: toBase64(password),
      },
    })
  );
  store.session = newTrustedToken(token);
  store.user = user;
};

export const createUserPasswordResetToken = async (params: {
  email: string;
}) => {
  await request("/user/token/password-reset", { method: "POST", json: params });
};

export const createUserSudoToken = async ({
  password,
}: {
  password: string;
}) => {
  await request("/user/token/sudo", {
    method: "POST",
    bearer: store.session,
    json: { password: toBase64(password) },
  });
};

export const getUser = async () => {
  const res = await request("/user", { method: "GET", bearer: store.session });
  if (res.status === 401) {
    store.user = null;
  } else {
    const { user }: { user: User } = await res.json();
    store.user = user;
  }
};

export const createUser = async ({
  password,
  ...other
}: {
  token: string;
  name: string;
  password: string;
}) => {
  const { session, user } = await getBody<{ session: string; user: User }>(
    request("/user", {
      method: "POST",
      json: { ...other, password: toBase64(password) },
    })
  );
  store.session = newTrustedToken(session);
  store.user = user;
};

export const editUserPassword = async ({
  password,
  newPassword,
}: {
  password: string;
  newPassword: string;
}) => {
  await request("/user/password", {
    method: "PUT",
    bearer: store.session,
    json: {
      password: toBase64(password),
      newPassword: toBase64(newPassword),
    },
  });
};

export const editUserEmail = async (params: { token: string }) => {
  await request("/user/email", {
    method: "PUT",
    bearer: store.sudo,
    json: params,
  });
};

export const resetUserPassword = async ({
  token,
  password,
}: {
  token: string;
  password: string;
}) => {
  const { session, user } = await getBody<{ session: string; user: User }>(
    request("/user/password-reset", {
      method: "POST",
      json: {
        token,
        password: toBase64(password),
      },
    })
  );
  store.session = newTrustedToken(session);
  store.user = user;
};
