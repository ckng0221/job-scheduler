import React from "react";
import JobList from "./JobList";
import { cookies } from "next/headers";
import { validateCookieToken } from "../../api/auth";
import { Unauthorized } from "../../components/error/ErrorComp";

async function page() {
  const cookieStore = cookies();
  const jwt = cookieStore.get("Authorization")?.value;
  let user;
  if (jwt) {
    user = await validateCookieToken(jwt || "");
  }
  if (!user) return <Unauthorized />;

  return (
    <div className="p-4">
      <JobList userId={user.ID} />
    </div>
  );
}

export default page;
