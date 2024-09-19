"use client";

import {
  ArrowUpDown,
  CheckCircle2,
  MoreHorizontal,
  XCircle,
} from "lucide-react";
import YamlEditor from "@focus-reactive/react-yaml";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  getPaginationRowModel,
  SortingState,
  getSortedRowModel,
  ColumnFiltersState,
  getFilteredRowModel,
} from "@tanstack/react-table";

import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { ColumnDef } from "@tanstack/react-table";
import { Role } from "../../_interfaces/role";
import { format } from "date-fns";
import { Checkbox } from "@/components/ui/checkbox";
import { toast } from "sonner";
import { ResponsiveDialog } from "./ResponsiveDialog";
import Link from "next/link";
import { AnimatePresence, motion } from "framer-motion";

const badges = [
  { name: "Create" },
  { name: "Delete" },
  { name: "Get" },
  { name: "List" },
  { name: "Patch" },
  { name: "Update" },
  { name: "Watch" },
];

const formSchema = z.object({
  nameOfRole: z
    .string()
    .min(2, "Please have at least 2 characters in your name.")
    .max(50, "50 characters is the max allowed")
    .refine(
      (value) => !/\s/.test(value),
      "Spaces are not allowed in the role name"
    ),
  namespaceForRole: z.string().min(1, "Please select an option"),
  badges: z.array(z.string()).refine((value) => value.some((item) => item), {
    message: "You have to select at least one item.",
  }),
});

