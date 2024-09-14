"use client";
import React, { useMemo, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { format } from "date-fns";
import { useQuery } from "@tanstack/react-query";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { apiClient } from "@/lib/apiClient";

interface ClusterRoleDetails {
  clusterRole: {
    metadata: {
      name: string;
      creationTimestamp: string;
      resourceVersion: string;
      uid: string;
    };
    rules: {
      apiGroups: string[];
      resources: string[];
      resourceNames: string[];
      verbs: string[];
    }[];
  };
}

const ClusterRoleDetailsPage = ({ params }: { params: { name: string } }) => {
  const { name } = params;
  const [currentPage, setCurrentPage] = useState(1);
  const [sortColumn, setSortColumn] = useState<string | null>(null);
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
  const [filter, setFilter] = useState("");
  const [expandedRows, setExpandedRows] = useState<Set<number>>(new Set());

  const {
    data: clusterRoleDetails,
    isLoading,
    error,
  } = useQuery<ClusterRoleDetails>({
    queryKey: ["clusterRoleDetails", name],
    queryFn: () => apiClient.getClusterRoleDetails(name),
  });

  const itemsPerPage = 10;

  const sortedAndFilteredRules = useMemo(() => {
    if (!clusterRoleDetails) return [];

    let filteredRules = clusterRoleDetails.clusterRole.rules.filter(
      (rule) =>
        rule.resources.some((resource) =>
          resource.toLowerCase().includes(filter.toLowerCase())
        ) ||
        rule.verbs.some((verb) =>
          verb.toLowerCase().includes(filter.toLowerCase())
        )
    );

    if (sortColumn) {
      filteredRules.sort((a, b) => {
        const aValue = a[sortColumn as keyof typeof a];
        const bValue = b[sortColumn as keyof typeof b];
        if (Array.isArray(aValue) && Array.isArray(bValue)) {
          return aValue.join(",").localeCompare(bValue.join(","));
        }
        return String(aValue).localeCompare(String(bValue));
      });

      if (sortDirection === "desc") {
        filteredRules.reverse();
      }
    }

    return filteredRules;
  }, [clusterRoleDetails, filter, sortColumn, sortDirection]);

  const paginatedRules = useMemo(() => {
    const startIndex = (currentPage - 1) * itemsPerPage;
    return sortedAndFilteredRules.slice(startIndex, startIndex + itemsPerPage);
  }, [sortedAndFilteredRules, currentPage]);

  const totalPages = Math.ceil(sortedAndFilteredRules.length / itemsPerPage);

  const toggleRowExpansion = (index: number) => {
    setExpandedRows((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(index)) {
        newSet.delete(index);
      } else {
        newSet.add(index);
      }
      return newSet;
    });
  };

  const handleSort = (column: string) => {
    if (sortColumn === column) {
      setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortColumn(column);
      setSortDirection("asc");
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {(error as Error).message}</div>;
  if (!clusterRoleDetails) return <div>No data available</div>;

  const { clusterRole } = clusterRoleDetails;

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40 p-4">
      <div className="grid gap-4 md:grid-cols-3">
        <Card className="overflow-hidden md:col-span-1">
          <CardHeader>
            <CardTitle>Cluster Role Summary</CardTitle>
          </CardHeader>
          <CardContent className="p-6 text-sm">
            <div className="grid gap-3">
              <div className="font-semibold flex items-center justify-between">
                Cluster Role Details
              </div>
              {clusterRole && (
                <ul className="grid gap-3">
                  <li className="flex items-center justify-between">
                    <span className="text-muted-foreground">Name:</span>
                    <span>{clusterRole.metadata.name}</span>
                  </li>
                  <li className="flex items-center justify-between">
                    <span className="text-muted-foreground">
                      Creation Date:
                    </span>
                    <span>
                      {format(
                        new Date(clusterRole.metadata.creationTimestamp),
                        "MM/dd - hh:mm:ss a"
                      )}
                    </span>
                  </li>
                  <li className="flex items-center justify-between">
                    <span className="text-muted-foreground">
                      Resource Version:
                    </span>
                    <span>{clusterRole.metadata.resourceVersion}</span>
                  </li>
                  <li className="flex items-center justify-between">
                    <span className="text-muted-foreground">Total Rules:</span>
                    <span>{clusterRole.rules.length}</span>
                  </li>
                </ul>
              )}
            </div>
          </CardContent>
        </Card>
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex items-center justify-between">
              <Input
                placeholder="Filter rules..."
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="max-w-sm"
              />
              <Select onValueChange={(value) => setCurrentPage(1)}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Items per page" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10">10 per page</SelectItem>
                  <SelectItem value="20">20 per page</SelectItem>
                  <SelectItem value="50">50 per page</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="rounded-md border">
              <Table>
                <TableHeader className="sticky top-0 bg-background">
                  <TableRow>
                    <TableHead className="w-[50px]"></TableHead>
                    <TableHead
                      onClick={() => handleSort("apiGroups")}
                      className="cursor-pointer"
                    >
                      API Groups{" "}
                      {sortColumn === "apiGroups" &&
                        (sortDirection === "asc" ? "↑" : "↓")}
                    </TableHead>
                    <TableHead
                      onClick={() => handleSort("resources")}
                      className="cursor-pointer"
                    >
                      Resources{" "}
                      {sortColumn === "resources" &&
                        (sortDirection === "asc" ? "↑" : "↓")}
                    </TableHead>
                    <TableHead
                      onClick={() => handleSort("verbs")}
                      className="cursor-pointer"
                    >
                      Verbs{" "}
                      {sortColumn === "verbs" &&
                        (sortDirection === "asc" ? "↑" : "↓")}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {paginatedRules.map((rule, index) => (
                    <React.Fragment key={index}>
                      <TableRow>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => toggleRowExpansion(index)}
                          >
                            {expandedRows.has(index) ? "−" : "+"}
                          </Button>
                        </TableCell>
                        <TableCell>{rule.apiGroups.join(", ")}</TableCell>
                        <TableCell>
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger>
                                {rule.resources.slice(0, 3).join(", ")}
                                {rule.resources.length > 3 && "..."}
                              </TooltipTrigger>
                              <TooltipContent>
                                <p>{rule.resources.join(", ")}</p>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </TableCell>
                        <TableCell>
                          {rule.verbs.map((verb) => (
                            <Badge
                              key={verb}
                              variant="secondary"
                              className="mr-1 mb-1"
                            >
                              {verb}
                            </Badge>
                          ))}
                        </TableCell>
                      </TableRow>
                      {expandedRows.has(index) && (
                        <TableRow>
                          <TableCell colSpan={4}>
                            {/* <Accordion type="single" collapsible>
                              <AccordionItem value="details">
                                <AccordionTrigger>
                                  Rule Details
                                </AccordionTrigger>
                                <AccordionContent> */}
                            <div className="grid grid-cols-2 gap-2">
                              <div>
                                <strong>API Groups:</strong>
                                <ul>
                                  {rule.apiGroups.map((group, i) => (
                                    <li key={i}>{group}</li>
                                  ))}
                                </ul>
                              </div>
                              <div>
                                <strong>Resources:</strong>
                                <ul>
                                  {rule.resources.map((resource, i) => (
                                    <li key={i}>{resource}</li>
                                  ))}
                                </ul>
                              </div>
                              <div>
                                <strong>Verbs:</strong>
                                <ul>
                                  {rule.verbs.map((verb, i) => (
                                    <li key={i}>{verb}</li>
                                  ))}
                                </ul>
                              </div>
                              {rule.resourceNames && (
                                <div>
                                  <strong>Resource Names:</strong>
                                  <ul>
                                    {rule.resourceNames.map((name, i) => (
                                      <li key={i}>{name}</li>
                                    ))}
                                  </ul>
                                </div>
                              )}
                            </div>
                            {/* </AccordionContent>
                              </AccordionItem>
                            </Accordion> */}
                          </TableCell>
                        </TableRow>
                      )}
                    </React.Fragment>
                  ))}
                </TableBody>
              </Table>
            </div>
            <div className="mt-4 flex items-center justify-between">
              <div>
                Showing {(currentPage - 1) * itemsPerPage + 1} to{" "}
                {Math.min(
                  currentPage * itemsPerPage,
                  sortedAndFilteredRules.length
                )}{" "}
                of {sortedAndFilteredRules.length} rules
              </div>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() =>
                    setCurrentPage((prev) => Math.max(prev - 1, 1))
                  }
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() =>
                    setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                  }
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ClusterRoleDetailsPage;
