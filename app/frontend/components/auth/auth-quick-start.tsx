import { testUsers } from "@/lib/user";
import { User2Icon } from "lucide-react";
import { QuickStartLoginButton } from "./quick-start-login-button";

export const AuthQuickStart = () => {
	return (
		<section className="h-full w-full flex flex-col justify-center lg:max-w-sm lg:mx-auto">
			<h3 className="flex items-center gap-2 leading-4 tracking-tight font-light text-xl mb-4">
				<User2Icon />
				<span>Test the app with the following accounts</span>
			</h3>
			<ul className="mt-4 list-none space-y-2">
				{testUsers.map((acc) => (
					<li className="flex flex-col" key={acc.name}>
						<QuickStartLoginButton name={acc.name} />
						<p className="text-muted-foreground">{acc.explanation}</p>
					</li>
				))}
			</ul>
		</section>
	);
};
