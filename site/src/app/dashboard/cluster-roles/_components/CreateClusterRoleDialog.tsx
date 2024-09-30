import React from "react";
import { Button } from "@/components/ui/button";
import { Trash } from "lucide-react";
import { ResponsiveDialog } from "@/components/ResponsiveDialog";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useFieldArray, useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useCreateClusterRole } from "@/hooks/useClusterRoles";

const clusterRoleSchema = z.object({
  name: z.string().min(1, "Name is required"),
  rules: z
    .array(
      z.object({
        apiGroups: z.array(z.string()),
        resources: z.array(z.string()),
        verbs: z.array(z.string()),
      })
    )
    .min(1, "At least one rule is required"),
});

type ClusterRoleFormValues = z.infer<typeof clusterRoleSchema>;

const CreateClusterRoleDialog = ({
  isOpen,
  setIsOpen,
}: {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}) => {
  const createClusterRoleMutation = useCreateClusterRole({
    onSuccess: () => {
      setIsOpen(false);
    },
  });

  const form = useForm<ClusterRoleFormValues>({
    resolver: zodResolver(clusterRoleSchema),
    defaultValues: {
      name: "",
      rules: [{ apiGroups: [], resources: [], verbs: [] }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "rules",
  });

  const onSubmit = async (data: ClusterRoleFormValues) => {
    const payload = {
      apiVersion: "rbac.authorization.k8s.io/v1",
      kind: "ClusterRole",
      metadata: {
        name: data.name,
      },
      rules: data.rules.map((rule) => ({
        apiGroups: rule.apiGroups,
        resources: rule.resources,
        verbs: rule.verbs,
      })),
    };
    await createClusterRoleMutation.mutateAsync(payload);
  };

  return (
    <ResponsiveDialog
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title="Create New Cluster Role"
      className="!max-w-none w-[60%] h-[60]"
    >
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <div className="flex items-center gap-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem className="flex-1">
                  <FormLabel>Cluster Role Name *</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter cluster role name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          {fields.map((field, index) => (
            <div key={field.id} className="space-y-4 p-4 border rounded-md">
              <div className="flex justify-between items-center">
                <h4 className="text-lg font-semibold">Rule {index + 1}</h4>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => remove(index)}
                >
                  <Trash className="h-4 w-4" />
                </Button>
              </div>

              <FormField
                control={form.control}
                name={`rules.${index}.apiGroups`}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>API Groups *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Enter API groups (comma-separated)"
                        {...field}
                        onChange={(e) =>
                          field.onChange(e.target.value.split(","))
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name={`rules.${index}.resources`}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Resources *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Enter resources (comma-separated)"
                        {...field}
                        onChange={(e) =>
                          field.onChange(e.target.value.split(","))
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name={`rules.${index}.verbs`}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Verbs *</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={(value) =>
                          field.onChange([...field.value, value])
                        }
                        value=""
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Select verbs" />
                        </SelectTrigger>
                        <SelectContent>
                          {[
                            "get",
                            "list",
                            "watch",
                            "create",
                            "update",
                            "patch",
                            "delete",
                          ].map((verb) => (
                            <SelectItem key={verb} value={verb}>
                              {verb}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <div className="flex flex-wrap gap-2 mt-2">
                      {field.value.map((verb, verbIndex) => (
                        <div
                          key={verbIndex}
                          className="bg-secondary text-secondary-foreground px-2 py-1 rounded-md flex items-center gap-2"
                        >
                          {verb}
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              field.onChange(
                                field.value.filter((_, i) => i !== verbIndex)
                              )
                            }
                          >
                            <Trash className="h-3 w-3" />
                          </Button>
                        </div>
                      ))}
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          ))}

          <Button
            type="button"
            variant="outline"
            onClick={() => append({ apiGroups: [], resources: [], verbs: [] })}
          >
            Add Rule
          </Button>

          <Button type="submit">Create Cluster Role</Button>
        </form>
      </Form>
    </ResponsiveDialog>
  );
};

export default CreateClusterRoleDialog;
