import React from "react";
import ScheduleForm from "@/components/scheduleForm";
import { cookies } from "next/headers";
import { validateCookieToken } from "../../../api/auth";
import { Unauthorized } from "../../../components/error/ErrorComp";
import { Breadcrumbs, Link } from "@mui/material";
import NextLink from "next/link";

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
      <Breadcrumbs aria-label="breadcrumb" className="mb-4">
        <Link
          underline="hover"
          color="inherit"
          href="/jobs"
          component={NextLink}
        >
          Jobs
        </Link>
        <Link
          underline="hover"
          color="text.primary"
          href={`/jobs/${params.id}`}
          aria-current="page"
          component={NextLink}
        >
          {params.id}
        </Link>
      </Breadcrumbs>
      <ScheduleForm userId={user.ID} jobId={params.id} existing />
    </div>
  );
}
