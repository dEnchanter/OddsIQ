import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/layout/Navbar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "OddsIQ - AI Sports Betting Analytics",
  description: "AI-powered sports betting prediction and analytics platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} antialiased min-h-screen bg-background`}>
        <Navbar />
        <main className="container mx-auto py-6 px-4">
          {children}
        </main>
      </body>
    </html>
  );
}
