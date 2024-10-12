import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import YamlEditor from "@/components/YamlEditor";
import { RoleDetails } from "@/interfaces/roleDetails";
import { format } from "date-fns";
import { Copy } from "lucide-react";
import yaml from "js-yaml";

const RoleDetailsCard: React.FC<{
  roleDetails: RoleDetails;
  onDuplicate: () => void;
  handleYamlUpdate: (updatedYaml: string) => void;
}> = ({ roleDetails, onDuplicate, handleYamlUpdate }) => (
  <Card className="overflow-hidden">
    <CardContent className="p-6 text-sm">
      <div className="grid gap-3">
        <div className="font-semibold flex items-center justify-between">
          Role Details
          <Button onClick={onDuplicate} variant="outline" size="sm">
            <Copy className="h-4 w-4 mr-2" />
            Duplicate Role
          </Button>
        </div>
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
            <span className="text-muted-foreground">Creation Date:</span>
            <span>
              {format(
                new Date(roleDetails.role.metadata.creationTimestamp),
                "MM/dd - hh:mm:ss a"
              )}
            </span>
          </li>
          <li className="flex items-center justify-between">
            <span className="text-muted-foreground">Resource Version:</span>
            <span>{roleDetails.role.metadata.resourceVersion}</span>
          </li>
        </ul>
        <Separator />
        <YamlEditor
          initialContent={yaml.dump({
            apiVersion: "rbac.authorization.k8s.io/v1",
            kind: "Role",
            metadata: {
              name: roleDetails.role.metadata.name,
              namespace: roleDetails.role.metadata.namespace,
            },
            rules: roleDetails.role.rules,
          })}
          onSave={handleYamlUpdate}
          readOnly={false}
          enableDiff={false}
          enableReset={false}
          enableFormat={false}
        />
      </div>
    </CardContent>
  </Card>
);
export default RoleDetailsCard;
