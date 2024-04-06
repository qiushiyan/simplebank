import { AuthQuickStart } from "@/components/auth/auth-quick-start";
import { AuthTabs } from "@/components/auth/auth-tabs";
import { LogoutButton } from "@/components/auth/logout-button";
import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth";
import { ping } from "@/lib/fetcher";
import { routes } from "@/lib/navigataion";
import { ChevronLeftIcon, CommandIcon } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { config } from "../../lib/config";

type Props = {
	searchParams?: unknown;
};

export default async function ({ searchParams }: Props) {
	const user = await getCurrentUser();
	if (!(await ping())) {
		return notFound();
	}
	const { tab } = routes.auth.$parseSearchParams(searchParams);

	return (
		<div className="w-screen h-screen grid flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0">
			<Link href="/" className={"absolute top-4 left-4 md:top-8 md:left-8"}>
				<Button variant="ghost">
					<ChevronLeftIcon />
					Home
				</Button>
			</Link>
			<div className="col-span-2 lg:p-8 lg:col-span-1">
				<div className="mx-auto flex w-full flex-col justify-center space-y-8 sm:w-[350px]">
					<div className="flex flex-col space-y-2 text-center">
						<div className="flex justify-center items-center">
							<Image alt="logo" src="/logo.jpg" width={64} height={64} />
						</div>
						<h1 className="text-2xl font-handwritten tracking-tight">
							{config.title}
						</h1>
					</div>
					{user ? (
						<section className="font-light space-y-2">
							<p>
								It seems you are already logged in as{" "}
								<span className="font-semibold">{user.name}</span>
							</p>
							<p>Do you mean to log out?</p>
							<LogoutButton className="w-full" />
						</section>
					) : (
						<AuthTabs tab={tab} />
					)}
				</div>
			</div>
			<div className="hidden h-full bg-gray-100 lg:col-span-1 lg:flex lg:items-center">
				<AuthQuickStart />
			</div>
		</div>
	);
}
