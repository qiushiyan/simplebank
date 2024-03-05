import { AuthQuickStart } from "@/components/auth/auth-quick-start";
import { AuthTabs } from "@/components/auth/auth-tabs";
import { Button } from "@/components/ui/button";
import { routes } from "@/lib/navigataion";
import { ChevronLeftIcon, CommandIcon } from "lucide-react";
import Link from "next/link";
import { config } from "../../lib/config";

type Props = {
	searchParams?: unknown;
};

export default function ({ searchParams }: Props) {
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
						<CommandIcon className="mx-auto size-6" />
						<h1 className="text-2xl font-semibold tracking-tight">
							{config.title}
						</h1>
					</div>
					<AuthTabs tab={tab} />
				</div>
			</div>
			<div className="hidden h-full bg-gray-100 lg:col-span-1 lg:flex lg:items-center">
				<AuthQuickStart />
			</div>
		</div>
	);
}
