import React, { useState } from "react";
import {
  ColumnDef,
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  flexRender,
} from "@tanstack/react-table";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ArrowUpDown } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useRouter } from "next/navigation";

interface ClusterRole {
  metadata: {
    name: string;
    creationTimestamp: string;
  };
  rules: {
    apiGroups: string[];
    resources: string[];
    verbs: string[];
  }[];
}

const DataTable = ({ clusterRoles }: { clusterRoles: ClusterRole[] }) => {
  const [sorting, setSorting] = useState([]);
  const [columnFilters, setColumnFilters] = useState([]);
  const router = useRouter();
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

  const handleRowClick = (clusterRoleName: string) => {
    router.push(`/dashboard/cluster-roles/${clusterRoleName}`);
  };

  const table = useReactTable({
    data: clusterRoles,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    state: {
      sorting,
      columnFilters,
    },
  });

  return (
    <Card className="h-full">
      <CardHeader>
        <div className="justify-between item-start flex">
          <div className="flex flex-col gap-4">
            <CardTitle className="font-bold">Cluster Roles</CardTitle>
            <CardDescription>
              Manage Cluster Roles and Permissions
            </CardDescription>
          </div>
          <div className="flex items-center gap-2"></div>
        </div>
        <div className="flex items-center py-4">
          <Input
            placeholder="Filter names..."
            value={
              (table.getColumn("metadata.name")?.getFilterValue() as string) ??
              ""
            }
            onChange={(event) =>
              table
                .getColumn("metadata.name")
                ?.setFilterValue(event.target.value)
            }
            className="max-w-sm"
          />
        </div>
      </CardHeader>
      <CardContent>
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
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() => handleRowClick(row.original.metadata.name)}
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
    </Card>
  );
};

export default DataTable;
