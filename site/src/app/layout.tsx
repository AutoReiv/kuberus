"use client";
import { Poppins } from "next/font/google";
import "./globals.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/components/ui/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import { cn } from "@/lib/utils";

const poppins = Poppins({
  weight: ["400", "500", "600", "700"],
  subsets: ["latin"],
});

const queryClient = new QueryClient();

export default function RootLayout({
  children,
}) {
  return (
    <QueryClientProvider client={queryClient}>
      <html lang="en" suppressHydrationWarning className="max-h-screen">
        <body className={cn(poppins.className, "max-h-screen")}>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            {children}
            <Toaster />
          </ThemeProvider>
        </body>
      </html>
    </QueryClientProvider>
  );
}
