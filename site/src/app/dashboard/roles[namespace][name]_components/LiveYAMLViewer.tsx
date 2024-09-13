import React, { useEffect, useState } from "react";
import yaml from "js-yaml";
import { Button } from "@/components/ui/button";
import { CheckCircle2, Copy, Edit, Save } from "lucide-react";
import { toast } from "sonner";
import MonacoEditor from "@monaco-editor/react";
import * as monaco from 'monaco-editor';

interface LiveYAMLViewerProps {
  rules;
  metadata: {
    name: string;
    namespace: string;
  };
  onUpdate: (updatedRules: any, updatedMetadata: any) => void;
}


monaco.editor.defineTheme('custom-dark', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: 'comment', foreground: '7f848e' },
      { token: 'keyword', foreground: 'c678dd' },
      { token: 'identifier', foreground: 'e06c75' },
      { token: 'string', foreground: '98c379' },
      // Add more token rules as needed
    ],
    colors: {
      'editor.background': '#282c34',
      'editor.foreground': '#abb2bf',
      'editor.lineHighlightBackground': '#2c313c',
      'editorCursor.foreground': '#528bff',
      'editorWhitespace.foreground': '#3b4048',
      'editorIndentGuide.background': '#3b4048',
      'editorIndentGuide.activeBackground': '#7f848e',
      // Add more color rules as needed
    },
  });

  
const LiveYAMLViewer: React.FC<LiveYAMLViewerProps> = ({ rules, metadata, onUpdate }) => {
  const [yamlContent, setYamlContent] = useState("");
  const [isEditing, setIsEditing] = useState(false);
  const [editedYaml, setEditedYaml] = useState("");

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
      setEditedYaml(yamlString);
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

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSave = () => {
    try {
      const updatedConfig = yaml.load(editedYaml) as any;
      onUpdate(updatedConfig.rules, updatedConfig.metadata);
      setIsEditing(false);
      toast(
        <div className="flex items-center justify-start gap-4">
          <CheckCircle2 className="text-green-500" />
          <span>Changes saved successfully.</span>
        </div>
      );
    } catch (error) {
      toast(
        <div className="flex items-center justify-start gap-4">
          <CheckCircle2 className="text-red-500" />
          <span>Invalid YAML. Please check your changes.</span>
        </div>
      );
    }
  };

  return (
    <div className="live-yaml-viewer mt-2">
      <div className="flex justify-between items-center mb-2">
        <h3 className="text-lg font-semibold">Live YAML View</h3>
        <div>
          {!isEditing && (
            <Button onClick={handleEdit} variant="outline" size="sm" className="mr-2">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
          )}
          {isEditing && (
            <Button onClick={handleSave} variant="outline" size="sm" className="mr-2">
              <Save className="h-4 w-4 mr-2" />
              Save
            </Button>
          )}
          <Button onClick={copyToClipboard} variant="outline" size="sm">
            <Copy className="h-4 w-4 mr-2" />
            Copy
          </Button>
        </div>
      </div>
      <MonacoEditor
        height="400px"
        language="yaml"
        theme="custom-dark"
        value={isEditing ? editedYaml : yamlContent}
        options={{ readOnly: !isEditing }}
        onChange={(value) => isEditing && setEditedYaml(value || "")}
      />
    </div>
  );
};

export default LiveYAMLViewer;
