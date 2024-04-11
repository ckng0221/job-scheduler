import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { DateProvider } from "./providers";
import NavBar from "../components/NavBar";
import { Toaster } from "react-hot-toast";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Job Scheduler",
  description: "A job scheduler",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <DateProvider>
        <html lang="en">
          <body className={`${inter.className} p-2`}>
            <>
              <NavBar />
              {children}
              <Toaster />
            </>
          </body>
        </html>
      </DateProvider>
    </>
  );
}
