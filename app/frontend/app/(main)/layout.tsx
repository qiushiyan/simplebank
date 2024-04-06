import { AccountLink } from "@/components/account/account-link";
import { Button, buttonVariants } from "@/components/ui/button";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { getAccounts } from "@/lib/account";
import { getCurrentUser } from "@/lib/auth";
import { env } from "@/lib/env.mjs";
import { ping } from "@/lib/fetcher";
import { ChevronsUpDown } from "lucide-react";
import Link from "next/link";
import { notFound } from "next/navigation";

export default async function ({ children }: { children: React.ReactNode }) {
	const user = await getCurrentUser();
	if (!user) {
		return (
			<section className="max-w-[960px] mx-auto flex flex-row gap-4 p-4">
				<Link
					href="/auth"
					className={buttonVariants({
						variant: "link",
					})}
				>
					log in
				</Link>
			</section>
		);
	}

	if (!(await ping())) {
		return notFound();
	}
	const accounts = await getAccounts(user);

	return (
		<section className="max-w-4xl mx-auto">
			<header className="text-center my-4">
				<h2 className="text-3xl font-handwritten text-primary hover:underline">
					<Link href="/">Simple Bank</Link>
				</h2>
			</header>
			<section className="flex flex-row gap-4 group">
				<aside className="w-[200px]">
					<nav>
						<ul className="space-y-4">
							<li>
								<Collapsible defaultOpen={true} className="space-y-2">
									<div className="flex items-center justify-between space-x-4">
										<h4 className="text-sm font-semibold">Accounts</h4>
										<CollapsibleTrigger asChild>
											<Button variant="ghost" size="sm" className="w-9 p-0">
												<ChevronsUpDown className="h-4 w-4" />
												<span className="sr-only">Toggle</span>
											</Button>
										</CollapsibleTrigger>
									</div>
									<CollapsibleContent className="flex flex-col space-y-2">
										{accounts ? (
											accounts.data.map((acc) => (
												<AccountLink id={acc.id} name={acc.name} key={acc.id} />
											))
										) : (
											<p>error loading accounts</p>
										)}
									</CollapsibleContent>
								</Collapsible>
							</li>
							<li>
								<Link href="/auth">Auth</Link>
							</li>
						</ul>
					</nav>
				</aside>
				<section className="flex-1 group-has-[[data-pending]]:animate-pulse">
					{children}
				</section>
			</section>
		</section>
	);
}
