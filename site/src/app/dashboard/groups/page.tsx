"use client";

import { apiClient } from "@/lib/apiClient";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useGroups } from "@/hooks/useGroups";

type Group = string;

const Groups = () => {
  // Get Groups
  const { data: groups, isLoading, isError } = useGroups();

  if (isError) {
    return <div>Error</div>;
  }

  // Define your columns
  const columns: ColumnDef<Group>[] = [
    {
      id: "name",
      header: "Name",
      cell: ({ row }) => <div>{row.original}</div>,
    },
  ];

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={groups}
          columns={columns}
          enableGridView={false}
          title="Groups"
          description="Groups are used to group users together. You can create groups to manage users and permissions."
        ></GenericDataTable>
      )}
    </motion.div>
  );
};

export default Groups;
