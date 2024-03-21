"use client";
import * as React from "react";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import {
  Button,
  FormControl,
  FormControlLabel,
  FormLabel,
  Paper,
  Radio,
  RadioGroup,
  Snackbar,
} from "@mui/material";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import dayjs from "dayjs";
import { IJob, submitJob } from "../api/job";

export default function ScheduleForm() {
  const initialJob = {
    JobName: "",
    IsRecurring: false,
    NextRunTime: 0,
    UserID: 1,
    Cron: "",
    IsDisabled: false,
  };
  const [job, setJob] = React.useState<IJob>(initialJob);
  const [scheduledDatetime, setScheduledDateTime] = React.useState(dayjs());
  const [openSnackbar, setOpenSnackBar] = React.useState(false);
  const [snackbarMessage, setSnackbarMessage] = React.useState("");

  async function submitForm(e: React.FormEvent) {
    e.preventDefault();
    let nextRunTimeUnix = 0;
    const payload: IJob = {
      ...job,
    };

    // update next run time for one-time job
    console.log(job);

    if (!job.IsRecurring) {
      nextRunTimeUnix = scheduledDatetime.unix();
      payload["NextRunTime"] = nextRunTimeUnix;
    } else {
      //TODO: create cron expression
      const cronExpression = "*/5 * * * *";
      payload["Cron"] = cronExpression;
    }

    const res = await submitJob(payload);
    if (res.ok) {
      setOpenSnackBar(true);
      setSnackbarMessage("Scheduled job created!");
      setJob(initialJob);
    } else {
      setOpenSnackBar(true);
      setSnackbarMessage("Failed to create schedule job");
    }
  }

  const handleClose = (
    event: React.SyntheticEvent | Event,
    reason?: string,
  ) => {
    if (reason === "clickaway") {
      return;
    }

    setOpenSnackBar(false);
  };

  return (
    <>
      <Paper elevation={3} className="p-16">
        <form className="" onSubmit={(e) => submitForm(e)}>
          <FormControl>
            <div>
              <TextField
                required
                name="job-name"
                id="outlined-required"
                label="Job Name"
                placeholder="Your job name"
                value={job.JobName}
                onChange={(e) => setJob({ ...job, JobName: e.target.value })}
              />
            </div>
            <div className="m-4">
              <FormLabel id="demo-row-radio-buttons-group-label">
                Frequency
              </FormLabel>
              <RadioGroup
                row
                aria-labelledby="demo-row-radio-buttons-group-label"
                name="row-radio-buttons-group"
              >
                <FormControlLabel
                  value="one-time"
                  control={<Radio />}
                  label="One-time"
                  checked={!job.IsRecurring}
                  onChange={() => setJob({ ...job, IsRecurring: false })}
                />
                <FormControlLabel
                  value="recurring"
                  control={<Radio />}
                  label="Recurring"
                  checked={job.IsRecurring}
                  onChange={() => setJob({ ...job, IsRecurring: true })}
                />
              </RadioGroup>
            </div>

            <FormLabel>Date & Time</FormLabel>
            <DateTimePicker
              className="mb-4"
              value={scheduledDatetime}
              onChange={(e) => setScheduledDateTime(e || dayjs(""))}
            />

            <Button variant="outlined" type="submit">
              Submit
            </Button>
          </FormControl>
        </form>
      </Paper>
      <Snackbar
        open={openSnackbar}
        autoHideDuration={3000}
        onClose={handleClose}
        message={snackbarMessage}
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        // action={action}
      />
    </>
  );
}
