"use client";

import { ArrowUpDown, MoreHorizontal } from "lucide-react";
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

import { object, z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { ColumnDef } from "@tanstack/react-table";
import { useRouter } from "next/navigation";
import { Role } from "../../_interfaces/role";
import { format } from "date-fns";
import { Checkbox } from "@/components/ui/checkbox";

export default function DataTable({ roles, namespace }) {
  const [data, setData] = useState(roles);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [rowSelection, setRowSelection] = useState({});

  const router = useRouter();
  const yaml = require("js-yaml");
  const yamlText = `
apiVersion: rbac.authorization.k8s.io/v1
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

  const formSchema = z.object({
    username: z.string().min(2).max(50),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
    },
  });

  const createRole = async (event) => {
    event.preventDefault();
    const jsonObject = yaml.load(formText);
    const jsonText = JSON.stringify(jsonObject, null, 2);
    const URL = "http://localhost:8080/api/roles";

    const response = await fetch(URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: jsonText,
    });

    const newRole = await response.json();
    setData((prevData) => [...prevData, newRole]);
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
  };

  const goToRolePage = async (role) => {
    await router.push(`/dashboard/roles/${role}`);
  };

  const columns: ColumnDef<Role>[] = [
    {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
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
        return format(new Date(timestamp), "hh:mm:ss a - MM/dd");
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
              <DropdownMenuItem onClick={() => deleteRole(namespace, name)}>
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
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

  const bulkDelete = () => {
    // TODO
  }

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
            <Button onClick={()=> bulkDelete()}>Delete</Button> 
            <Dialog>
              <DialogTrigger type="button">Create</DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create a Role</DialogTitle>
                  <DialogDescription>
                    Type in your information for creating a form.
                  </DialogDescription>
                </DialogHeader>

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
                        onSubmit={form.handleSubmit(createRole)}
                        className="space-y-8"
                      >
                        <FormField
                          control={form.control}
                          name="username"
                          render={({ field }) => (
                            <div>
                              <FormItem className="mb-4">
                                <FormLabel>Name</FormLabel>
                                <FormControl>
                                  <Input placeholder="shadcn" {...field} />
                                </FormControl>
                                <FormDescription>
                                  This is your public display name.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>

                              <FormItem>
                                <FormLabel>Namespace</FormLabel>
                                <FormControl>
                                  <Select>
                                    <SelectTrigger>
                                      <SelectValue placeholder="Namespace" />
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
                                <FormDescription>
                                  Select the namespace.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>

                              <FormItem>
                                <FormLabel>Verbs</FormLabel>
                                <FormControl>
                                  <div className="flex items-center justify-start gap-2 flex-wrap">
                                    <Badge variant="outline">Create</Badge>
                                    <Badge variant="outline">Delete</Badge>
                                    <Badge variant="outline">Get</Badge>
                                    <Badge variant="outline">List</Badge>
                                    <Badge variant="outline">Patch</Badge>
                                    <Badge variant="outline">Update</Badge>
                                    <Badge variant="outline">Watch</Badge>
                                  </div>
                                </FormControl>
                                <FormDescription>
                                  Select verb(s)
                                </FormDescription>
                                <FormMessage />
                              </FormItem>
                              <FormItem>
                                <FormLabel>Resources</FormLabel>
                                <FormControl>
                                  <Input placeholder="shadcn" {...field} />
                                </FormControl>
                                <FormDescription>
                                  Select the namespace.
                                </FormDescription>
                                <FormMessage />
                              </FormItem>
                            </div>
                          )}
                        />
                        <Button type="submit">Submit</Button>
                      </form>
                    </Form>
                  </TabsContent>
                  <TabsContent value="yaml">
                    <YamlEditor text={formText} onChange={handleChange} />
                    <Button onClick={() => createRole(event)}>
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
