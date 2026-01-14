import type { Metadata } from "next";
import "./globals.css";

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
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
