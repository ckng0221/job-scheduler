"use client";
import React from "react";
import { loginRedirectAction } from "../actions/authActions";

export default function LoginBtn() {
  return (
    <button
      onClick={async () => {
        loginRedirectAction();
      }}
    >
      Login
    </button>
  );
}
