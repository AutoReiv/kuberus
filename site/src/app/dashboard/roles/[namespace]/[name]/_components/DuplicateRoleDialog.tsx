import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

const DuplicateRoleDialog: React.FC<{
  isOpen: boolean;
  onClose: () => void;
  onDuplicate: (newNamespace: string, newName: string) => void;
}> = ({ isOpen, onClose, onDuplicate }) => {
  const [newNamespace, setNewNamespace] = useState("");
  const [newName, setNewName] = useState("");

  const handleDuplicate = () => {
    onDuplicate(newNamespace, newName);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Duplicate Role</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="new-namespace" className="text-right">
              New Namespace
            </Label>
            <Input
              id="new-namespace"
              value={newNamespace}
              onChange={(e) => setNewNamespace(e.target.value)}
              className="col-span-3"
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="new-name" className="text-right">
              New Name
            </Label>
            <Input
              id="new-name"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              className="col-span-3"
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={handleDuplicate}>Duplicate</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default DuplicateRoleDialog;
