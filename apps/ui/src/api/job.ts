const BASE_URL = "http://localhost:8000";

export interface IJob {
  JobName: string;
  IsRecurring: boolean;
  NextRunTime: Number;
  UserID: Number;
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
