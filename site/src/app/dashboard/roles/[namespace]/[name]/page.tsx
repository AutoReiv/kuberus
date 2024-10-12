"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { z } from "zod";
import yaml from "js-yaml";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { apiClient } from "@/lib/apiClient";
import { Rules } from "@/interfaces/rules";
import { RoleDetails } from "@/interfaces/roleDetails";
import { useResources } from "@/hooks/useResources";
import { useRoleDetails } from "@/hooks/useRoleDetails";
import { verbs } from "@/interfaces/verbs";
import { Badge } from "@/components/ui/badge";
import { Plus } from "lucide-react";
import GenericDataTable from "@/components/GenericDataTable";
import AddRuleDialog from "./_components/AddRuleDialog";
import DuplicateRoleDialog from "./_components/DuplicateRoleDialog";
import RoleDetailsCard from "./_components/RoleDetailsCard";

const newRuleSchema = z.object({
  resources: z.string().min(1, "You must select a resource"),
  verbs: z.array(z.string()).min(1, "You must select at least one verb"),
});

type NewRuleFormValues = z.infer<typeof newRuleSchema>;

const RoleDetailsPage = ({
  params,
}: {
  params: { namespace: string; name: string };
}) => {
  const { namespace, name } = params;
  const router = useRouter();
  const [isAddRuleDialogOpen, setIsAddRuleDialogOpen] = useState(false);
  const [isDuplicateDialogOpen, setIsDuplicateDialogOpen] = useState(false);

  const {
    data: roleDetails,
    isLoading,
    error,
    refetch: refetchRoleDetails,
  } = useRoleDetails(namespace, name);
  const { data: resources } = useResources();

  const form = useForm<NewRuleFormValues>({
    resolver: zodResolver(newRuleSchema),
    defaultValues: {
      resources: "",
      verbs: [],
    },
  });

  const columns = [
    {
      accessorKey: "resources",
      header: "Resources",
      cell: ({ row }) => (
        <div>
          {row.original.resources.map((resource) => (
            <Badge key={resource} variant="default">
              {resource}
            </Badge>
          ))}
        </div>
      ),
      accessorFn: (row) => row.resources.join(", "),
    },
    {
      accessorKey: "verbs",
      header: "Verbs",
      cell: ({ row }) => (
        <div className="flex gap-2 flex-wrap">
          {verbs.map((verb) => (
            <Badge
              key={verb.name}
              variant={
                row.original.verbs.includes(verb.name) ? "success" : "secondary"
              }
              className={`cursor-pointer ${
                row.original.verbs.includes(verb.name) ? "" : "opacity-50"
              }`}
              onClick={() => toggleVerb(row.index, verb.name)}
            >
              {verb.name}
            </Badge>
          ))}
        </div>
      ),
    },
  ];

  const toggleVerb = async (ruleIndex: number, verb: string) => {
    const updatedRules = [...roleDetails.role.rules];
    const ruleVerbs = updatedRules[ruleIndex].verbs;
    updatedRules[ruleIndex].verbs = ruleVerbs.includes(verb)
      ? ruleVerbs.filter((v) => v !== verb)
      : [...ruleVerbs, verb];

    try {
      await updateRole(updatedRules);
      await refetchRoleDetails();
      toast.success("Rule updated successfully");
    } catch (error) {
      console.error("Failed to update role:", error);
      toast.error("Failed to update rule");
    }
  };

  const onSubmit = async (data: NewRuleFormValues) => {
    const existingResources = roleDetails.role.rules.flatMap(
      (rule) => rule.resources
    );

    if (existingResources.includes(data.resources)) {
      toast.error("This resource already exists in the rules.");
      return;
    }

    const newRule: Rules = {
      apiGroups: [""],
      resources: [data.resources],
      resourceNames: [],
      verbs: data.verbs,
    };

    const updatedRules = [...roleDetails.role.rules, newRule];

    try {
      await updateRole(updatedRules);
      refetchRoleDetails();
      setIsAddRuleDialogOpen(false);
      form.reset();
      toast.success("New rule added successfully");
    } catch (error) {
      console.error("Failed to add new rule:", error);
      toast.error("Failed to add new rule");
    }
  };
  const updateRole = async (updatedRules: Rules[]) => {
    const roleData = {
      metadata: {
        name: name,
        namespace: namespace,
      },
      rules: updatedRules,
    };

    await apiClient.updateRole(namespace, name, roleData);
  };

  const duplicateRole = async (newNamespace: string, newName: string) => {
    const newRoleData = {
      metadata: {
        ...(() => {
          const { resourceVersion, ...rest } = roleDetails.role.metadata;
          return rest;
        })(),
        name: newName,
        namespace: newNamespace,
      },
      rules: roleDetails.role.rules,
    };

    try {
      await apiClient.createRoles(newRoleData);
      toast.success(`Successfully duplicated ${newName} in ${newNamespace}`);
      router.push("/dashboard/roles");
    } catch (error) {
      toast.error(`Failed to duplicate role: ${error}`);
    }
  };

  const handleYamlUpdate = async (updatedYaml: string) => {
    try {
      const updatedRole = yaml.load(updatedYaml) as RoleDetails;
      await apiClient.updateRole(namespace, name, updatedRole);
      refetchRoleDetails();
      toast.success("Role updated successfully");
    } catch (error) {
      toast.error("Failed to update role");
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="flex min-h-screen w-full flex-col bg-muted/40">
      <div className="flex flex-col sm:gap-4">
        <main className="grid flex-1 items-start gap-4 p-4 sm:px-6 md:gap-8 lg:grid-cols-3 xl:grid-cols-3">
          <RoleDetailsCard
            roleDetails={roleDetails}
            onDuplicate={() => setIsDuplicateDialogOpen(true)}
            handleYamlUpdate={handleYamlUpdate}
          />
          <div className="grid auto-rows-max items-start gap-4 md:gap-8 lg:col-span-2">
            <Tabs defaultValue="rules">
              <TabsList>
                <TabsTrigger value="rules">Rules</TabsTrigger>
              </TabsList>
              <TabsContent value="rules">
                <GenericDataTable
                  data={roleDetails.role.rules}
                  columns={columns}
                  title="Rules"
                  description="Here are the rules for this role."
                  enableSorting={true}
                  enableFiltering={true}
                  enablePagination={true}
                  enableGridView={false}
                  enableQuickActions={false}
                  actionButton={
                    <Button
                      onClick={() => setIsAddRuleDialogOpen(true)}
                      variant="outline"
                      size="sm"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Add Rule
                    </Button>
                  }
                />
              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
      <AddRuleDialog
        isOpen={isAddRuleDialogOpen}
        onClose={() => setIsAddRuleDialogOpen(false)}
        form={form}
        onSubmit={onSubmit}
        resources={resources}
        existingResources={roleDetails.role.rules.flatMap(
          (rule) => rule.resources
        )}
      />
      <DuplicateRoleDialog
        isOpen={isDuplicateDialogOpen}
        onClose={() => setIsDuplicateDialogOpen(false)}
        onDuplicate={duplicateRole}
      />
    </div>
  );
};

export default RoleDetailsPage;
