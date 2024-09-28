"use client";

import React, { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/apiClient";
import GenericDataTable from "@/components/GenericDataTable";
import { Checkbox } from "@/components/ui/checkbox";
import { ColumnDef } from "@tanstack/react-table";
import { Role } from "../_interfaces/role";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, Eye, MoreHorizontal, Plus, Trash } from "lucide-react";
import { format } from "date-fns";
import { SkeletonPage } from "@/components/SkeletonPage";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { z } from "zod";
import YamlEditor from "@/components/YamlEditor";
import { toast } from "sonner";
import { cn } from "@/lib/utils";
import CreateRole from "./_components/CreateRole";
import yaml from "js-yaml";

/**
 * Renders a component that displays a list of roles and namespaces.
 *
 * The component uses the `useQuery` hook from `@tanstack/react-query` to fetch the list of roles and namespaces from the API.
 * If the data is still being fetched, a skeleton loader is displayed. Otherwise, a `DataTable` component is rendered with the fetched roles and namespaces.
 */
const dnsNameRegex =
  /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

const resourceSchema = z.object({
  name: z.string().min(1, "Resource name is required"),
  verbs: z.array(z.string()).min(1, "At least one verb is required"),
});

const roleSchema = z
  .object({
    roleName: z
      .string()
      .min(1, "Role name is required")
      .regex(dnsNameRegex, "Role name must be DNS compliant"),
    namespace: z
      .string()
      .min(1, "Namespace is required")
      .regex(dnsNameRegex, "Role name must be DNS compliant"),
    resources: z
      .array(resourceSchema)
      .min(1, "At least one resource is required"),
    apiGroup: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if (data.resources.some((resource) => resource.verbs.length === 0)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Each resource must have at least one verb selected",
        path: ["resources"],
      });
    }
  });

