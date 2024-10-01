"use client";

import React from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import {
  useClusterRoles,
  useDeleteClusterRoles,
} from "@/hooks/useClusterRoles";
import CreateClusterRoleDialog from "./_components/CreateClusterRoleDialog";
import { ClusterRole } from "@/interfaces/clusterRole";

const ClusterRoles = () => {
  const router = useRouter();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = React.useState(false);
  const { data: clusterRoles, isLoading, isError } = useClusterRoles();
  const deleteClusterRoleMutation = useDeleteClusterRoles();

  const columns: ColumnDef<ClusterRole>[] = [
    {
      id: "metadata.name",
      accessorKey: "metadata.name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
    },
    {
      accessorKey: "rules",
      header: "Rules",
      cell: ({ row }) => {
        const rules = row.original.rules;
        return <span>{rules.length} rule(s)</span>;
      },
    },
  ];

  const handleDeleteRole = (selectedRows: ClusterRole[]) => {
    selectedRows.forEach((row) => {
      deleteClusterRoleMutation.mutate(row);
    });
  };

  const routeToDetails = (name: string) => {
    router.push(`/dashboard/cluster-roles/${name}`);
  };

  if (isError) {
    return <div>Error</div>;
  }

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Cluster Roles</h1>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> Create Cluster Role
        </Button>
        <CreateClusterRoleDialog
          isOpen={isCreateDialogOpen}
          setIsOpen={setIsCreateDialogOpen}
        />
      </div>
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoles.items}
          columns={columns}
          enableGridView={false}
          rowActions={(row) => [
            <Trash
              key="delete"
              size={20}
              onClick={() => handleDeleteRole([row])}
            >
              Delete
            </Trash>,
            <Eye
              key="view"
              onClick={() => routeToDetails(row.metadata.name)}
              size={20}
            >
              View Details
            </Eye>,
          ]}
          enableRowSelection={true}
          bulkActions={(selectedRows) => [
            <Button
              key="delete"
              onClick={() => handleDeleteRole(selectedRows)}
              variant="destructive"
            >
              Delete Selected Cluster Roles
            </Button>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};
export default ClusterRoles;
