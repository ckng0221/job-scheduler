const BASE_URL = "http://localhost:8000";

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
  const url = `${BASE_URL}/scheduler/jobs`;
  const res = fetch(url, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: { "Content-Type": "application/json" },
  });

  return res;
}

export async function uploadTaskScript(jobId: string, file: File) {
  const url = `${BASE_URL}/scheduler/jobs/${jobId}/task-script`;
  const formdata = new FormData();
  formdata.append("file", file, file.name);

  const res = await fetch(url, {
    method: "POST",
    body: formdata,
  });
  return res;
}
