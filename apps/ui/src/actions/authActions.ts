"use server";
import { cookies } from "next/headers";

import { redirect } from "next/navigation";

import { getGoogleLoginUrl, login } from "../api/auth";

export async function loginRedirectAction() {
  const res = await getGoogleLoginUrl();
  if (res.ok) {
    const data = await res.json();
    const url = data.url;

    cookies().set({
      name: "state",
      value: data.state,
      httpOnly: false,
      path: "/",
    });
    cookies().set({
      name: "nonce",
      value: data.nonce,
      httpOnly: false,
      path: "/",
    });

    redirect(url);
  }
}

export async function loginAction(
  authorizationCode: string,
  state: string,
  cookieState: string,
  nounce: string,
) {
  try {
    console.log("running login actions..");
    console.log(authorizationCode);

    const res = await login(authorizationCode, state, cookieState, nounce);

    //   console.log(res);
    if (res.ok) {
      const data = await res.json();
      // set cookies
      cookies().set({
        name: "Authorization",
        value: data.access_token,
        httpOnly: false,
        path: "/",
      });
      // console.log(data);
      console.log("Login successful");

      return { message: "success", accesToken: data.access_token };
    } else {
      console.log("Failed");
      return { message: "Failed to login." };
    }
  } catch (err) {
    return { error: "Login Error" };
  }
}

export async function logoutAction() {
  try {
    cookies().delete("Authorization");

    // revalidatePath("/");
    // redirect("/");
  } catch (err) {
    return {
      error: "Failed to logout",
    };
  }
}
