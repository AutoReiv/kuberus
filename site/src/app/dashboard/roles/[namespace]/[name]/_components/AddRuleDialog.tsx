import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { verbs } from "@/interfaces/verbs";
import { z } from "zod";
import { Badge } from "@/components/ui/badge";

const newRuleSchema = z.object({
  resources: z.string().min(1, "You must select a resource"),
  verbs: z.array(z.string()).min(1, "You must select at least one verb"),
});
type NewRuleFormValues = z.infer<typeof newRuleSchema>;

const AddRuleDialog = ({
  isOpen,
  onClose,
  form,
  onSubmit,
  resources,
  existingResources,
}) => {
  const availableResources = resources.resources.filter(
    (resource) => !existingResources.includes(resource)
  );

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Rule</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="resources"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Resources</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select a resource" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {availableResources.map((resource) => (
                        <SelectItem key={resource} value={resource}>
                          {resource}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    Choose the resource for this rule.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="verbs"
              render={() => (
                <FormItem>
                  <div className="mb-4">
                    <FormLabel className="text-base">Verbs</FormLabel>
                    <FormDescription>
                      Select the verbs for this rule.
                    </FormDescription>
                  </div>
                  {verbs.map((verb) => (
                    <FormField
                      key={verb.name}
                      control={form.control}
                      name="verbs"
                      render={({ field }) => {
                        const isSelected = field.value?.includes(verb.name);
                        return (
                          <Badge
                            key={verb.name}
                            variant={isSelected ? "success" : "secondary"}
                            className="mr-2 mb-2 cursor-pointer"
                            onClick={() => {
                              const newValue = isSelected
                                ? field.value.filter((v) => v !== verb.name)
                                : [...field.value, verb.name];
                              field.onChange(newValue);
                            }}
                          >
                            {verb.name}
                          </Badge>
                        );
                      }}
                    />
                  ))}
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit">Add Rule</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default AddRuleDialog;
