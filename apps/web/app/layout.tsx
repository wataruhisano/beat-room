import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Beatroom",
  description: "Private beat packs for each artist",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
