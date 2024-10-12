import React, { useState, useEffect, useRef } from "react";
import yaml from "js-yaml";
import MonacoEditor, { DiffEditor } from "@monaco-editor/react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import {
  Copy,
  Edit,
  Save,
  RotateCcw,
  FileSymlink,
  Download,
  SplitSquareHorizontal,
} from "lucide-react";

interface YamlEditorProps {
  initialContent: string;
  onSave?: (content: string) => void;
  readOnly?: boolean;
  enableEdit?: boolean;
  enableSave?: boolean;
  enableReset?: boolean;
  enableFormat?: boolean;
  enableCopy?: boolean;
  enableDiff?: boolean;
  enableExport?: boolean;
}

const YamlEditor: React.FC<YamlEditorProps> = ({
  initialContent,
  onSave,
  readOnly = false,
  enableEdit = true,
  enableSave = true,
  enableReset = true,
  enableFormat = true,
  enableCopy = true,
  enableDiff = true,
  enableExport = true,
}) => {
  const [content, setContent] = useState(initialContent);
  const [isEditing, setIsEditing] = useState(!readOnly);
  const [showDiff, setShowDiff] = useState(false);
  const [history, setHistory] = useState<string[]>([initialContent]);
  const [historyIndex, setHistoryIndex] = useState(0);
  const { theme } = useTheme();
  const editorRef = useRef(null);

  useEffect(() => {
    setContent(initialContent);
    setHistory([initialContent]);
    setHistoryIndex(0);
  }, [initialContent]);

  const handleEdit = () => setIsEditing(true);

  const handleSave = () => {
    try {
      yaml.load(content); // Validate YAML
      setIsEditing(false);
      onSave?.(content);
      addToHistory(content);
      toast.success("YAML saved successfully");
    } catch (error) {
      toast.error("Invalid YAML. Please check your changes.");
    }
  };

  const handleReset = () => {
    setContent(initialContent);
    toast.info("YAML reset to original content");
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(content).then(() => {
      toast.success("YAML content copied to clipboard");
    });
  };

  const formatYaml = () => {
    try {
      const formattedYaml = yaml.dump(yaml.load(content), { indent: 2 });
      setContent(formattedYaml);
      toast.success("YAML formatted successfully");
    } catch (error) {
      toast.error("Failed to format YAML. Please check for syntax errors.");
    }
  };

  const toggleDiffView = () => setShowDiff(!showDiff);

  const addToHistory = (newContent: string) => {
    setHistory((prev) => [...prev.slice(0, historyIndex + 1), newContent]);
    setHistoryIndex((prev) => prev + 1);
  };

  const exportYaml = () => {
    const blob = new Blob([content], { type: "text/yaml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "rbac_config.yaml";
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="yaml-editor">
      <div className="flex justify-between items-center mb-2">
        <h3 className="text-lg font-semibold">YAML Editor</h3>
        <div>
          {!isEditing && !readOnly && enableEdit && (
            <Button
              onClick={handleEdit}
              variant="outline"
              size="sm"
              className="mr-2"
            >
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
          )}
          {isEditing && (
            <>
              {enableSave && (
                <Button
                  onClick={handleSave}
                  variant="outline"
                  size="sm"
                  className="mr-2"
                >
                  <Save className="h-4 w-4 mr-2" />
                  Save
                </Button>
              )}
              {enableReset && (
                <Button
                  onClick={handleReset}
                  variant="outline"
                  size="sm"
                  className="mr-2"
                >
                  <RotateCcw className="h-4 w-4 mr-2" />
                  Reset
                </Button>
              )}
              {enableFormat && (
                <Button
                  onClick={formatYaml}
                  variant="outline"
                  size="sm"
                  className="mr-2"
                >
                  <FileSymlink className="h-4 w-4 mr-2" />
                  Format
                </Button>
              )}
            </>
          )}
          {enableCopy && (
            <Button
              onClick={copyToClipboard}
              variant="outline"
              size="sm"
              className="mr-2"
            >
              <Copy className="h-4 w-4 mr-2" />
              Copy
            </Button>
          )}
          {enableDiff && (
            <Button
              onClick={toggleDiffView}
              variant="outline"
              size="sm"
              className="mr-2"
            >
              <SplitSquareHorizontal className="h-4 w-4 mr-2" />
              {showDiff ? "Hide Diff" : "Show Diff"}
            </Button>
          )}
          {enableExport && (
            <Button
              onClick={exportYaml}
              variant="outline"
              size="sm"
              className="mr-2"
            >
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          )}
        </div>
      </div>
      <div className="flex">
        {showDiff ? (
          <DiffEditor
            height="400px"
            language="yaml"
            theme={theme === "dark" ? "vs-dark" : "light"}
            original={initialContent}
            modified={content}
            options={{ readOnly: true }}
          />
        ) : (
          <MonacoEditor
            height="400px"
            language="yaml"
            theme={theme === "dark" ? "vs-dark" : "light"}
            value={content}
            options={{
              readOnly: !isEditing,
              minimap: { enabled: false },
              automaticLayout: true,
            }}
            onChange={(value) => setContent(value || "")}
            onMount={(editor) => {
              editorRef.current = editor;
            }}
          />
        )}
      </div>
    </div>
  );
};

export default YamlEditor;
