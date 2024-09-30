"use client";

import { useQuery, useQueryClient, useMutation } from "@tanstack/react-query";
import { ArrowUpDown, Package2, Plus, Trash } from "lucide-react";
import { Button } from "@/components/ui/button";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ColumnDef } from "@tanstack/react-table";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { motion } from "framer-motion";
import CreateUserDialog from "./_components/CreateUserDialog";

export interface User {
  username: string;
}

const Admin = () => {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedRow, setSelectedRow] = useState(null);
  const queryClient = useQueryClient();
  const [isCreateUserDialogOpen, setIsCreateUserDialogOpen] = useState(false);

  const {
    data: users,
    isFetched: usersFetched,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["users"],
    queryFn: () => apiClient.getAdminUsers(),
  });

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "username",
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
      cell: ({ row }) => <div>{row.original.username}</div>,
    },
  ];

  const handleCreateUser = (userData: {
    username: string;
    password: string;
    passwordConfirm: string;
  }) => {
    // Handle the creation of the user here
    // You can use the userData object to create the user
    // For example, you can send a request to your backend to create the user
    // and then update the users state with the new user
    // You can also use the queryClient to invalidate the users query to refetch the data
    // and update the UI
  };

  const deleteUserMutation = useMutation({
    mutationFn: (username: string) => apiClient.getAdminUserDelete(username),
    onMutate: (username) => {
      toast.loading(`Deleting user ${username}...`);
      // Optionally, you can implement optimistic updates here
    },
    onSuccess: (_, username) => {
      toast.dismiss();
      toast.success(`User ${username} has been deleted successfully.`);
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: (error) => {
      toast.dismiss();
      toast.error(
        error.message || "An unexpected error occurred while deleting the user."
      );
    },
  });

  const handleDeleteUser = (username: string) => {
    deleteUserMutation.mutate(username);
  };

  return (
    <div className="flex min-h-screen w-full flex-col">
      <main className="flex flex-1 flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-3xl font-bold tracking-tight">Admin Users</h2>
          <Button onClick={() => setIsCreateUserDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create User
          </Button>
          <CreateUserDialog
            isOpen={isCreateUserDialogOpen}
            setIsOpen={setIsCreateUserDialogOpen}
            onCreateUser={handleCreateUser}
          />
        </div>
        {isLoading && !usersFetched ? (
          <p>Loading Table...</p>
        ) : isError ? (
          <p>Error loading users</p>
        ) : (
          <GenericDataTable
            className="col-span-2"
            data={users}
            columns={columns}
            enableGridView={false}
            rowActions={(row) => [
              <Trash
                key="delete"
                size={16}
                onClick={() => handleDeleteUser(row.username)}
              >
                Delete
              </Trash>,
            ]}
          ></GenericDataTable>
        )}
      </main>
    </div>
  );
};

export default Admin;
