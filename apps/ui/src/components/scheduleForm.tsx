"use client";
import {
  Button,
  Checkbox,
  Chip,
  FormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Paper,
  Radio,
  RadioGroup,
  Select,
  SelectChangeEvent,
  Snackbar,
} from "@mui/material";
import Box from "@mui/material/Box";
import { Theme, styled, useTheme } from "@mui/material/styles";
import { DateTimeValidationError } from "@mui/x-date-pickers";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import dayjs from "dayjs";
import { useRouter } from "next/navigation";
import {
  Dispatch,
  FormEvent,
  ReactNode,
  SetStateAction,
  SyntheticEvent,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { loginAction } from "../actions/authActions";
import {
  IJob,
  IJobRead,
  getOneJob,
  readTaskScript,
  submitJob,
  uploadTaskScript,
} from "../api/job";
import { getCookie } from "../utils/common";
import InfoDialog from "./Dialog";

function generateCronExpression({
  scheduledDatetime,
  frequency,
  weekdays,
  isEveryMonth,
  months,
  dates,
}: {
  scheduledDatetime: dayjs.Dayjs;
  frequency: string;
  weekdays: string[];
  isEveryMonth: boolean;
  months: string[];
  dates: string[];
}) {
  let cronMin = "*",
    cronHour = "*",
    cronMonthDay = "*",
    cronMonth = "*",
    cronWeekDay = "*";

  cronMin = scheduledDatetime.minute().toString();
  cronHour = scheduledDatetime.hour().toString();

  switch (frequency) {
    case "daily":
      break;

    case "weekly":
      cronWeekDay = weekdays.join(",");
      break;

    case "monthly":
      cronMonthDay = dates.join(",");
      cronMonth = isEveryMonth ? "*" : months.map(getMonthId).join(",");

      break;

    default:
      break;
  }

  const cronExpression = `${cronMin} ${cronHour} ${cronMonthDay} ${cronMonth} ${cronWeekDay}`;

  return cronExpression;
}

function parseCronExpression(cron: string) {
  // let weekdays: number[] = [];
  // let monthdays: number[] = [];
  // let months: number[] = [];
  let triggerFrenquency = "daily";
  const cronArray = cron.split(" ");
  const cronMonthdays = cronArray[2];
  const cronMonths = cronArray[3];
  const cronWeekdays = cronArray[4];

  if (cronWeekdays !== "*") {
    triggerFrenquency = "weekly";
  } else if (cronMonthdays !== "*") {
    triggerFrenquency = "monthly";
  }

  // try {
  //   const interval = cronParser.parseExpression(cron);
  //   const fields = JSON.parse(JSON.stringify(interval.fields));
  //   weekdays = fields.dayOfWeek;
  //   monthdays = fields.dayOfMonth;
  //   months = fields.month;
  // } catch (err) {
  //   console.error(err);
  // }
  return { cronWeekdays, cronMonthdays, cronMonths, triggerFrenquency };
}

export default function ScheduleForm({
  userId,
  jobId,
  edit = false,
}: {
  userId: string;
  jobId?: string;
  edit?: boolean;
}) {
  const initialJob = {
    JobName: "",
    IsRecurring: false,
    FirstScheduledTime: 0,
    NextRunTime: 0,
    UserID: userId,
    Cron: "",
    IsDisabled: false,
  };
  const [job, setJob] = useState<IJob>(initialJob);
  const [scheduledDatetime, setScheduledDateTime] = useState(
    dayjs().add(1, "minute"),
  );
  const [openSnackbar, setOpenSnackBar] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [recurringFrequency, setRecurringFrequency] = useState("daily");
  const [error, setError] = useState<DateTimeValidationError | null>(null);
  // Recuring options
  const [selectedDays, setSelectedDays] = useState<string[]>([]);
  const [isEveryMonth, setIsEveryMonth] = useState(true);
  const [selectedMonths, setSelectedMonths] = useState<string[]>([]);
  const [selectedDates, setSelectedDates] = useState<string[]>([]);
  // file upload
  const [file, setFile] = useState<any>();
  const fileRef = useRef<any>(null);

  const router = useRouter();

  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    if (queryParams.has("code")) {
      const code = queryParams.get("code") || "";
      const state = queryParams.get("state") || "";

      window.history.replaceState({}, document.title, "/");
      console.log("login...");
      const cookieState = getCookie("state") || "";
      const nonce = getCookie("nonce") || "";
      // console.log("state", state);
      // console.log("nonce", nonce);

      loginAction(code, state, cookieState, nonce);
      // const token = getCookie("Authorization");
      // console.log(token);

      // To remove query parameters from url
      router.push("/");
    }
  }, [router]);

  // Fetch job
  useEffect(() => {
    async function fetchJob() {
      if (jobId) {
        const res = await getOneJob(jobId);
        if (res.ok) {
          const job = await res.json();
          setJob(job);
          if (job.IsRecurring) {
            const {
              cronWeekdays,
              cronMonthdays,
              cronMonths,
              triggerFrenquency,
            } = parseCronExpression(job.Cron);

            setRecurringFrequency(triggerFrenquency);
            if (triggerFrenquency == "weekly") {
              setSelectedDays(cronWeekdays.split(","));
            } else if (triggerFrenquency == "monthly") {
              if (cronMonths != "*") {
                setIsEveryMonth(false);
                setSelectedMonths(cronMonths.split(","));
              }
              setSelectedDates(cronMonthdays.split(","));
            }
          }
        }
      }
    }
    fetchJob();
  }, [jobId]);
  const errorMessage = useMemo(() => {
    switch (error) {
      case "disablePast": {
        return "Date time cannot earlier than current time.";
      }

      default: {
        return "";
      }
    }
  }, [error]);

  async function submitForm(e: FormEvent) {
    e.preventDefault();

    // form validation
    if (!userId) {
      setOpenSnackBar(true);
      setSnackbarMessage("Please login first");
      return;
    }

    if (scheduledDatetime.unix() < dayjs().unix()) {
      setError("disablePast");
      return;
    }

    let nextRunTimeUnix = scheduledDatetime.unix();
    const payload: IJob = {
      ...job,
      UserID: userId,
    };

    // update next run time for both one-time and recurring job
    payload["FirstScheduledTime"] = nextRunTimeUnix;
    payload["NextRunTime"] = nextRunTimeUnix;
    // console.log(job);

    if (job.IsRecurring) {
      const cronExpression = generateCronExpression({
        scheduledDatetime,
        frequency: recurringFrequency,
        isEveryMonth: isEveryMonth,
        months: selectedMonths,
        weekdays: selectedDays,
        dates: selectedDates,
      });
      console.log(recurringFrequency);
      console.log(cronExpression);

      payload["Cron"] = cronExpression;
    }

    const res = await submitJob(payload);
    if (res?.ok) {
      const data = await res.json();
      if (file && file.size > 0) {
        const res = await uploadTaskScript(data.ID, file);
        if (!res.ok) {
          alert("Failed to upload script.");
        }
      }

      setOpenSnackBar(true);
      setSnackbarMessage("Scheduled job created!");
      setJob(initialJob);
      fileRef.current.value = null;
    } else {
      setOpenSnackBar(true);
      setSnackbarMessage("Failed to create schedule job");
    }
  }

  const handleClose = (event: SyntheticEvent | Event, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }

    setOpenSnackBar(false);
  };

  const VisuallyHiddenInput = styled("input")({
    clip: "rect(0 0 0 0)",
    clipPath: "inset(50%)",
    height: 1,
    overflow: "hidden",
    position: "absolute",
    bottom: 0,
    left: 0,
    whiteSpace: "nowrap",
    width: 1,
  });

  return (
    <>
      <Paper elevation={3} className="p-16">
        <h1 className="font-medium text-lg mb-4">Job Scheduler</h1>
        <form className="" onSubmit={(e) => submitForm(e)}>
          <FormControl>
            <div>
              <div>
                <label
                  htmlFor="job-name-id"
                  className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Job Name
                </label>
                <input
                  type="text"
                  id="job-name-id"
                  className=" border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="Job Name"
                  required
                  value={job.JobName}
                  onChange={(e) => setJob({ ...job, JobName: e.target.value })}
                  onInvalid={(e) =>
                    (e.target as HTMLInputElement).setCustomValidity(
                      "Please enter your job name",
                    )
                  }
                  onInput={(e) =>
                    (e.target as HTMLInputElement).setCustomValidity("")
                  }
                />
              </div>
            </div>
            <div className="m-4">
              <FormLabel id="frequency-radio-btn">Frequency</FormLabel>
              <RadioGroup
                row
                aria-labelledby="frequency-radio-btn"
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
            {job.IsRecurring && (
              <>
                <div className="">
                  <FormLabel id="frequency-radio-btn">
                    Trigger Frequency
                  </FormLabel>
                  <RadioGroup
                    row
                    aria-labelledby="frequency-radio-btn"
                    name="row-radio-buttons-group"
                  >
                    <FormControlLabel
                      value="daily"
                      control={<Radio />}
                      label="Daily"
                      checked={recurringFrequency == "daily"}
                      onChange={() => setRecurringFrequency("daily")}
                    />
                    <FormControlLabel
                      value="weekly"
                      control={<Radio />}
                      label="Weekly"
                      checked={recurringFrequency == "weekly"}
                      onChange={() => setRecurringFrequency("weekly")}
                    />
                    <FormControlLabel
                      value="monthly"
                      control={<Radio />}
                      label="Monthly"
                      checked={recurringFrequency == "monthly"}
                      onChange={() => setRecurringFrequency("monthly")}
                    />
                  </RadioGroup>
                </div>
              </>
            )}

            <FormLabel id="datetimepicker">
              {job.IsRecurring ? "Start on" : "Scheduled on"}
            </FormLabel>
            <DateTimePicker
              className="mb-4"
              value={scheduledDatetime}
              onChange={(e) => setScheduledDateTime(e || dayjs(""))}
              disablePast
              format="DD/MM/YYYY hh:mm A"
              onError={(newError) => setError(newError)}
              slotProps={{
                textField: {
                  helperText: errorMessage,
                },
              }}
            />

            {recurringFrequency == "weekly" && job.IsRecurring && (
              <WeeklyOption
                selectedDays={selectedDays}
                setSelectedDays={setSelectedDays}
              />
            )}
            {recurringFrequency == "monthly" && job.IsRecurring && (
              <MonthlyOption
                isEveryMonth={isEveryMonth}
                setIsEveryMonth={setIsEveryMonth}
                selectedMonths={selectedMonths}
                setSelectedMonths={setSelectedMonths}
                selectedDates={selectedDates}
                setSelectedDates={setSelectedDates}
              />
            )}

            <TaskFileUpload setFile={setFile} fileRef={fileRef} job={job} />

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

const days = [
  { id: 1, day: "Monday" },
  { id: 2, day: "Tuesday" },
  { id: 3, day: "Wednesday" },
  { id: 4, day: "Thursday" },
  { id: 5, day: "Friday" },
  { id: 6, day: "Saturday" },
  { id: 0, day: "Sunday" },
];

const months = [
  { id: 1, name: "January" },
  { id: 2, name: "February" },
  { id: 3, name: "March" },
  { id: 4, name: "April" },
  { id: 5, name: "May" },
  { id: 6, name: "June" },
  { id: 7, name: "July" },
  { id: 8, name: "August" },
  { id: 9, name: "September" },
  { id: 10, name: "October" },
  { id: 11, name: "November" },
  { id: 12, name: "December" },
];

function WeeklyOption({
  selectedDays,
  setSelectedDays,
}: {
  selectedDays: string[];
  setSelectedDays: Dispatch<SetStateAction<string[]>>;
}) {
  function updateDays(dayId: string) {
    if (selectedDays.includes(dayId)) {
      setSelectedDays([...selectedDays.filter((id) => id != dayId)]);
    } else {
      setSelectedDays([...selectedDays, dayId]);
    }
  }

  return (
    <div className="mb-4">
      <FormGroup>
        {days.map((day, idx) => {
          return (
            <FormControlLabel
              key={idx}
              control={
                <Checkbox
                  value={day.id}
                  checked={selectedDays.includes(String(day?.id))}
                  onChange={(e) => updateDays(e.target.value)}
                />
              }
              label={day.day}
            />
          );
        })}
      </FormGroup>
    </div>
  );
}

interface IMonthlyProps {
  isEveryMonth: boolean;
  setIsEveryMonth: Dispatch<SetStateAction<boolean>>;
  selectedMonths: string[];
  setSelectedMonths: Dispatch<SetStateAction<string[]>>;
  selectedDates: string[];
  setSelectedDates: Dispatch<SetStateAction<string[]>>;
}

export function MonthlyOption(props: IMonthlyProps) {
  return (
    <>
      <MonthsOption
        isEveryMonth={props.isEveryMonth}
        setIsEveryMonth={props.setIsEveryMonth}
        selectedMonths={props.selectedMonths}
        setSelectedMonths={props.setSelectedMonths}
      />
      <DatesOption
        selectedDates={props.selectedDates}
        setSelectedDates={props.setSelectedDates}
      />
    </>
  );
}

export function MonthsOption({
  isEveryMonth,
  setIsEveryMonth,
  selectedMonths,
  setSelectedMonths,
}: {
  isEveryMonth: boolean;
  setIsEveryMonth: Dispatch<SetStateAction<boolean>>;
  selectedMonths: string[];
  setSelectedMonths: Dispatch<SetStateAction<string[]>>;
}) {
  const theme = useTheme();
  const ITEM_HEIGHT = 48;
  const ITEM_PADDING_TOP = 8;
  const MenuProps = {
    PaperProps: {
      style: {
        maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
        width: 250,
      },
    },
  };

  function getStyles(
    day: string,
    selectedDays: readonly string[],
    theme: Theme,
  ) {
    return {
      fontWeight:
        selectedDays.indexOf(day) === -1
          ? theme.typography.fontWeightRegular
          : theme.typography.fontWeightBold,
    };
  }

  const handleChange = (event: SelectChangeEvent<typeof selectedMonths>) => {
    const {
      target: { value },
    } = event;

    setSelectedMonths(
      // On autofill we get a stringified value.
      typeof value === "string" ? value.split(",") : value,
    );
  };

  return (
    <div>
      <FormControlLabel
        control={
          <Checkbox
            value={isEveryMonth}
            onChange={() => setIsEveryMonth(!isEveryMonth)}
            checked={isEveryMonth}
          />
        }
        label="Every month"
      />
      <br />
      <FormControl sx={{ m: 1, width: 300 }}>
        <InputLabel id="month-chip-label">Month</InputLabel>
        <Select
          labelId="month-chip-label"
          id="multiple-chip-month"
          multiple
          value={selectedMonths}
          onChange={handleChange}
          required={!isEveryMonth}
          input={<OutlinedInput id="select-multiple-chip-month" label="Chip" />}
          renderValue={(selected) => (
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected
                .sort((a, b) => (getMonthId(a) || 0) - (getMonthId(b) || 0))
                .map((value) => (
                  <Chip key={value} label={value.slice(0, 3)} />
                ))}
            </Box>
          )}
          MenuProps={MenuProps}
          disabled={isEveryMonth}
        >
          {months.map((month) => (
            <MenuItem
              key={month.id}
              value={month.name}
              style={getStyles(month.name, selectedMonths, theme)}
            >
              {month.name}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </div>
  );
}

export function DatesOption({
  selectedDates,
  setSelectedDates,
}: {
  selectedDates: string[];
  setSelectedDates: Dispatch<SetStateAction<string[]>>;
}) {
  const theme = useTheme();
  const ITEM_HEIGHT = 48;
  const ITEM_PADDING_TOP = 8;
  const MenuProps = {
    PaperProps: {
      style: {
        maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
        width: 250,
      },
    },
  };
  const calendarDates = Array.from({ length: 31 }, (_, i) => String(i + 1));

  function getStyles(
    day: string,
    selectedDays: readonly string[],
    theme: Theme,
  ) {
    return {
      fontWeight:
        selectedDays.indexOf(day) === -1
          ? theme.typography.fontWeightRegular
          : theme.typography.fontWeightBold,
    };
  }

  const handleChange = (event: SelectChangeEvent<typeof selectedDates>) => {
    const {
      target: { value },
    } = event;

    setSelectedDates(
      // On autofill we get a stringified value.
      typeof value === "string" ? value.split(",") : value,
    );
  };

  return (
    <div>
      <FormControl sx={{ m: 1, width: 300 }}>
        <InputLabel id="date-chip-label">Date</InputLabel>
        <Select
          labelId="date-chip-label"
          id="multiple-chip"
          required
          multiple
          value={selectedDates}
          onChange={handleChange}
          input={<OutlinedInput id="select-multiple-chip" label="Chip" />}
          renderValue={(selected) => (
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected
                .sort((a, b) => Number(a) - Number(b))
                .map((value) => (
                  <Chip key={value} label={value} />
                ))}
            </Box>
          )}
          MenuProps={MenuProps}
        >
          {calendarDates.map((date) => (
            <MenuItem
              key={date}
              value={date}
              style={getStyles(date, selectedDates, theme)}
            >
              {date}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </div>
  );
}

function getMonthId(monthName: string) {
  return months.find((month) => month.name == monthName)?.id;
}

function TaskFileUpload({
  job,
  fileRef,
  setFile,
}: {
  job: IJob;
  fileRef: any;
  setFile: Dispatch<any>;
}) {
  const acceptedExts = process.env.NEXT_PUBLIC_SUPPORTED_EXTENSIONS || "";
  let acceptedExsString;
  if (acceptedExts) {
    acceptedExsString = acceptedExts.split(",").sort().join(", ");
  }
  const [openDialog, setOpenDialog] = useState(false);
  const [scriptText, setScriptText] = useState<ReactNode>(<></>);

  const currentScriptFilenameArray = job?.TaskPath?.split("/");
  const currentScriptFilename =
    currentScriptFilenameArray?.[currentScriptFilenameArray?.length - 1];
  useEffect(() => {
    async function fetchScriptText() {
      const res = await readTaskScript(String(job?.ID));
      if (res.ok) {
        const script = await res.text();
        console.log(script);
        setScriptText(<pre className="text-sm">{script}</pre>);
      }
    }
    fetchScriptText();
  }, [job.ID]);

  return (
    <div className="my-4">
      <div className="block mb-4 text-sm font-medium text-gray-900 dark:text-white">
        <div>
          <div>Current Script</div>
          <div className="text-gray-500 dark:text-white">
            <a
              onClick={() => {
                setOpenDialog(true);
              }}
              href="#"
              className="text-blue-600 dark:text-blue-500 hover:underline"
            >
              {currentScriptFilename}
            </a>
          </div>
          <InfoDialog
            open={openDialog}
            setOpen={setOpenDialog}
            title={currentScriptFilename}
            body={scriptText}
          />
        </div>
      </div>
      <label
        className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
        htmlFor="file_input"
      >
        Upload Task Script
      </label>
      <input
        required
        className="block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 dark:text-gray-400 focus:outline-none dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400"
        id="file_input"
        name="task script"
        type="file"
        accept={acceptedExts}
        onChange={(e) => setFile(e?.target?.files?.[0])}
        ref={fileRef}
      />
      <p
        className="mt-1 text-sm text-gray-500 dark:text-gray-300"
        id="file_input_help"
      >
        {acceptedExsString}
      </p>
    </div>
  );
}
