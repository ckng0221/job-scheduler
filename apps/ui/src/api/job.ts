import { getCookie } from "../utils/common";

const BACKEND_HOST = process.env.NEXT_PUBLIC_BACKEND_HOST;

export interface IJob {
  JobName: string;
  IsRecurring: boolean;
  NextRunTime: Number;
  UserID: string;
  Cron: string;
  IsDisabled: boolean;
}

export interface IJobRead extends IJob {
  ID: string;
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
