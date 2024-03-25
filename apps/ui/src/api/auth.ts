import { getCookie } from "../utils/common";

const BACKEND_HOST = "http://localhost:8000";

const MODULE = "auth";

export async function getGoogleLoginUrl() {
  const url = `${BACKEND_HOST}/${MODULE}/google-login`;
  const res = fetch(url, {
    method: "GET",
  });

  return res;
}

export async function exchangeProfile(authorizationCode: string) {
  const url = `${BACKEND_HOST}/${MODULE}/google-token-exchange?code=${authorizationCode}`;
  const res = fetch(url, {
    method: "POST",
  });

  return res;
}
export async function login(
  authorizationCode: string,
  state: string,
  cookieState: string,
  nonce: string,
) {
  const url = `${BACKEND_HOST}/auth/login`;
  const headers = new Headers();
  headers.append("Cookie", `state=${cookieState}`);
  headers.append("Content-Type", "application/json");
  const res = await fetch(url, {
    method: "POST",
    body: JSON.stringify({
      code: authorizationCode,
      state: state,
      nonce: nonce,
    }),
    headers,
  });

  return res;
}

// Validate JWT token
export async function validateCookieToken(id_token: string) {
  const endpoint = `${BACKEND_HOST}/${MODULE}/validate`;

  const headers = new Headers();
  headers.append("Authorization", `Bearer ${id_token}`);

  const res = await fetch(endpoint, { headers: headers });
  if (res.ok) {
    const user = await res.json();
    return user;
  }
}

export async function logout() {
  const endpoint = `${BACKEND_HOST}/${MODULE}/logout`;

  const headers = new Headers();
  headers.append("Authorization", `Bearer ${getCookie("Authorization")}`);
  const res = await fetch(endpoint, {
    method: "POST",
    headers,
  });

  const data = await res.json();
  return data;
}
