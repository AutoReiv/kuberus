"use client";
import Link from "next/link";

import { Bell, Blend, Package2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { usePathname } from "next/navigation";

const SideNav = () => {
  const path = usePathname();

  const MenuList = [
    {
      name: "Roles",
      icon: "",
      path: "/dashboard/roles",
    },
    {
      name: "Role Bindings",
      icon: "",
      path: "/dashboard/role-bindings",
    },
    {
      name: "Cluster Roles",
      icon: "",
      path: "/dashboard/cluster-roles",
    },
    {
      name: "Cluster Role Bindings",
      icon: "",
      path: "/dashboard/cluster-role-bindings",
    },
    {
      name: "Service Accounts",
      icon: "",
      path: "/dashboard/service-accounts",
    },
    {
      name: "Groups",
      icon: "",
      path: "/dashboard/groups",
    },
  ];
  return (
    <div className="hidden border-r bg-muted/40 md:block max-h-screen">
      <div className="flex h-full max-h-screen flex-col gap-2">
        <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
          <Link href="/" className="flex items-center gap-2 font-semibold">
            <Package2 className="h-6 w-6" />
            <span className="">Phasing.</span>
          </Link>
          <Button variant="outline" size="icon" className="ml-auto h-8 w-8">
            <Bell className="h-4 w-4" />
            <span className="sr-only">Toggle notifications</span>
          </Button>
        </div>
        <div className="flex-1">
          <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
            {MenuList.map((menu, index) => (
              <Link
                key={index}
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary ${
                  path === menu.path &&
                  "bg-accent-foreground text-secondary hover:text-secondary"
                }`}
                href={menu.path}
              >
                {/* <menu.icon className="h-6 w-6" /> */}
                <h6 className="text-sm">{menu.name}</h6>
              </Link>
            ))}
          </nav>
        </div>
        <div className="mt-auto p-4">
          <Card>
            <CardHeader>
              <CardTitle>
                <div className="flex items-center gap-2">
                  <Blend className="h-4 w-4" />
                  <span className="text-sm font-medium">Suhhhhh</span>
                </div>
              </CardTitle>
            </CardHeader>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default SideNav;
