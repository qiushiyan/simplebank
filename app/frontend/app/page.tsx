import { LogoutButton } from "@/components/auth/logout-button";
import { buttonVariants } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth";
import Link from "next/link";

export default async function Home() {
	const user = await getCurrentUser();

	return (
		<main className="max-w-4xl mx-auto ">
			{user ? (
				<div>
					<h1>Welcome, {user.name}!</h1>
					<LogoutButton />
					<p>token {user.access_token}</p>
				</div>
			) : (
				<Link href={"/auth"} className={buttonVariants()}>
					sign in{" "}
				</Link>
			)}
		</main>
	);
}
