"use client";
import Link from "next/link";
import {
  UserCog,
  UserCheck,
  Shield,
  ShieldCheck,
  UserCircle,
  Users,
  Package2,
  Bell,
  Blend,
  Group,
  GroupIcon,
  Folder,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

const sidebarVariants = {
  expanded: { width: 280 },
  collapsed: { width: 64 },
};

const menuItemVariants = {
  expanded: { opacity: 1, x: 0 },
  collapsed: { opacity: 0, x: -20 },
};

const SideNav = ({ onExpand }: { onExpand: (expanded: boolean) => void }) => {
  const path = usePathname();
  const [isExpanded, setIsExpanded] = useState(false);

  const handleExpand = (expanded: boolean) => {
    setIsExpanded(expanded);
    onExpand(expanded);
  };

  const MenuList = [
    {
      name: "Roles",
      icon: <UserCog className="h-6 w-6" />,
      path: "/dashboard/roles",
    },
    {
      name: "Role Bindings",
      icon: <UserCheck className="h-6 w-6" />,
      path: "/dashboard/role-bindings",
    },
    {
      name: "Cluster Roles",
      icon: <Shield className="h-6 w-6" />,
      path: "/dashboard/cluster-roles",
    },
    {
      name: "Cluster Role Bindings",
      icon: <ShieldCheck className="h-6 w-6" />,
      path: "/dashboard/cluster-role-bindings",
    },
    {
      name: "Namespaces",
      icon: <Folder className="h-6 w-6" />,
      path: "/dashboard/namespaces",
    },
    {
      name: "Service Accounts",
      icon: <UserCircle className="h-6 w-6" />,
      path: "/dashboard/service-accounts",
    },
    {
      name: "Groups",
      icon: <Group className="h-6 w-6" />,
      path: "/dashboard/groups",
    },
    {
      name: "Users",
      icon: <Users className="h-6 w-6" />,
      path: "/dashboard/users",
    },
  ];

  return (
    <motion.div
      className="hidden border-r bg-muted/40 md:block max-h-screen"
      initial="collapsed"
      animate={isExpanded ? "expanded" : "collapsed"}
      variants={sidebarVariants}
      transition={{ duration: 0.3, type: "spring", stiffness: 100 }}
      onMouseEnter={() => handleExpand(true)}
      onMouseLeave={() => handleExpand(false)}
    >
      <div className="flex h-full max-h-screen flex-col gap-2">
        <motion.div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
          <Link href="/" className="flex items-center gap-2 font-semibold">
            <Package2 className="h-6 w-6" />
            <AnimatePresence>
              {isExpanded && (
                <motion.span
                  initial={{ opacity: 0, width: 0 }}
                  animate={{ opacity: 1, width: "auto" }}
                  exit={{ opacity: 0, width: 0 }}
                  transition={{ duration: 0.2 }}
                >
                  Phasing.
                </motion.span>
              )}
            </AnimatePresence>
          </Link>
          <AnimatePresence>
            {isExpanded && (
              <motion.div
                initial={{ opacity: 0, scale: 0 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0 }}
                transition={{ duration: 0.2 }}
              >
                <Button
                  variant="outline"
                  size="icon"
                  className="ml-auto h-8 w-8"
                >
                  <Bell className="h-4 w-4" />
                  <span className="sr-only">Toggle notifications</span>
                </Button>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
        <div className="flex-1">
          <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
            {MenuList.map((menu, index) => (
              <motion.div
                key={index}
                variants={menuItemVariants}
                initial="initial"
                animate="animate"
                exit="exit"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Link
                  className={`flex items-center gap-3 rounded-lg p-2 text-muted-foreground transition-all ${
                    path.startsWith(menu.path)
                      ? "bg-accent-foreground text-secondary"
                      : isExpanded
                      ? "hover:text-primary hover:bg-accent/50"
                      : ""
                  }`}
                  href={menu.path}
                >
                  <motion.div
                    whileHover={{ rotate: 10, scale: 1.1 }}
                    transition={{ type: "spring", stiffness: 300 }}
                  >
                    {menu.icon}
                  </motion.div>
                  <AnimatePresence>
                    {isExpanded && (
                      <motion.h6
                        className="text-sm"
                        initial={{ opacity: 0, width: 0 }}
                        animate={{ opacity: 1, width: "auto" }}
                        exit={{ opacity: 0, width: 0 }}
                        transition={{ duration: 0.2 }}
                      >
                        {menu.name}
                      </motion.h6>
                    )}
                  </AnimatePresence>
                </Link>
              </motion.div>
            ))}
          </nav>
        </div>
        <motion.div
          className="mt-auto p-4"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
        >
          <Card>
            <CardHeader>
              <CardTitle>
                <div className="flex items-center gap-2">
                  <AnimatePresence>
                    {isExpanded && (
                      <motion.span
                        className="text-sm font-medium flex items-center gap-4"
                        initial={{ opacity: 0, width: 0 }}
                        animate={{ opacity: 1, width: "auto" }}
                        exit={{ opacity: 0, width: 0 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Link
                          href="/dashboard/admin"
                          className="flex items-center gap-4"
                        >
                          <Blend className="h-4 w-4" />
                          Suhhhhh
                        </Link>
                      </motion.span>
                    )}
                  </AnimatePresence>
                </div>
              </CardTitle>
            </CardHeader>
          </Card>
        </motion.div>
      </div>
    </motion.div>
  );
};

export default SideNav;
