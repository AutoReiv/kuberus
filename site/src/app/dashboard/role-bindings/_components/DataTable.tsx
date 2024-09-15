import React, { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal, ArrowUpDown } from "lucide-react";
import {
  useReactTable,
  ColumnDef,
  getCoreRowModel,
  getSortedRowModel,
  SortingState,
  ColumnFiltersState,
  getFilteredRowModel,
  flexRender,
  getPaginationRowModel,
} from "@tanstack/react-table";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Copy,
  FileText,
  UserPlus,
  UserMinus,
  RefreshCw,
  Clock,
  FileCode,
} from "lucide-react";
import { toast } from "@/hooks/use-toast";
import yaml from "js-yaml";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";

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

const DataTable = ({ roleBindings }: { roleBindings: RoleBinding[] }) => {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [rowSelection, setRowSelection] = useState({});
  const [viewMode, setViewMode] = useState<"grid" | "table">("table");
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
      id: "name", // Add this line
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
    {
      id: "actions",
      cell: ({ row }) => {
        const roleBinding = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem>
                <FileText className="mr-2 h-4 w-4" />
                <Link
                  href={`/dashboard/role-bindings/${roleBinding.metadata.namespace}/${roleBinding.metadata.name}`}
                >
                  View
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Clone", roleBinding.metadata.name)}
              >
                <Copy className="mr-2 h-4 w-4" />
                Clone
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Add Subject", roleBinding.metadata.name)
                }
              >
                <UserPlus className="mr-2 h-4 w-4" />
                Add Subject
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Remove Subject", roleBinding.metadata.name)
                }
              >
                <UserMinus className="mr-2 h-4 w-4" />
                Remove Subject
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Change Role", roleBinding.metadata.name)
                }
              >
                <RefreshCw className="mr-2 h-4 w-4" />
                Change Role
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Export YAML", roleBinding.metadata.name)
                }
              >
                <FileCode className="mr-2 h-4 w-4" />
                Export YAML
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() =>
                  console.log("Temporary Access", roleBinding.metadata.name)
                }
              >
                <Clock className="mr-2 h-4 w-4" />
                Temporary Access
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Edit", roleBinding.metadata.name)}
              >
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => console.log("Delete", roleBinding.metadata.name)}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const table = useReactTable({
    data: roleBindings,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      rowSelection,
    },
  });

  const router = useRouter();

  const viewDetails = (roleBinding: string) => {
    router.push(`/dashboard/role-bindings/${roleBinding}`);
  };

  const cloneRoleBinding = (roleBinding: RoleBinding) => {
    // Implement cloning logic here
    console.log("Cloning role binding:", roleBinding.metadata.name);
    toast({
      title: "Role Binding Cloned",
      description: `A copy of ${roleBinding.metadata.name} has been created.`,
    });
  };

  const addSubject = (roleBinding: RoleBinding) => {
    // Implement add subject logic here
    console.log("Adding subject to role binding:", roleBinding.metadata.name);
    toast({
      title: "Add Subject",
      description: "Subject addition interface opened.",
    });
  };

  const removeSubject = (roleBinding: RoleBinding) => {
    // Implement remove subject logic here
    console.log(
      "Removing subject from role binding:",
      roleBinding.metadata.name
    );
    toast({
      title: "Remove Subject",
      description: "Subject removal interface opened.",
    });
  };

  const changeRole = (roleBinding: RoleBinding) => {
    // Implement change role logic here
    console.log("Changing role for role binding:", roleBinding.metadata.name);
    toast({
      title: "Change Role",
      description: "Role change interface opened.",
    });
  };

  const exportYAML = (roleBinding: RoleBinding) => {
    const yamlString = yaml.dump(roleBinding);
    const blob = new Blob([yamlString], { type: "text/yaml" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${roleBinding.metadata.name}.yaml`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    toast({
      title: "YAML Exported",
      description: `${roleBinding.metadata.name}.yaml has been downloaded.`,
    });
  };

  const setTemporaryAccess = (roleBinding: RoleBinding) => {
    // Implement temporary access logic here
    console.log(
      "Setting temporary access for role binding:",
      roleBinding.metadata.name
    );
    toast({
      title: "Temporary Access",
      description: "Temporary access configuration opened.",
    });
  };

  const editRoleBinding = (roleBinding: RoleBinding) => {
    router.push(
      `/dashboard/role-bindings/${roleBinding.metadata.namespace}/${roleBinding.metadata.name}/edit`
    );
  };

  const deleteRoleBinding = (roleBinding: RoleBinding) => {
    // Implement delete logic here
    console.log("Deleting role binding:", roleBinding.metadata.name);
    toast({
      title: "Role Binding Deleted",
      description: `${roleBinding.metadata.name} has been deleted.`,
    });
  };

  const GridView: React.FC<{ data: RoleBinding[] }> = ({ data }) => {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {data.map((role) => (
          <motion.div
            key={role.metadata.uid}
            whileHover={{ scale: 1.05 }}
            transition={{ type: "spring", stiffness: 100, mass: 0.5 }}
          >
            <Card className="h-full flex flex-col justify-between shadow-lg">
              <CardHeader className="bg-primary/10 rounded-t-lg">
                <CardTitle className="text-xl font-bold">
                  {role.metadata.name}
                </CardTitle>
                <CardDescription className="text-sm opacity-70">
                  {role.metadata.namespace}
                </CardDescription>
              </CardHeader>
              <CardContent className="flex-grow">
                <p className="text-sm mt-2">
                  Created:{" "}
                  {new Date(role.metadata.creationTimestamp).toLocaleString()}
                </p>
                {/* Add more role details here */}
              </CardContent>
              <CardFooter className="bg-secondary/10 rounded-b-lg">
                <Link
                  href={`/dashboard/roles/${role.metadata.namespace}/${role.metadata.name}`}
                  className="w-full"
                >
                  <Button className="w-full">View Details</Button>
                </Link>
              </CardFooter>
            </Card>
          </motion.div>
        ))}
      </div>
    );
  };

  return (
    <Card className="h-full">
      <CardHeader>
        <div className="justify-between item-start flex">
          <div className="flex flex-col gap-4">
            <CardTitle className="font-bold">Role Bindings</CardTitle>
            <CardDescription>Manage Role Bindings</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between py-4">
          <Input
            placeholder="Filter role bindings..."
            value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
            onChange={(event) =>
              table.getColumn("name")?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
          <Button
            onClick={() => setViewMode(viewMode === "table" ? "grid" : "table")}
          >
            {viewMode === "table" ? "Switch to Grid" : "Switch to Table"}
          </Button>
        </div>
        <AnimatePresence mode="wait">
          {viewMode === "table" ? (
            <motion.div
              key="table"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              <Table>
                <TableHeader>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map((header) => (
                        <TableHead key={header.id}>
                          {header.isPlaceholder
                            ? null
                            : flexRender(
                                header.column.columnDef.header,
                                header.getContext()
                              )}
                        </TableHead>
                      ))}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody>
                  {table.getRowModel().rows?.length ? (
                    table.getRowModel().rows.map((row) => (
                      <TableRow
                        key={row.id}
                        data-state={row.getIsSelected() && "selected"}
                      >
                        {row.getVisibleCells().map((cell) => (
                          <TableCell key={cell.id}>
                            {flexRender(
                              cell.column.columnDef.cell,
                              cell.getContext()
                            )}
                          </TableCell>
                        ))}
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell
                        colSpan={columns.length}
                        className="h-24 text-center"
                      >
                        No results.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </motion.div>
          ) : (
            <motion.div
              key="grid"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              <GridView data={roleBindings} />
            </motion.div>
          )}
        </AnimatePresence>

        <CardFooter>
          <div className="flex-1 text-sm text-muted-foreground">
            {table.getFilteredSelectedRowModel().rows.length} of{" "}
            {table.getFilteredRowModel().rows.length} row(s) selected.
          </div>
          <div className="flex items-center justify-end space-x-2 py-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              Next
            </Button>
          </div>
        </CardFooter>
      </CardContent>
    </Card>
  );
};

export default DataTable;
