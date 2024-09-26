"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { motion, useAnimation } from "framer-motion";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import {
  ShieldCheck,
  UserPlus,
  Info,
  Mail,
  Play,
  HelpCircle,
} from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

const featureVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

const features = [
  { icon: ShieldCheck, title: "Secure Access Control" },
  { icon: UserPlus, title: "Easy User Management" },
  { icon: Info, title: "Detailed Audit Logs" },
];

export default function HomePage() {
  return (
    <motion.div className="min-h-screen bg-background flex flex-col items-center justify-center p-4 relative overflow-hidden">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md z-10"
      >
        <motion.div
          whileHover={{
            scale: 1.05,
            boxShadow: "0px 0px 8px rgb(255,255,255)",
          }}
          className="bg-card text-card-foreground"
        >
          <Card className="border-none shadow-xl">
            <CardHeader className="text-center">
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", stiffness: 260, damping: 20 }}
              >
                <ShieldCheck className="w-16 h-16 mx-auto text-primary mb-4" />
              </motion.div>
              <CardTitle className="text-3xl font-bold">RBAC Manager</CardTitle>
              <CardDescription className="text-lg">
                Secure. Efficient. Role-Based Access Control.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Link href="/login">
                  <Button className="w-full text-lg py-6">
                    <motion.div transition={{ duration: 0.3 }}>
                      <ShieldCheck className="mr-2 h-6 w-6" />
                    </motion.div>
                    Login
                  </Button>
                </Link>
              </motion.div>
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Link href="/register">
                  <Button variant="outline" className="w-full text-lg py-6">
                    <motion.div transition={{ duration: 0.3 }}>
                      <UserPlus className="mr-2 h-6 w-6" />
                    </motion.div>
                    Register
                  </Button>
                </Link>
              </motion.div>
            </CardContent>
          </Card>
        </motion.div>
        <div className="mt-8 flex justify-center space-x-4">
          <Link href="/about">
            <Button variant="ghost" size="sm">
              <Info className="mr-2 h-4 w-4" />
              About
            </Button>
          </Link>
          <Link href="/contact">
            <Button variant="ghost" size="sm">
              <Mail className="mr-2 h-4 w-4" />
              Contact
            </Button>
          </Link>
          {/* <ModeToggle /> */}
          <Dialog>
            <DialogTrigger asChild>
              <Button variant="outline" size="sm">
                <Play className="mr-2 h-4 w-4" />
                Quick Demo
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>RBAC Manager Demo</DialogTitle>
              </DialogHeader>
              <div className="aspect-video">Something Here...</div>
            </DialogContent>
          </Dialog>
        </div>
      </motion.div>
      <motion.div
        initial="hidden"
        animate="visible"
        transition={{ staggerChildren: 0.1 }}
        className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6 z-10"
      >
        {features.map((feature, index) => (
          <motion.div
            key={index}
            variants={featureVariants}
            className="text-center"
          >
            <feature.icon className="w-12 h-12 mx-auto text-primary mb-2" />
            <h3 className="font-semibold">{feature.title}</h3>
          </motion.div>
        ))}
      </motion.div>
      <motion.div
        className="fixed bottom-4 right-4 z-20"
        whileHover={{ scale: 1.1 }}
        whileTap={{ scale: 0.9 }}
      >
        <Button size="lg" className="rounded-full">
          <HelpCircle className="h-6 w-6" />
        </Button>
      </motion.div>
    </motion.div>
  );
}
