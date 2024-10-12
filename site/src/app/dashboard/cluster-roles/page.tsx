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
import { ActionsDropdown } from "@/components/ActionsDropdown";
import { DeletionConfirmationDialog } from "@/components/DeletionConfirmationDialog";

const ClusterRoles = () => {
  const router = useRouter();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = React.useState(false);
  const [deletionDialog, setDeletionDialog] = React.useState<{
    isOpen: boolean;
    role: ClusterRole | null;
  }>({
    isOpen: false,
    role: null,
  });

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
    {
      id: "actions",
      cell: ({ row }) => {
        const role = row.original;
        const actions = [
          {
            label: "View Details",
            icon: <Eye className="mr-2 h-4 w-4" />,
            onClick: () => routeToDetails(role.metadata.name),
          },
          {
            label: "Delete",
            icon: <Trash className="mr-2 h-4 w-4" />,
            onClick: () => handleDeleteRole(role),
          },
        ];
        return <ActionsDropdown actions={actions} />;
      },
    },
  ];

  const handleDeleteRole = (role: ClusterRole) => {
    setDeletionDialog({ isOpen: true, role });
  };

  const confirmDelete = () => {
    if (deletionDialog.role) {
      deleteClusterRoleMutation.mutate(deletionDialog.role);
    }
    setDeletionDialog({ isOpen: false, role: null });
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
          enableRowSelection={true}
          bulkActions={(selectedRows) => [
            <Button
              key="delete"
              onClick={() => handleDeleteRole(selectedRows[0])}
              variant="destructive"
            >
              Delete Selected Cluster Role
            </Button>,
          ]}
        ></GenericDataTable>
      )}
      <DeletionConfirmationDialog
        isOpen={deletionDialog.isOpen}
        onClose={() => setDeletionDialog({ isOpen: false, role: null })}
        onConfirm={confirmDelete}
        itemName={deletionDialog.role?.metadata.name || ""}
        itemType="Cluster Role"
      />
    </motion.div>
  );
};
export default ClusterRoles;
