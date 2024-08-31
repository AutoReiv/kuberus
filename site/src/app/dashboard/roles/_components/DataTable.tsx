"use client";

import { ArrowUpDown, CheckCircle2, MoreHorizontal } from "lucide-react";
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
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { ResponsiveDialog } from "./ResponsiveDialog";
import Link from "next/link";

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
    .max(50, "50 characters is the max allowed"),
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
  const [deleteConfirmationDialog, setDeleteConfirmationDialog] = useState(false);

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
    const URL = `http://localhost:8080/api/roles?namespace=${namespace}&name=${name}`;
    await fetch(URL, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    setData((prevData) =>
      prevData.filter((role) => role.metadata.name !== name)
    );

    toast(
      <div className="flex items-center justify-start gap-4">
        <CheckCircle2 className="text-green-500" />
        <span>{`Successfully deleted ${name} in the ${namespace} namespace.`}</span>
      </div>
    );

    setDeleteConfirmationDialog(false);
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
        const id = row.original.metadata.name;
        return (
          <>
            <ResponsiveDialog
              isOpen={deleteConfirmationDialog}
              setIsOpen={setDeleteConfirmationDialog}
              title="Delete Role"
              description="Are you sure you want to delete this role?"
            >
              <Button
                variant="destructive"
                onClick={() => deleteRole(namespace, name)}
              >
                Confirm
              </Button>
              <Button
                variant="ghost"
                onClick={() => setDeleteConfirmationDialog(false)}
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
                {/* <DropdownMenuItem onClick={() => deleteRole(namespace, name)}>
                Delete
              </DropdownMenuItem> */}
                <DropdownMenuItem>
                  <button onClick={() => setDeleteConfirmationDialog(true)}>
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

  return (
    <Card className="h-full">
      <CardHeader>
        <div className="justify-between item-start flex">
          <div className="flex flex-col gap-4">
            <CardTitle className="font-bold">Roles</CardTitle>
            <CardDescription>Manage User Roles and Permissions</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            {/* TODO bulk delete */}
            {/* <Button onClick={() => bulkDelete()}>Delete</Button> */}
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
        <div className="flex items-center py-4">
          <Input
            placeholder="Filter name..."
            value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
            onChange={(event) =>
              table.getColumn("name")?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
        </div>
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
