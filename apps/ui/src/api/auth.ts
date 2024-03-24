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

export async function exchangeProfile(access_token: string) {
  const url = `${BACKEND_HOST}/${MODULE}/google-token-exchange?code=${access_token}`;
  const res = fetch(url, {
    method: "POST",
  });

  return res;
}
export async function login(access_token: string) {
  const url = `${BACKEND_HOST}/auth/login`;
  const res = await fetch(url, {
    method: "POST",
    body: JSON.stringify({ code: access_token }),
    headers: { "Content-Type": "application/json" },
  });

  return res;
}

export async function validateCookieToken(access_token: string) {
  const endpoint = `${BACKEND_HOST}/${MODULE}/validate`;

  const headers = new Headers();
  headers.append("Cookie", `Authorization=${access_token}`);

  const res = await fetch(endpoint, { headers: headers });
  if (res.ok) {
    const user = await res.json();
    return user;
  }
}

export async function logout() {
  const endpoint = `${BACKEND_HOST}/${MODULE}/logout`;

  const headers = new Headers();
  headers.append("Cookie", `Authorization=${getCookie("Authorization")}`);
  const res = await fetch(endpoint, {
    method: "POST",
    headers,
  });

  const data = await res.json();
  return data;
}
