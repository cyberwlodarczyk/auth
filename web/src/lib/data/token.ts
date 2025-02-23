import { decodeJwt, errors } from "jose";

export interface Token {
  raw: string;
  expiresAt: number;
}

export const newToken = (raw: string): Token | null => {
  try {
    const payload = decodeJwt(raw);
    const expiresAt = payload.exp;
    if (!expiresAt || expiresAt <= Date.now() / 1000) {
      return null;
    }
    return { raw, expiresAt };
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
  const payload = decodeJwt(raw);
  const expiresAt = payload.exp;
  if (!expiresAt) {
    throw new Error("trusted token has no expiration time");
  }
  return { raw, expiresAt };
};
