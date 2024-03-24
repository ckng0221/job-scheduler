"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "react-hot-toast";
import { logoutAction } from "../actions/authActions";
import Modal from "./Modal";

export default function LogoutBtn({ className }: { className?: string }) {
  const [showModal, setShowModal] = useState(false);
  const router = useRouter();
  async function handleConfirm() {
    const result = await logoutAction();
    if (result?.error) {
      toast.error("Failed to log out.");
    }
    toast.success("Logged out");
    router.refresh();
  }

  return (
    <div className={className}>
      <button
        className="block py-2 px-3 text-gray-900 rounded hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-white md:dark:hover:text-blue-500 dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent"
        onClick={async () => {
          setShowModal(true);
        }}
      >
        Logout
      </button>
      <Modal
        title="Logout"
        description="Are you sure you want to log out?"
        openModal={showModal}
        setOpenModal={setShowModal}
        confirmAction={handleConfirm}
      />
    </div>
  );
}
