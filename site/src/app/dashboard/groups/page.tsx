"use client";

import { apiClient } from "@/lib/apiClient";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";

type Group = string;

const Groups = () => {
  // Get Groups
  const {
    data: groups,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["groups"],
    queryFn: () => apiClient.getGroups(),
  });

  if (isError) {
    return <div>Error</div>;
  }

  // Define your columns
  const columns: ColumnDef<Group>[] = [
    {
      id: "name",
      header: "Name",
      cell: ({ row }) => <div>{row.original}</div>,
    }
  ];

  return (
    <div className="flex w-full flex-col">
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={groups}
          columns={columns}
          title="Groups"
          description="Groups are used to group users together. You can create groups to manage users and permissions." 
          // Add in row action to route to details
        ></GenericDataTable>
      )}
    </div>
  );
};

export default Groups;
