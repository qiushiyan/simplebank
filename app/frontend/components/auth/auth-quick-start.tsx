import { testUsers } from "@/lib/user";
import { User2Icon } from "lucide-react";
import { Separator } from "../ui/separator";
import { QuickStartLoginButton } from "./quick-start-login-button";

export const AuthQuickStart = () => {
	return (
		<section className="h-full w-full flex flex-col justify-center lg:max-w-sm lg:mx-auto space-y-4">
			<div>
				<h3 className="flex items-center gap-2 leading-4 tracking-tight font-light text-xl mb-4">
					<User2Icon />
					<span>Test it with the following accounts</span>
				</h3>
				<ul className="mt-4 list-none space-y-2">
					{testUsers.map((acc) => (
						<li className="flex flex-col" key={acc.name}>
							<QuickStartLoginButton name={acc.name} />
							<p className="text-muted-foreground">{acc.explanation}</p>
						</li>
					))}
				</ul>
			</div>
			<Separator className="bg-slate-500/50" />
			<div className="space-y-2">
				<h3 className="font-semibold">Privacy Notice</h3>
				<p className="text-sm text-muted-foreground">
					Email is optional to create accounts. If you provide an email the app
					will be able to send you emails upon your request, no spam will be
					sent. If you don't want to receive emails, leave the email field empty
					or use one of the accounts above.
				</p>
				<h3 className="font-semibold">Terms of Service</h3>
				<p className="text-sm text-muted-foreground">
					This is a fake banking app with no terms of service. It's possible to
					lose your data.
				</p>
			</div>
		</section>
	);
};
