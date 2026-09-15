import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Hello World Acceptance 6",
  description: "Persisted greeting",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
