import React from "react";
import ScheduleForm from "@/components/scheduleForm";
import { cookies } from "next/headers";
import { validateCookieToken } from "../../../api/auth";

export default async function page({ params }: { params: { id: string } }) {
  const cookieStore = cookies();
  const jwt = cookieStore.get("Authorization")?.value;
  let user;
  if (jwt) {
    user = await validateCookieToken(jwt || "");
  }

  return (
    <div className="p-4">
      <ScheduleForm userId={user?.id} jobId={params.id} edit />
    </div>
  );
}
