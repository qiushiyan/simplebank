import { cn } from "@/lib/utils";
import type { Metadata } from "next";
import {
	Inter as FontSans,
	Luckiest_Guy as FontHandwritten,
} from "next/font/google";

import { config } from "@/lib/config";
import "./globals.css";

const fontSans = FontSans({
	subsets: ["latin"],
	variable: "--font-sans",
});

const fontHandwritten = FontHandwritten({
	subsets: ["latin"],
	weight: "400",
	variable: "--font-handwritten",
});

export const metadata: Metadata = {
	title: config.title,
	description: config.description,
};

export default async function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body
				className={cn(
					"min-h-screen bg-background font-sans antialiased",
					fontSans.variable,
					fontHandwritten.variable,
				)}
			>
				{children}
			</body>
		</html>
	);
}
