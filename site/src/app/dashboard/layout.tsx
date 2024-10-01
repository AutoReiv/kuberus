"use client";

import {
  Menu,
  Moon,
  Sun,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetTrigger } from "@/components/ui/sheet";
import SideNav from "./_components/SideNav";
import { useTheme } from "next-themes";
import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Toggle } from "@/components/ui/toggle";

export const pageVariants = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 },
  transition: { type: "spring", stiffness: 100, damping: 20 },
};

const Layout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const { setTheme, theme } = useTheme();
  const [isSideNavExpanded, setIsSideNavExpanded] = useState(false);
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    if (pathname === "/dashboard") {
      router.push("/dashboard/roles");
    }
  }, [pathname, router]);

  return (
    <motion.div
      className="grid min-h-screen w-full md:grid-cols-[0_1fr]"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.3 }}
    >
      <SideNav onExpand={setIsSideNavExpanded} />
      <motion.div
        className="flex flex-col"
        animate={{
          marginLeft: isSideNavExpanded ? "280px" : "64px",
          width: isSideNavExpanded ? "calc(100% - 280px)" : "calc(100% - 64px)",
        }}
        transition={{ duration: 0.4, ease: "easeInOut" }}
      >
        <motion.header
          className="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6"
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.4 }}
        >
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon"
                className="shrink-0 md:hidden"
              >
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle navigation menu</span>
              </Button>
            </SheetTrigger>
          </Sheet>
          <div className="w-full flex-1">{/*  */}</div>
          <Toggle
            aria-label="Toggle theme"
            pressed={theme === "dark"}
            onPressedChange={(pressed) => setTheme(pressed ? "dark" : "light")}
          >
            <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          </Toggle>
        </motion.header>
        <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
          <div className="flex flex-1">{children}</div>
        </main>
      </motion.div>
    </motion.div>
  );
};
export default Layout;
