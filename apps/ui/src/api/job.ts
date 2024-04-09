import { getCookie } from "../utils/common";

const BACKEND_HOST = process.env.NEXT_PUBLIC_BACKEND_HOST;

export interface IJob {
  JobName: string;
  IsRecurring: boolean;
  FirstScheduledTime: Number;
  NextRunTime: Number;
  UserID: string;
  Cron: string;
  IsDisabled: boolean;
}

export interface IJobRead extends IJob {
  ID: string;
}

export interface IJobUpdate extends Partial<IJob> {}

export async function getUserJobs(userId: string) {
  const url = `${BACKEND_HOST}/scheduler/jobs?`;
  const headers = new Headers();
  headers.append("Content-Type", "application/json");
  headers.append("Authorization", `Bearer ${getCookie("Authorization")}` ?? "");

  const res = fetch(
    url +
      new URLSearchParams({
        user_id: userId,
      }),
    {
      method: "GET",
      headers: headers,
    },
  );

  return res;
}

export async function submitJob(payload: IJob) {
  const url = `${BACKEND_HOST}/scheduler/jobs`;
  const headers = new Headers();
  headers.append("Content-Type", "application/json");
  headers.append("Authorization", `Bearer ${getCookie("Authorization")}` ?? "");

  const res = fetch(url, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: headers,
  });

  return res;
}

export async function updateJob(jobId: string, payload: IJobUpdate) {
  const url = `${BACKEND_HOST}/scheduler/jobs/${jobId}`;
  const headers = new Headers();
  headers.append("Content-Type", "application/json");
  headers.append("Authorization", `Bearer ${getCookie("Authorization")}` ?? "");

  const res = fetch(url, {
    method: "PATCH",
    body: JSON.stringify(payload),
    headers: headers,
  });

  return res;
}

export async function uploadTaskScript(jobId: string, file: File) {
  const url = `${BACKEND_HOST}/scheduler/jobs/${jobId}/task-script`;
  const formdata = new FormData();
  formdata.append("file", file, file.name);
  const headers = new Headers();
  headers.append("Authorization", `Bearer ${getCookie("Authorization")}` ?? "");

  const res = await fetch(url, {
    method: "POST",
    body: formdata,
    headers: headers,
  });
  return res;
}
