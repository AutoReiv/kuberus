"use client";

import { apiClient } from "@/lib/apiClient";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { useRouter } from "next/navigation";

type User = {
  username: string;
  source: string;
};

const Users = () => {
  const router = useRouter()
  const {
    data: users,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["users"],
    queryFn: () => apiClient.getUsers(),
  });

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => {
        return (
          <div
            className="flex items-center gap-2"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Username
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </div>
        );
      },
    },
    {
      accessorKey: "source",
      header: ({ column }) => {
        return (
          <div
            className="flex items-center gap-2"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Source
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </div>
        );
      },
    },
  ];

  const routeToDetails = (name: string) => {
    router.push(`/dashboard/users/${name}`);
  };

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      {isLoading ? (
        <SkeletonPage />
      ) : isError ? (
        <div>Error loading users</div>
      ) : (
        <GenericDataTable
          data={users}
          columns={columns}
          title="Users"
          description="Manage user accounts and permissions"
          enableSorting
          enableFiltering
          enablePagination
          enableGridView
          onRowClick={(row) =>
            routeToDetails(row.username)
          }
        />
      )}
    </motion.div>
  );
};

export default Users;
