import React, { useState, useCallback, useEffect } from "react";
import {
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  flexRender,
  ColumnDef,
  SortingState,
  ColumnFiltersState,
  VisibilityState,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
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
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import { motion, AnimatePresence } from "framer-motion";
import useInfiniteScroll from "react-infinite-scroll-hook";
import { CSVLink } from "react-csv";
import {
  ChevronDown,
  Dot,
  DotSquare,
  Grid,
  MoreHorizontal,
  TableIcon,
} from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip";

interface GenericDataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  title?: string;
  description?: string;
  enableSorting?: boolean;
  enableFiltering?: boolean;
  enablePagination?: boolean;
  enableGridView?: boolean;
  enableColumnVisibility?: boolean;
  enableRowSelection?: boolean;
  isLoading?: boolean;
  error?: string;
  gridViewRenderer?: (item: T) => React.ReactNode;
  rowActions?: (row: T) => React.ReactNode[];
  bulkActions?: (selectedRows: T[]) => React.ReactNode[];
  subComponent?: (row: T) => React.ReactNode;
  onRowClick?: (item: T) => void;
  infiniteScroll?: boolean;
  loadMore?: () => void;
  hasMore?: boolean;
  className?: string;
  enableQuickActions?: boolean;
  quickActions?;
}

function GenericDataTable<T>({
  data,
  columns,
  title,
  description,
  enableSorting = true,
  enableFiltering = true,
  enablePagination = true,
  enableGridView = true,
  enableColumnVisibility = false,
  enableRowSelection = false,
  isLoading = false,
  error,
  gridViewRenderer,
  rowActions,
  bulkActions,
  subComponent,
  onRowClick,
  infiniteScroll = false,
  loadMore,
  hasMore = false,
  className,
  enableQuickActions = false,
  quickActions,
}: GenericDataTableProps<T>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});
  const [viewMode, setViewMode] = useState<"table" | "grid">("table");
  const [globalFilter, setGlobalFilter] = useState("");

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: enablePagination
      ? getPaginationRowModel()
      : undefined,
    getSortedRowModel: enableSorting ? getSortedRowModel() : undefined,
    getFilteredRowModel: enableFiltering ? getFilteredRowModel() : undefined,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    globalFilterFn: useCallback((row, columnId, filterValue) => {
      const value = row.getValue(columnId);
      return String(value).toLowerCase().includes(filterValue.toLowerCase());
    }, []),
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      globalFilter,
    },
  });

  const [sentryRef] = useInfiniteScroll({
    loading: isLoading,
    hasNextPage: hasMore,
    onLoadMore: loadMore,
    disabled: !infiniteScroll,
    rootMargin: "0px 0px 400px 0px",
  });

  const saveConfig = () => {
    const config = {
      sorting,
      columnFilters,
      columnVisibility,
    };
    localStorage.setItem("tableConfig", JSON.stringify(config));
  };

  const loadConfig = () => {
    const savedConfig = localStorage.getItem("tableConfig");
    if (savedConfig) {
      const config = JSON.parse(savedConfig);
      setSorting(config.sorting);
      setColumnFilters(config.columnFilters);
      setColumnVisibility(config.columnVisibility);
    }
  };

  const exportCSV = () => {
    const csvData = table.getRowModel().rows.map((row) => {
      const rowData: Record<string, any> = {};
      row.getVisibleCells().forEach((cell) => {
        rowData[cell.column.id] = cell.getValue();
      });
      return rowData;
    });
    return csvData;
  };

  const exportJSON = () => {
    const jsonData = table.getRowModel().rows.map((row) => {
      const rowData: Record<string, any> = {};
      row.getVisibleCells().forEach((cell) => {
        rowData[cell.column.id] = cell.getValue();
      });
      return rowData;
    });
    return JSON.stringify(jsonData, null, 2);
  };

  const renderQuickActions = (row: T) => {
    if (!enableQuickActions || !quickActions) return null;
    return (
      <div className="flex space-x-2">
        {quickActions(row).map((action, index) => (
          <React.Fragment key={index}>{action}</React.Fragment>
        ))}
      </div>
    );
  };

  const GridView = () => (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {data.map((item, index) => (
        <motion.div
          key={index}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: index * 0.1 }}
          whileHover={{ scale: 1.05 }}
          className="cursor-pointer"
          onClick={() => onRowClick && onRowClick(item)}
        >
          {gridViewRenderer ? (
            gridViewRenderer(item)
          ) : (
            <Card className="shadow-sm hover:shadow-md transition-shadow duration-300">
              <CardHeader>
                <CardTitle>
                  {(item as any).name || `Item ${index + 1}`}
                </CardTitle>
                <CardDescription>
                  {(item as any).description || "No description available"}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground">
                  {Object.entries(item)
                    .slice(0, 3)
                    .map(([key, value]) => (
                      <p key={key}>
                        <strong>{key}:</strong> {String(value)}
                      </p>
                    ))}
                </div>
              </CardContent>
            </Card>
          )}
        </motion.div>
      ))}
    </div>
  );

  const renderRowActions = (row: T) => {
    if (!rowActions) return null;
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="h-8 w-8 p-0">
            <span className="sr-only">Open menu</span>
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-auto bg-white dark:bg-gray-800 rounded-md shadow-lg flex items-center justify-center">
          {rowActions(row).map((action: any, index) => (
            <TooltipProvider key={index}>
              <Tooltip>
                <TooltipTrigger asChild>
                  <DropdownMenuItem className="cursor-pointer w-full px-4 py-2 rounded-sm flex items-center gap-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700">
                    {action}
                  </DropdownMenuItem>
                </TooltipTrigger>
                <TooltipContent>
                  <p>{`${action.props.children}`}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  };

  const renderBulkActions = () => {
    if (!bulkActions || Object.keys(rowSelection).length === 0) return null;
    const selectedRows = table
      .getSelectedRowModel()
      .rows.map((row) => row.original);
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button>Bulk Actions</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          {bulkActions(selectedRows).map((action, index) => (
            <DropdownMenuItem key={index}>{action}</DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <Card className={className}>
      <CardHeader className={title ? "inherit" : "py-1"}>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between py-4">
          <Input
            placeholder="Global search..."
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            className="max-w-sm"
          />
          <div className="flex items-center space-x-2">
            {enableColumnVisibility && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline">Columns</Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  {table
                    .getAllColumns()
                    .filter((column) => column.getCanHide())
                    .map((column) => {
                      return (
                        <DropdownMenuCheckboxItem
                          key={column.id}
                          className="capitalize"
                          checked={column.getIsVisible()}
                          onCheckedChange={(value) =>
                            column.toggleVisibility(!!value)
                          }
                        >
                          {column.id}
                        </DropdownMenuCheckboxItem>
                      );
                    })}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline">Actions</Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem onClick={saveConfig}>
                  Save Config
                </DropdownMenuItem>
                <DropdownMenuItem onClick={loadConfig}>
                  Load Config
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <CSVLink data={exportCSV()} filename="table_data.csv">
                    Export CSV
                  </CSVLink>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => {
                    const jsonString = `data:text/json;chatset=utf-8,${encodeURIComponent(
                      exportJSON()
                    )}`;
                    const link = document.createElement("a");
                    link.href = jsonString;
                    link.download = "table_data.json";
                    link.click();
                  }}
                >
                  Export JSON
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            {enableGridView && (
              <motion.div
                initial="initial"
                animate="animate"
                exit="exit"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Button
                  onClick={() =>
                    setViewMode(viewMode === "table" ? "grid" : "table")
                  }
                >
                  {viewMode === "table" ? <Grid /> : <TableIcon />}
                </Button>
              </motion.div>
            )}
            {renderBulkActions()}
          </div>
        </div>
        <AnimatePresence mode="wait">
          {viewMode === "table" ? (
            <motion.div
              key="table"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.4 }}
            >
              <Table>
                <TableHeader>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map((header) => (
                        <TableHead key={header.id}>
                          <motion.div
                            initial={{ opacity: 0, y: -20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.3 }}
                          >
                            {header.isPlaceholder
                              ? null
                              : flexRender(
                                  header.column.columnDef.header,
                                  header.getContext()
                                )}
                          </motion.div>
                        </TableHead>
                      ))}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody>
                  {table.getRowModel().rows?.length ? (
                    table.getRowModel().rows.map((row) => (
                      <React.Fragment key={row.id}>
                        <TableRow
                          key={row.id}
                          data-state={row.getIsSelected() && "selected"}
                          onClick={(e) => {
                            onRowClick && onRowClick(row.original);
                          }}
                          className="hover:bg-muted/50"
                        >
                          {row.getVisibleCells().map((cell) => (
                            <TableCell key={cell.id}>
                              {flexRender(
                                cell.column.columnDef.cell,
                                cell.getContext()
                              )}
                            </TableCell>
                          ))}
                          {enableQuickActions && (
                            <TableCell>
                              {renderQuickActions(row.original)}
                            </TableCell>
                          )}
                          {rowActions && (
                            <TableCell>
                              {renderRowActions(row.original)}
                            </TableCell>
                          )}
                        </TableRow>
                        {subComponent && (
                          <TableRow>
                            <TableCell
                              colSpan={columns.length + (rowActions ? 1 : 0)}
                            >
                              {subComponent(row.original)}
                            </TableCell>
                          </TableRow>
                        )}
                      </React.Fragment>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell
                        colSpan={columns.length + (rowActions ? 1 : 0)}
                        className="h-24 text-center"
                      >
                        No results.
                      </TableCell>
                    </TableRow>
                  )}
                  {infiniteScroll && hasMore && (
                    <TableRow ref={sentryRef}>
                      <TableCell
                        colSpan={columns.length + (rowActions ? 1 : 0)}
                        className="h-24 text-center"
                      >
                        Loading more...
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
              transition={{ duration: 0.4 }}
            >
              <GridView />
            </motion.div>
          )}
        </AnimatePresence>
      </CardContent>
      {enablePagination && !infiniteScroll && table.getPageCount() > 1 && (
        <CardFooter className="flex items-center justify-between py-4">
          <div className="flex items-center space-x-2">
            <span className="text-sm text-muted-foreground">
              Showing{" "}
              {table.getState().pagination.pageIndex *
                table.getState().pagination.pageSize +
                1}{" "}
              to{" "}
              {Math.min(
                (table.getState().pagination.pageIndex + 1) *
                  table.getState().pagination.pageSize,
                table.getFilteredRowModel().rows.length
              )}{" "}
              of {table.getFilteredRowModel().rows.length} entries
            </span>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  Show {table.getState().pagination.pageSize}
                  <ChevronDown className="ml-2 h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                {[10, 20, 30, 40, 50].map((pageSize) => (
                  <DropdownMenuItem
                    key={pageSize}
                    onSelect={() => table.setPageSize(pageSize)}
                  >
                    Show {pageSize}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.setPageIndex(0)}
              disabled={!table.getCanPreviousPage()}
            >
              First
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              Previous
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {table.getState().pagination.pageIndex + 1} of{" "}
              {table.getPageCount()}
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              Next
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => table.setPageIndex(table.getPageCount() - 1)}
              disabled={!table.getCanNextPage()}
            >
              Last
            </Button>

            <Input
              type="number"
              defaultValue={table.getState().pagination.pageIndex + 1}
              onChange={(e) => {
                const page = e.target.value ? Number(e.target.value) - 1 : 0;
                table.setPageIndex(page);
              }}
              className="w-16"
            />
          </div>
        </CardFooter>
      )}
    </Card>
  );
}

export default GenericDataTable;
