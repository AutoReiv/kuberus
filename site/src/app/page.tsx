"use client";

import { motion } from "framer-motion";

export default function HomePage() {
  return (
    <motion.div className="min-h-screen bg-background flex flex-col items-center justify-center p-4 relative overflow-hidden">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md z-10"
      >
        -- INSERT DASHBOARD HERE --
      </motion.div>
    </motion.div>
  );
}
