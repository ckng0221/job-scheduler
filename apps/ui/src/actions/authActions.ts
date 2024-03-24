"use server";
import { cookies } from "next/headers";

import { redirect } from "next/navigation";

import { getGoogleLoginUrl, login } from "../api/auth";

export async function loginRedirectAction() {
  const res = await getGoogleLoginUrl();
  if (res.ok) {
    const data = await res.json();
    const url = data.url;

    redirect(url);
  }
}

export async function loginAction(accesToken: string) {
  try {
    console.log("running login actions..");
    console.log(accesToken);
    const res = await login(accesToken);

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
