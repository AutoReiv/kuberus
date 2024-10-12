"use client";

import React, { useState } from "react";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import {
  useClusterRoleBindings,
  useDeleteClusterRoleBinding,
} from "@/hooks/useClusterRoleBinding";
import { ClusterRoleBinding } from "@/interfaces/clusterRoleBinding";
import CreateClusterRoleBindingsDialog from "./_components/CreateClusterRoleBindingsDialog";
import { ActionsDropdown } from "@/components/ActionsDropdown";
import { DeletionConfirmationDialog } from "@/components/DeletionConfirmationDialog";

const ClusterRoleBindings = () => {
  const router = useRouter();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [deletionDialog, setDeletionDialog] = useState<{
    isOpen: boolean;
    itemName: string;
    itemType: string;
  }>({ isOpen: false, itemName: "", itemType: "" });

  const {
    data: clusterRoleBindings,
    isPending: isPendingRoles,
    isError,
  } = useClusterRoleBindings();

  const deleteClusterRoleBindingMutation = useDeleteClusterRoleBinding();

  const columns: ColumnDef<ClusterRoleBinding>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      cell: ({ row }) => <div>{row.original.metadata.name}</div>,
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
      cell: ({ row }) => <div>{row.original.roleRef.name}</div>,
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1">
          {row.original.subjects?.map((subject, index) => (
            <Badge key={index} variant="secondary">
              {subject.kind}: {subject.name}
            </Badge>
          )) ?? "No subjects"}
        </div>
      ),
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => (
        <div>
          {new Date(row.original.metadata.creationTimestamp).toLocaleString()}
        </div>
      ),
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
            onClick: () => handleDeleteClusterRoleBinding(role),
          },
        ];
        return <ActionsDropdown actions={actions} />;
      },
    },
  ];

  const handleDeleteClusterRoleBinding = (row: ClusterRoleBinding) => {
    setDeletionDialog({
      isOpen: true,
      itemName: row.metadata.name,
      itemType: "Cluster Role Binding",
    });
  };

  const confirmDelete = async () => {
    await deleteClusterRoleBindingMutation.mutate(deletionDialog.itemName);
    setDeletionDialog({ isOpen: false, itemName: "", itemType: "" });
  };
  const routeToDetails = (name: string) => {
    router.push(`/dashboard/cluster-role-bindings/${name}`);
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
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">
          Cluster Role Bindings
        </h2>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Cluster Role Binding
        </Button>
        <CreateClusterRoleBindingsDialog
          isOpen={isCreateDialogOpen}
          setIsOpen={setIsCreateDialogOpen}
        />
      </div>
      {isPendingRoles ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={clusterRoleBindings.items}
          columns={columns}
          enableGridView={false}
        />
      )}
      <DeletionConfirmationDialog
        isOpen={deletionDialog.isOpen}
        onClose={() => setDeletionDialog({ ...deletionDialog, isOpen: false })}
        onConfirm={confirmDelete}
        itemName={deletionDialog.itemName}
        itemType={deletionDialog.itemType}
      />
    </motion.div>
  );
};

export default ClusterRoleBindings;