const Roles = () => {
  const router = useRouter();
  const [isCreateRoleDialogOpen, setIsCreateRoleDialogOpen] = useState(false);
  const queryClient = useQueryClient();
  const [yamlContent, setYamlContent] = useState("");
  const initialData = `apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]`;
  // Get Roles
  const {
    data: roles,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["roles"],
    queryFn: () => apiClient.getRoles(),
  });

  const columns: ColumnDef<Role>[] = [
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
      id: "name",
      header: ({ column }) => {
        return (
          <span
            className="flex items-center gap-2"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </span>
        );
      },
      accessorKey: "metadata.name",
      cell: ({ row }) => {
        const isActive = row.original.active;
        const name = row.getValue("name") as string;
        return (
          <TooltipProvider>
            <Tooltip delayDuration={300}>
              <TooltipTrigger asChild>
                <div className="flex items-center gap-2">
                  <div
                    className={cn(
                      "h-3 w-3 rounded-full",
                      isActive ? "bg-green-500 animate-pulse" : "bg-gray-500"
                    )}
                  />
                  {name}
                </div>
              </TooltipTrigger>
              <TooltipContent side="top" align="start" sideOffset={5}>
                {isActive ? "Active" : "Inactive"}
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        );
      },
    },
    {
      accessorKey: "metadata.namespace",
      header: ({ column }) => {
        return (
          <span
            className="flex items-center gap-2"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Namespace
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </span>
        );
      },
      cell: ({ row }) => row.original.metadata?.namespace || "N/A",
    },
    {
      id: "createdAt",
      header: ({ column }) => {
        return (
          <div
            className="flex items-center gap-2"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Created At
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </div>
        );
      },
      accessorKey: "metadata.creationTimestamp",
      cell: ({ getValue }) => {
        const timestamp: any = getValue();
        return format(new Date(timestamp), "MM/dd - hh:mm:ss a");
      },
    },
  ];

  const handleCreateRoleFromYaml = (content: string) => {
    try {
      const roleData = yaml.load(content);
      createRoleMutation.mutate(roleData);
    } catch (error) {
      toast.error(`Error parsing YAML: ${error.message}`);
    }
  };

  const routeToDetails = (namespace: string, name: string) => {
    router.push(`/dashboard/roles/${namespace}/${name}`);
    toast(`Routing to details page for ${name}`);
  };

  const createRoleMutation = useMutation({
    mutationFn: (roleData: any) => apiClient.createRole(roleData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
      setIsCreateRoleDialogOpen(false);
      toast.success(`Role has been created successfully`);
    },
    onError: (error: any) => {
      if (error.response && error.response.status === 409) {
        toast.error(
          `Error: A role with this name already exists in the specified namespace.`
        );
      } else {
        toast.error(`Error creating role: ${error.message}`);
      }
    },
  });
  const handleCreateRole = (data: z.infer<typeof roleSchema>) => {
    const existingRole = roles.find(
      (role) =>
        role.metadata.name === data.roleName &&
        role.metadata.namespace === data.namespace
    );

    if (existingRole) {
      toast.error(
        `A role with the name "${data.roleName}" already exists in the "${data.namespace}" namespace.`
      );
      return;
    }

    const roleData = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "Role",
      metadata: {
        name: data.roleName,
        namespace: data.namespace,
      },
      rules: data.resources.map((resource) => ({
        apiGroups: [data.apiGroup || ""],
        resources: [resource.name],
        verbs: resource.verbs,
      })),
    };
    createRoleMutation.mutate(roleData);
  };

  const handleDeleteRole = (row: Role[]) => {
    deleteRolesMutation.mutate(row);
  };

  const deleteRolesMutation = useMutation({
    mutationFn: (roles: Role | Role[]) => {
      if (Array.isArray(roles)) {
        return Promise.all(
          roles.map((role) =>
            apiClient.deleteRole(role.metadata.namespace, role.metadata.name)
          )
        );
      } else {
        return apiClient.deleteRole(
          roles.metadata.namespace,
          roles.metadata.name
        );
      }
    },
    onSuccess: (_, roles) => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
      const message = Array.isArray(roles)
        ? `Selected roles have been deleted successfully`
        : `Role ${roles.metadata.name} has been deleted successfully`;
      toast(message);
    },
    onSettled: (data, error, variables) => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
      if (error) {
        toast(`Error deleting role(s): ${error.message}`);
      } else {
        const message = Array.isArray(variables)
          ? `Selected roles have been deleted successfully`
          : `Role ${variables.metadata.name} has been deleted successfully`;
        toast(message);
      }
    },
  });

  if (isError) {
    return <div>Error</div>;
  }

  return (
    <motion.div
      className="flex-1 space-y-4"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
    >
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">Roles</h2>
        <Button onClick={() => setIsCreateRoleDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Role
        </Button>
        <ResponsiveDialog
          isOpen={isCreateRoleDialogOpen}
          setIsOpen={setIsCreateRoleDialogOpen}
          title="Create New Role"
          className="!max-w-none w-[60%] h-[60]"
        >
          <Tabs defaultValue="form">
            <TabsList className="w-full">
              <TabsTrigger value="form" className="w-full">
                FORM
              </TabsTrigger>
              <TabsTrigger value="yaml" className="w-full">
                YAML
              </TabsTrigger>
            </TabsList>
            <TabsContent value="form">
              <CreateRole onSubmit={handleCreateRole} />
            </TabsContent>
            <TabsContent value="yaml">
              <YamlEditor
                initialContent={initialData}
                onSave={handleCreateRoleFromYaml}
              />
            </TabsContent>
          </Tabs>
        </ResponsiveDialog>
      </div>
      {isLoading ? (
        <SkeletonPage></SkeletonPage>
      ) : (
        <GenericDataTable
          data={roles}
          columns={columns}
          enableGridView={false}
          enableRowSelection={true}
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
              onClick={() =>
                routeToDetails(row.metadata.namespace, row.metadata.name)
              }
              size={20}
            >
              View Details
            </Eye>,
          ]}
          bulkActions={(selectedRows) => [
            <Button
              key="delete"
              onClick={() => deleteRolesMutation.mutate(selectedRows)}
              variant="destructive"
            >
              Delete Selected Roles
            </Button>,
          ]}
        ></GenericDataTable>
      )}
    </motion.div>
  );
};
export default Roles;
