import React from "react";
import ScheduleForm from "@/components/scheduleForm";
import { cookies } from "next/headers";
import { validateCookieToken } from "../../../api/auth";
import { Unauthorized } from "../../../components/error/ErrorComp";

export default async function page({ params }: { params: { id: string } }) {
  const cookieStore = cookies();
  const jwt = cookieStore.get("Authorization")?.value;
  let user;
  if (jwt) {
    user = await validateCookieToken(jwt || "");
  }
  if (!user) return <Unauthorized />;

  return (
    <div className="p-4">
      <ScheduleForm userId={user.ID} jobId={params.id} existing />
    </div>
  );
}
