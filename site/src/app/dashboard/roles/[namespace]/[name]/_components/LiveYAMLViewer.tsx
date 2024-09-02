import React, { useEffect, useState } from "react";
import yaml from "js-yaml";
import { Button } from "@/components/ui/button";
import { CheckCircle2, Copy } from "lucide-react";
import { toast } from "sonner";

interface LiveYAMLViewerProps {
  rules;
  metadata: {
    name: string;
    namespace: string;
  };
}

const LiveYAMLViewer: React.FC<LiveYAMLViewerProps> = ({ rules, metadata }) => {
  const [yamlContent, setYamlContent] = useState("");

  useEffect(() => {
    const updateYAML = () => {
      const roleConfig = {
        apiVersion: "rbac.authorization.k8s.io/v1",
        kind: "Role",
        metadata: metadata,
        rules: rules,
      };

      const yamlString = yaml.dump(roleConfig);
      setYamlContent(yamlString);
    };

    updateYAML();
  }, [rules, metadata]);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(yamlContent).then(() => {
      toast(
        <div className="flex items-center justify-start gap-4">
          <CheckCircle2 className="text-green-500" />
          <span>YAML content has been copied to your clipboard.</span>
        </div>
      );
    });
  };

  return (
    <div className="live-yaml-viewer mt-2">
      <div className="flex justify-between items-center mb-2">
        <h3 className="text-lg font-semibold">Live YAML View</h3>
        <Button onClick={copyToClipboard} variant="outline" size="sm">
          <Copy className="h-4 w-4 mr-2" />
          Copy
        </Button>
      </div>
      <pre className="bg-muted p-4 rounded-md overflow-auto dark:bg-muted-foreground/10 dark:text-muted-foreground">
        {yamlContent}
      </pre>
    </div>
  );
};

export default LiveYAMLViewer;
