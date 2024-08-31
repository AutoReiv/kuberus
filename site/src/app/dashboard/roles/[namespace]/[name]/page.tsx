"use client";

import { useQuery } from "@tanstack/react-query";
import { Copy } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { format } from "date-fns";
import { useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";

interface RoleDetails {
  role: {
    metadata: {
      creationTimestamp: string;
      managedFields: [
        {
          apiVersion: string;
          fieldsType: string;
          fieldsV1: {
            f: {
              rules: [
                {
                  apiGroups: string[];
                  resources: string[];
                  resourceNames: string[];
                  verbs: string[];
                }
              ];
            };
          };
          manager: string;
          operation: string;
          time: string;
        }
      ];
      name: string;
      namespace: string;
      resourceVersion: string;
      uid: string;
    };
    rules: [
      {
        apiGroups: string[];
        resources: string[];
        resourceNames: string[];
        verbs: string[];
      }
    ];
  };
}

const newRuleSchema = z.object({
  resources: z.string().min(1, "Resources are required"),
  verbs: z.array(z.string()).refine((value) => value.some((item) => item), {
    message: "You have to select at least one item.",
  }),
});

type NewRuleFormValues = z.infer<typeof newRuleSchema>;

const verbs = [
  { name: "Create" },
  { name: "Delete" },
  { name: "Get" },
  { name: "List" },
  { name: "Patch" },
  { name: "Update" },
  { name: "Watch" },
];

const fetchRoleDetails = async (namespace: string, name: string) => {
  const URL = `http://localhost:8080/api/roles/details?roleName=${name}&namespace=${namespace}`;
  const response = await fetch(URL, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });
  if (!response.ok) {
    throw new Error("Failed to fetch role details");
  }
  const data = await response.json();
  console.log(data, " DATA ****");
  return data;
};

const RoleDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const form = useForm<NewRuleFormValues>({
    resolver: zodResolver(newRuleSchema),
    defaultValues: {
      resources: "",
      verbs: ["create"],
    },
  });

  const {
    data: roleDetails,
    isLoading,
    error,
  } = useQuery<RoleDetails, Error>({
    queryKey: ["roleDetails", namespace, name],
    queryFn: () => fetchRoleDetails(namespace, name),
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  const onSubmit = (data: NewRuleFormValues) => {
    console.log(data);
    setIsDialogOpen(false);
    form.reset();
  };

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40">
      {/* <aside className="fixed inset-y-0 left-0 z-10 hidden w-14 flex-col border-r bg-background sm:flex">
        <nav className="flex flex-col items-center gap-4 px-2 sm:py-5">
          <Link
            href="#"
            className="group flex h-9 w-9 shrink-0 items-center justify-center gap-2 rounded-full bg-primary text-lg font-semibold text-primary-foreground md:h-8 md:w-8 md:text-base"
          >
            <Package2 className="h-4 w-4 transition-all group-hover:scale-110" />
            <span className="sr-only">Acme Inc</span>
          </Link>
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <Home className="h-5 w-5" />
                <span className="sr-only">Dashboard</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Dashboard</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent text-accent-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <ShoppingCart className="h-5 w-5" />
                <span className="sr-only">Orders</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Orders</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <Package className="h-5 w-5" />
                <span className="sr-only">Products</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Products</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <Users2 className="h-5 w-5" />
                <span className="sr-only">Customers</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Customers</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <LineChart className="h-5 w-5" />
                <span className="sr-only">Analytics</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Analytics</TooltipContent>
          </Tooltip>
        </nav>
        <nav className="mt-auto flex flex-col items-center gap-4 px-2 sm:py-5">
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="#"
                className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
              >
                <Settings className="h-5 w-5" />
                <span className="sr-only">Settings</span>
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">Settings</TooltipContent>
          </Tooltip>
        </nav>
      </aside> */}
      <div className="flex flex-col sm:gap-4">
        <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 md:gap-8 lg:grid-cols-3 xl:grid-cols-3">
          <div>
            <Card className="overflow-hidden" x-chunk="dashboard-05-chunk-4">
              <CardHeader className="flex flex-row items-start bg-muted/50">
                <div className="grid gap-0.5">
                  <CardTitle className="group flex items-center gap-2 text-lg">
                    {roleDetails.role.metadata.name}
                    <Button
                      size="icon"
                      variant="outline"
                      className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                    >
                      <Copy className="h-3 w-3" />
                      <span className="sr-only">Copy Order ID</span>
                    </Button>
                  </CardTitle>
                  <CardDescription>
                    UID: {roleDetails.role.metadata.uid}
                  </CardDescription>
                </div>
              </CardHeader>
              <CardContent className="p-6 text-sm">
                <div className="grid gap-3">
                  <div className="font-semibold">Role Details</div>
                  <ul className="grid gap-3">
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">Name:</span>
                      <span>{roleDetails.role.metadata.name}</span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">Namespace:</span>
                      <span>{roleDetails.role.metadata.namespace}</span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">
                        Creation Date:
                      </span>
                      <span>
                        {format(
                          new Date(roleDetails.role.metadata.creationTimestamp),
                          "MM/dd - hh:mm:ss a"
                        )}
                      </span>
                    </li>
                    <li className="flex items-center justify-between">
                      <span className="text-muted-foreground">
                        Resource Version:
                      </span>
                      <span>{roleDetails.role.metadata.resourceVersion}</span>
                    </li>
                  </ul>
                  <Separator></Separator>
                </div>
              </CardContent>
              {/* <CardFooter className="flex flex-row items-center border-t bg-muted/50 px-6 py-3">
                <div className="text-xs text-muted-foreground">
                  Updated <time dateTime="2023-11-23">November 23, 2023</time>
                </div>
                <Pagination className="ml-auto mr-0 w-auto">
                  <PaginationContent>
                    <PaginationItem>
                      <Button size="icon" variant="outline" className="h-6 w-6">
                        <ChevronLeft className="h-3.5 w-3.5" />
                        <span className="sr-only">Previous Order</span>
                      </Button>
                    </PaginationItem>
                    <PaginationItem>
                      <Button size="icon" variant="outline" className="h-6 w-6">
                        <ChevronRight className="h-3.5 w-3.5" />
                        <span className="sr-only">Next Order</span>
                      </Button>
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              </CardFooter> */}
            </Card>
          </div>
          <div className="grid auto-rows-max items-start gap-4 md:gap-8 lg:col-span-2">
            <Tabs defaultValue="week">
              <div className="flex items-center">
                <TabsList>
                  <TabsTrigger value="week">Week</TabsTrigger>
                  <TabsTrigger value="month">Month</TabsTrigger>
                  <TabsTrigger value="year">Year</TabsTrigger>
                </TabsList>
                {/* <div className="ml-auto flex items-center gap-2">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-7 gap-1 text-sm"
                      >
                        <ListFilter className="h-3.5 w-3.5" />
                        <span className="sr-only sm:not-sr-only">Filter</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>Filter by</DropdownMenuLabel>
                      <DropdownMenuSeparator />
                      <DropdownMenuCheckboxItem checked>
                        Fulfilled
                      </DropdownMenuCheckboxItem>
                      <DropdownMenuCheckboxItem>
                        Declined
                      </DropdownMenuCheckboxItem>
                      <DropdownMenuCheckboxItem>
                        Refunded
                      </DropdownMenuCheckboxItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-7 gap-1 text-sm"
                  >
                    <File className="h-3.5 w-3.5" />
                    <span className="sr-only sm:not-sr-only">Export</span>
                  </Button>
                </div> */}
              </div>
              <TabsContent value="week">
                <Card x-chunk="dashboard-05-chunk-3">
                  <CardHeader className="px-7 flex-row items-center justify-between">
                    <div>
                      <CardTitle>Rules</CardTitle>
                      <CardDescription>
                        Here are the rules for this role.
                      </CardDescription>
                    </div>
                    <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                      <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                          <Plus className="h-4 w-4 mr-2" />
                          Add Rule
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Add New Rule</DialogTitle>
                        </DialogHeader>
                        <Form {...form}>
                          <form
                            onSubmit={form.handleSubmit(onSubmit)}
                            className="space-y-8"
                          >
                            <FormField
                              control={form.control}
                              name="resources"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel>Resources</FormLabel>
                                  <FormControl>
                                    <Input
                                      placeholder="e.g. pods, services"
                                      {...field}
                                    />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                            <FormField
                              control={form.control}
                              name="verbs"
                              render={({ field }) => (
                                <FormItem className="columns-2">
                                  {verbs.map((badge) => (
                                    <FormField
                                      key={badge.name}
                                      control={form.control}
                                      name="verbs"
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
                            <Button type="submit">Add Rule</Button>
                          </form>
                        </Form>
                      </DialogContent>
                    </Dialog>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead className="w-[calc(10%)]">#</TableHead>
                          <TableHead className="hidden sm:table-cell w-1/3">
                            Resources
                          </TableHead>
                          <TableHead className="hidden sm:table-cell w-1/3">
                            Verbs
                          </TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {roleDetails.role.rules.map((rule, index) => (
                          <TableRow key={index} className="bg-accent">
                            <TableCell className="font-medium">
                              {index + 1}
                            </TableCell>
                            <TableCell>
                              {rule.resources.map((resource) => (
                                <Badge key={resource} variant="default">
                                  {resource}
                                </Badge>
                              ))}
                            </TableCell>
                            <TableCell className="flex gap-2 flex-wrap">
                              {rule.verbs.map((verb) => (
                                <Badge key={verb} variant="default">
                                  {verb}
                                </Badge>
                              ))}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </div>
  );
};

export default RoleDetailsPage;
