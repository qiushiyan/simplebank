import { cn } from "@/lib/utils";
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
	title: "Simplebank",
	description: "A banking system",
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body className={cn("max-w-3xl mx-auto")}>{children}</body>
		</html>
	);
}
