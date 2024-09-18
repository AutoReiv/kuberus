"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import React from "react";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { ArrowUpDown, Eye, Plus, Trash } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { ColumnDef } from "@tanstack/react-table";
import { SkeletonPage } from "@/components/SkeletonPage";
import { pageVariants } from "../layout";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

interface RoleBinding {
  metadata: {
    name: string;
    namespace: string;
    uid: string;
    resourceVersion: string;
    creationTimestamp: string;
    managedFields: {
      manager: string;
      operation: string;
      apiVersion: string;
      time: string;
      fieldsType: string;
      fieldsV1: {
        [key: string]: any;
      };
    }[];
  };
  subjects: {
    kind: string;
    apiGroup: string;
    name: string;
  }[];
  roleRef: {
    apiGroup: string;
    kind: string;
    name: string;
  };
}

const formSchema = z.object({
  roleBindingName: z.string().min(1, "Name is required"),
  namespace: z.string().min(1, "Namespace is required"),
  roleName: z.string().min(1, "Role name is required"),
  subjectKind: z.enum(["User", "Group", "ServiceAccount"]),
  subjectName: z.string().min(1, "Subject name is required"),
});

const RoleBindings = () => {
  const router = useRouter();
  const [isCreateRoleBindingDialogOpen, setIsCreateRoleBindingDialogOpen] =
    React.useState(false);
  const queryClient = useQueryClient();

  // Get Roles
  const { data: roleBindings, isLoading } = useQuery({
    queryKey: ["roleBindings"],
    queryFn: () => apiClient.getRoleBindings(),
  });

  const columns: ColumnDef<RoleBinding>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected()}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "metadata.name",
      id: "name",
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
    },
    {
      accessorKey: "metadata.namespace",
      header: "Namespace",
    },
    {
      accessorKey: "roleRef.name",
      header: "Role",
    },
    {
      accessorKey: "subjects",
      header: "Subjects",
      cell: ({ row }) => {
        const subjects = row.original.subjects;
        return (
          <div>
            {subjects.map((subject, index) => (
              <div key={index}>{`${subject.kind}: ${subject.name}`}</div>
            ))}
          </div>
        );
      },
    },
    {
      accessorKey: "metadata.creationTimestamp",
      header: "Created At",
      cell: ({ row }) => {
        return new Date(
          row.original.metadata.creationTimestamp
        ).toLocaleString();
      },
    },
  ];

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      roleBindingName: "",
      namespace: "",
      roleName: "",
      subjectKind: "User",
      subjectName: "",
    },
  });

  const routeToDetails = (roleBinding: RoleBinding) => {
    router.push(
      `/dashboard/role-bindings/${roleBinding.metadata.namespace}/${roleBinding.metadata.name}`
    );
    toast(
      `Routing to details page for ${roleBinding.metadata.name} in ${roleBinding.metadata.namespace}`
    );
  };

  const onDelete = async (roleBinding: RoleBinding) => {
    try {
      await apiClient.deleteRoleBinding(
        roleBinding.metadata.namespace,
        roleBinding.metadata.name
      );
      // Invalidate and refetch the roleBindings query
      queryClient.invalidateQueries({ queryKey: ["roleBindings"] });
      toast.success("Role binding deleted successfully");
    } catch (error) {
      toast.error("Failed to delete role binding");
    }
  };

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    const roleBindingData = {
      metadata: {
        name: values.roleBindingName,
        namespace: values.namespace,
      },
      roleRef: {
        kind: "Role",
        name: values.roleName,
        apiGroup: "rbac.authorization.k8s.io",
      },
      subjects: [
        {
          kind: values.subjectKind,
          name: values.subjectName,
          apiGroup: "rbac.authorization.k8s.io",
        },
      ],
    };
    await apiClient.createRoleBinding(roleBindingData);
    setIsCreateRoleBindingDialogOpen(false);
    form.reset();
    // Refetch role bindings
    queryClient.invalidateQueries({ queryKey: ["roleBindings"] });
  };

  return (
    <motion.div
      className="flex w-full flex-col"
      initial="initial"
      animate="animate"
      exit="exit"
      variants={pageVariants}
    >
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-3xl font-bold tracking-tight">Role Bindings</h2>
        <Button onClick={() => setIsCreateRoleBindingDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Role Binding
        </Button>
        <ResponsiveDialog
          isOpen={isCreateRoleBindingDialogOpen}
          setIsOpen={setIsCreateRoleBindingDialogOpen}
          title="Create New Role Binding"
        >
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="roleBindingName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter role binding name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="namespace"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Namespace</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter namespace" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="roleName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter role name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="subjectKind"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Subject Kind</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select subject kind" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="User">User</SelectItem>
                        <SelectItem value="Group">Group</SelectItem>
                        <SelectItem value="ServiceAccount">
                          Service Account
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="subjectName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Subject Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter subject name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">Create Role Binding</Button>
            </form>
          </Form>
        </ResponsiveDialog>
      </div>
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={roleBindings}
          columns={columns}
          rowActions={(row) => [
            <Trash key="delete" size={20} onClick={() => onDelete(row)}>
              Delete
            </Trash>,
            <Eye key="view" size={20} onClick={() => routeToDetails(row)}>
              View Details
            </Eye>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};

export default RoleBindings;
