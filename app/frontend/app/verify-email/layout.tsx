import { Metadata } from "next";

export const metadata: Metadata = {
	title: "Verify email",
	description: "Verify email",
};

export default function ({ children }: { children: React.ReactNode }) {
	return (
		<section className="min-h-screen w-screen flex items-center justify-center">
			{children}
		</section>
	);
}
