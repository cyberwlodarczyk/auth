import { decodeJwt, errors, type JWTPayload } from "jose";

export interface Token {
  raw: string;
  payload: JWTPayload;
}

export const newToken = (raw: string): Token | null => {
  try {
    const payload = decodeJwt(raw);
    if (!payload.exp || payload.exp <= Date.now() / 1000) {
      return null;
    }
    return { raw, payload };
  } catch (error) {
    if (error instanceof errors.JWTInvalid) {
      return null;
    }
    throw error;
  }
};

export const newQueryToken = (): Token | null | undefined => {
  const params = new URLSearchParams(window.location.search);
  const raw = params.get("token");
  if (!raw) {
    return undefined;
  }
  return newToken(raw);
};

export const newTrustedToken = (raw: string): Token => {
  return { raw, payload: decodeJwt(raw) };
};
