import { Metadata } from "next";

export const metadata: Metadata = {
	title: "Get your account",
	description: "Register or login to your account",
};

export default function ({ children }: { children: React.ReactNode }) {
	return <section className="min-h-screen">{children}</section>;
}
