"use client";

import { Chip, Switch } from "@mui/material";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import dayjs from "dayjs";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { IJob, getUserJobs, updateJob } from "../../api/job";
import { renameKey } from "../../utils/common";

export default function JobList({ userId }: { userId: string }) {
  const [jobs, setjobs] = useState<IJob[]>([]);

  useEffect(() => {
    async function fetchJobs() {
      const res = await getUserJobs(userId);
      if (res.ok) {
        let jobs = await res.json();
        jobs.forEach((obj: IJob) => renameKey(obj, "ID", "id"));
        setjobs(jobs);
      } else {
        throw new Error("Failed to fetch jobs");
      }
    }
    fetchJobs();
  }, [userId]);

  return (
    <>
      <h1>Job List</h1>
      <DataTable jobs={jobs} setJobs={setjobs} />
    </>
  );
}

function DataTable({
  jobs,
  setJobs,
}: {
  jobs: IJob[];
  setJobs: Dispatch<SetStateAction<IJob[]>>;
}) {
  const columns: GridColDef[] = [
    { field: "id", headerName: "Job ID", width: 70 },
    { field: "JobName", headerName: "Job Name", width: 130 },
    {
      field: "CreatedAt",
      headerName: "Created Date",
      width: 180,
      renderCell: (params) => {
        const date = dayjs(params.value).format("DD/MM/YYYY hh:mm A");
        return <>{date}</>;
      },
    },
    {
      field: "IsRecurring",
      headerName: "Is Recurring",
      width: 130,
      renderCell: (params) => {
        const color: any = params.value == true ? "success" : "error";
        const labelText = params.value == true ? "Yes" : "No";
        return <Chip color={color} label={labelText} />;
      },
    },
    {
      field: "FirstScheduledTime",
      headerName: "Start Time",
      width: 180,
      renderCell: (params) => {
        let date = "";
        if (params.value > 0) {
          date = dayjs.unix(params.value).format("DD/MM/YYYY hh:mm A");
        }
        return <>{date}</>;
      },
    },
    {
      field: "NextRunTime",
      headerName: "Next Run Time",
      width: 180,
      renderCell: (params) => {
        let date = "";
        if (params.value > 0) {
          date = dayjs.unix(params.value).format("DD/MM/YYYY hh:mm A");
        }
        return <>{date}</>;
      },
    },
    {
      field: "IsDisabled",
      headerName: "Enabled",
      width: 130,
      renderCell: (params) => {
        const enabled = !jobs.find((job: any) => job?.id == params.id)
          ?.IsDisabled;

        return (
          <Switch
            checked={enabled}
            onClick={() => {
              updateEnabled(String(params.id), !enabled);
            }}
          />
        );
      },
    },
    {
      field: "IsCompleted",
      headerName: "Completed",
      width: 130,
      renderCell: (params) => {
        const color: any = params.value == true ? "success" : "error";
        const labelText = params.value == true ? "Yes" : "No";
        return <Chip color={color} label={labelText} />;
      },
    },
  ];

  function updateEnabled(jobId: string, enabled: boolean) {
    const idx = jobs.findIndex((job: any) => job.id == jobId);
    jobs[idx].IsDisabled = !enabled;
    setJobs(jobs);
    console.log("isDisabled", !enabled);
    updateJob(jobId, { IsDisabled: !enabled });
  }

  return (
    <div className="m-4" style={{ height: 800, width: "100%" }}>
      <DataGrid
        rows={jobs}
        columns={columns}
        initialState={{
          pagination: {
            paginationModel: { page: 0, pageSize: 10 },
          },
        }}
        pageSizeOptions={[10, 20]}
        // checkboxSelection
      />
    </div>
  );
}