export default function DataTable({ roles, namespace }) {
  const [data, setData] = useState(roles);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [rowSelection, setRowSelection] = useState({});
  const [createRoleDialogOpen, setCreateRoleDialogOpen] = useState(false);
  const [deleteConfirmationDialogs, setDeleteConfirmationDialogs] = useState<
    Record<string, boolean>
  >({});
  const [viewMode, setViewMode] = useState<"table" | "grid">("table");
  const setDeleteConfirmationDialog = (key: string, isOpen: boolean) => {
    setDeleteConfirmationDialogs((prev) => ({ ...prev, [key]: isOpen }));
  };
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nameOfRole: "",
      namespaceForRole: "",
      badges: ["Create"],
    },
  });

  const yaml = require("js-yaml");
  const yamlText = `apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role123
  namespace: default
rules:
- apiGroups:
    - ""
  resources:
    - pods
  verbs:
    - get
    - watch
    - list
`;
  const [formText, setformText] = useState(yamlText);

  useEffect(() => {
    setData(roles); // Update local state if props change
  }, [roles]);

  const handleChange = ({ text }) => {
    setformText(text);
  };

  const createRole = async (values, isForm?) => {
    const jsonObject = yaml.load(formText);
    const jsonText = JSON.stringify(jsonObject, null, 2);
    const URL = "http://localhost:8080/api/roles";
    const formPayload = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "Role",
      metadata: {
        name: values.nameOfRole,
        namespace: values.namespaceForRole,
      },
      rules: [
        {
          apiGroups: [""],
          resources: ["pods"],
          verbs: values.badges,
        },
      ],
    };

    const response = await fetch(URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        cookie: "session_token=eb1b485f-af8b-4631-9708-a87562ea6806",
      },
      credentials: "include",
      body: isForm ? JSON.stringify(formPayload) : jsonText,
    });

    const newRole = await response.json();
    setData((prevData) => [...prevData, newRole]);
    setCreateRoleDialogOpen(false);
    toast(
      <div className="flex items-center justify-start gap-4">
        <CheckCircle2 className="text-green-500" />
        <span>{`Successfully created ${newRole.metadata.name} in the ${newRole.metadata.namespace} namespace.`}</span>
      </div>
    );
  };

  const deleteRole = async (namespace, name) => {
    const dialogKey = `${namespace}-${name}`;
    const URL = `http://localhost:8080/api/roles?namespace=${namespace}&name=${name}`;
    const response = await fetch(URL, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (response.ok) {
      setData((prevData) =>
        prevData.filter(
          (role) =>
            !(
              role.metadata.name === name &&
              role.metadata.namespace === namespace
            )
        )
      );

      toast(
        <div className="flex items-center justify-start gap-4">
          <CheckCircle2 className="text-green-500" />
          <span>{`Successfully deleted ${name} in the ${namespace} namespace.`}</span>
        </div>
      );
    } else {
      toast(
        <div className="flex items-center justify-start gap-4">
          <XCircle className="text-red-500" />
          <span>{`Failed to delete ${name} in the ${namespace} namespace.`}</span>
        </div>
      );
    }

    setDeleteConfirmationDialog(dialogKey, false);
  };

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
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Name
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.name",
      cell: ({ row }) => {
        const name = row.getValue("name") as string;
        const namespace = row.original.metadata.namespace;
        return (
          <Link
            href={`/dashboard/roles/${namespace}/${name}`}
            className="hover:underline"
          >
            {name}
          </Link>
        );
      },
    },
    {
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Namespace
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.namespace",
    },
    {
      id: "createdAt",
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          >
            Created At
            <ArrowUpDown className="ml-2 h-4 w-4" />
          </Button>
        );
      },
      accessorKey: "metadata.creationTimestamp",
      cell: ({ getValue }) => {
        const timestamp: any = getValue();
        return format(new Date(timestamp), "MM/dd - hh:mm:ss a");
      },
    },
    {
      id: "actions",
      header: "",
      accessorKey: "metadata",
      cell: ({ row }) => {
        const { namespace, name } = row.original.metadata;
        const dialogKey = `${namespace}-${name}`;
        return (
          <>
            <ResponsiveDialog
              isOpen={deleteConfirmationDialogs[dialogKey] || false}
              setIsOpen={(isOpen) =>
                setDeleteConfirmationDialog(dialogKey, isOpen === true)
              }
              title="Delete Role"
              description="Are you sure you want to delete this role?"
            >
              <Button
                variant="destructive"
                onClick={() => {
                  deleteRole(namespace, name);
                  setDeleteConfirmationDialog(dialogKey, false);
                }}
              >
                Confirm
              </Button>
              <Button
                variant="ghost"
                onClick={() => setDeleteConfirmationDialog(dialogKey, false)}
              >
                Cancel
              </Button>
            </ResponsiveDialog>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button aria-haspopup="true" size="icon" variant="ghost">
                  <MoreHorizontal className="h-4 w-4" />
                  <span className="sr-only">Toggle menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem>Edit</DropdownMenuItem>
                <DropdownMenuItem>
                  <button
                    onClick={() => setDeleteConfirmationDialog(dialogKey, true)}
                  >
                    Delete
                  </button>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </>
        );
      },
    },
  ];

  const table = useReactTable({
    data,
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

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    createRole(values, true);
  };

  const GridView: React.FC<{ data: Role[] }> = ({ data }) => {
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
                  {format(
                    new Date(role.metadata.creationTimestamp),
                    "MM/dd - hh:mm:ss a"
                  )}
                </p>
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
            <CardTitle className="font-bold">Roles</CardTitle>
            <CardDescription>Manage User Roles and Permissions</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Dialog
              onOpenChange={() => {
                form.reset();
                setCreateRoleDialogOpen(true);
              }}
              open={createRoleDialogOpen}
            >
              <DialogTrigger
                type="button"
                className="bg-primary py-2 px-4 text-secondary rounded-md font-semibold"
              >
                Create
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create a Role</DialogTitle>
                  <DialogDescription>
                    Type in your information for creating a form.
                  </DialogDescription>
                </DialogHeader>
                <DialogClose
                  onClick={() => {
                    console.log("clicked");
                  }}
                ></DialogClose>
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
                    <Form {...form}>
                      <form
                        onSubmit={form.handleSubmit(onSubmit)}
                        className="space-y-8"
                      >
                        <FormField
                          control={form.control}
                          name="nameOfRole"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Name</FormLabel>
                              <FormControl>
                                <Input placeholder="Name" {...field} />
                              </FormControl>
                              <FormDescription>
                                This will be the name of your role.
                              </FormDescription>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name="namespaceForRole"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Namespace</FormLabel>
                              <FormControl>
                                <Select onValueChange={field.onChange}>
                                  <SelectTrigger>
                                    <SelectValue placeholder="Select namespace" />
                                  </SelectTrigger>
                                  <SelectContent>
                                    {namespace.map((namespace) => (
                                      <SelectItem
                                        value={namespace.metadata.name}
                                        key={namespace.metadata.uid}
                                      >
                                        {namespace.metadata.name}
                                      </SelectItem>
                                    ))}
                                  </SelectContent>
                                </Select>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name="badges"
                          render={({ field }) => (
                            <FormItem className="columns-2">
                              {badges.map((badge) => (
                                <FormField
                                  key={badge.name}
                                  control={form.control}
                                  name="badges"
                                  render={({ field }) => {
                                    return (
                                      <FormItem
                                        key={badge.name}
                                        className="flex flex-row items-start space-x-3 space-y-0"
                                      >
                                        <FormControl>
                                          <Checkbox
                                            checked={field.value?.includes(
                                              badge.name
                                            )}
                                            onCheckedChange={(checked) => {
                                              return checked
                                                ? field.onChange([
                                                    ...field.value,
                                                    badge.name,
                                                  ])
                                                : field.onChange(
                                                    field.value?.filter(
                                                      (value) =>
                                                        value !== badge.name
                                                    )
                                                  );
                                            }}
                                          />
                                        </FormControl>
                                        <FormLabel className="text-sm font-normal">
                                          {badge.name}
                                        </FormLabel>
                                      </FormItem>
                                    );
                                  }}
                                />
                              ))}
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <Button type="submit">Submit</Button>
                      </form>
                    </Form>
                  </TabsContent>
                  <TabsContent value="yaml">
                    <YamlEditor text={formText} onChange={handleChange} />
                    <Button variant="ghost" onClick={() => createRole(event)}>
                      Create Role{" "}
                    </Button>
                  </TabsContent>
                </Tabs>
              </DialogContent>
            </Dialog>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between py-4">
          <Input
            placeholder="Filter name..."
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
              <Table className="h-full">
                <TableHeader>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map((header) => {
                        return (
                          <TableHead key={header.id}>
                            {header.isPlaceholder
                              ? null
                              : flexRender(
                                  header.column.columnDef.header,
                                  header.getContext()
                                )}
                          </TableHead>
                        );
                      })}
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
              <GridView
                data={table
                  .getFilteredRowModel()
                  .rows.map((row) => row.original)}
              />
            </motion.div>
          )}
        </AnimatePresence>
      </CardContent>
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
    </Card>
  );
}
