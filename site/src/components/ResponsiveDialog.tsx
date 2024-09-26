import { motion, AnimatePresence, Variants } from "framer-motion";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer";
import { useMediaQuery } from "@/hooks/use-media-query";
import { Send, X } from "lucide-react";
import { useEffect } from "react";
import { cn } from "@/lib/utils";

const backdropVariants: Variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

const contentVariants: Variants = {
  hidden: { opacity: 0, scale: 0.95, y: 20 },
  visible: { opacity: 1, scale: 1, y: 0 },
};

const childVariants: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

export function ResponsiveDialog({
  children,
  isOpen,
  setIsOpen,
  title,
  description,
  isLoading = false,
  onSubmit,
  isCloseEnabled,
  className,
}: {
  children: React.ReactNode;
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  title: string;
  description?: string;
  isLoading?: boolean;
  onSubmit?: () => void;
  isCloseEnabled?: boolean;
  className?: string;
}) {
  const isDesktop = useMediaQuery("(min-width: 768px)");

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setIsOpen(false);
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [setIsOpen]);

  const content = (
    <>
      <motion.div variants={childVariants}>
        {isDesktop ? (
          <DialogHeader>
            <DialogTitle className="mb-4">{title}</DialogTitle>
            {description && (
              <DialogDescription>{description}</DialogDescription>
            )}
          </DialogHeader>
        ) : (
          <DrawerHeader>
            <DrawerTitle className="mb-4">{title}</DrawerTitle>
            {description && (
              <DialogDescription>{description}</DialogDescription>
            )}
          </DrawerHeader>
        )}
      </motion.div>
      <motion.div variants={childVariants} className="custom-scrollbar">
        {isLoading ? <SkeletonContent /> : children}
      </motion.div>
      <motion.div
        variants={childVariants}
        className="flex gap-4 items-center mt-4"
      >
        {onSubmit && (
          <Button onClick={onSubmit} variant="default">
            <Send className="h-4 w-4" />
            Submit
          </Button>
        )}
        {isCloseEnabled && (<Button onClick={() => setIsOpen(false)} variant="outline">
          <X className="h-4 w-4" />
          Close
        </Button>) }
        
      </motion.div>
    </>
  );

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            className="fixed inset-0 bg-black/30 backdrop-blur-sm"
            variants={backdropVariants}
            initial="hidden"
            animate="visible"
            exit="hidden"
          />
          {isDesktop ? (
            <Dialog open={isOpen} onOpenChange={setIsOpen}>
              <DialogContent
                className={cn("sm:max-w-[425px]", className)}
                aria-describedby={undefined}
              >
                <motion.div
                  variants={contentVariants}
                  initial="hidden"
                  animate="visible"
                  exit="hidden"
                >
                  {content}
                </motion.div>
              </DialogContent>
            </Dialog>
          ) : (
            <Drawer open={isOpen} onOpenChange={setIsOpen}>
              <DrawerContent>
                <motion.div
                  variants={contentVariants}
                  initial="hidden"
                  animate="visible"
                  exit="hidden"
                >
                  {content}
                </motion.div>
              </DrawerContent>
            </Drawer>
          )}
        </>
      )}
    </AnimatePresence>
  );
}

const SkeletonContent = () => (
  <div className="space-y-2">
    <div className="h-4 bg-gray-200 rounded animate-pulse"></div>
    <div className="h-4 bg-gray-200 rounded animate-pulse"></div>
    <div className="h-4 bg-gray-200 rounded animate-pulse"></div>
  </div>
);
