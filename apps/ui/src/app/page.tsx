import { validateCookieToken } from "../api/auth";

import ScheduleForm from "../components/scheduleForm";
import { cookies } from "next/headers";

export default async function Home() {
  const cookieStore = cookies();
  const jwt = cookieStore.get("Authorization")?.value;
  let user;
  if (jwt) {
    user = await validateCookieToken(jwt || "");
  }
  // console.log(user);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="m-4">
        <ScheduleForm userId={user?.ID || ""} />
      </div>
      <div className="m-4 grid text-center lg:max-w-5xl lg:w-full lg:mb-0 lg:grid-cols-4 lg:text-left"></div>
    </main>
  );
}
