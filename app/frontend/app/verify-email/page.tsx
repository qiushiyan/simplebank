import { getCurrentUser } from "@/lib/auth";
import { fetcher } from "@/lib/fetcher";
import { routes } from "@/lib/navigataion";
import { delay } from "@/lib/utils";
import { verifyEmail } from "@/lib/verify-email";
import { notFound } from "next/navigation";

type Props = {
	searchParams?: unknown;
};

export default async function ({ searchParams }: Props) {
	const { code, id } = routes.verifyEmail.$parseSearchParams(searchParams);
	const user = await getCurrentUser();
	await delay(1000);
	if (!user) {
		return notFound();
	}

	const response = await verifyEmail(user, id, code);
	if (!response.success) {
		return notFound();
	}

	if ("error" in response.data) {
		return (
			<div className="text-center">
				<p className="text-base font-semibold">404</p>
				<h1 className="mt-4 text-3xl font-bold tracking-tight text-primary/80 sm:text-5xl">
					{response.data.error}
				</h1>
				<div className="mt-10 flex items-center justify-center gap-x-6">
					<a href="/" className="text-sm font-semibold ">
						Go to home <span aria-hidden="true">&rarr;</span>
					</a>
				</div>
			</div>
		);
	}

	return (
		<div className="text-center">
			<p className="text-base font-semibold">404</p>
			<h1 className="mt-4 text-3xl font-bold tracking-tight text-primary/80 sm:text-5xl">
				Email verified
			</h1>
			<div className="mt-10 flex items-center justify-center gap-x-6">
				<a href="/" className="text-sm font-semibold ">
					Go to home <span aria-hidden="true">&rarr;</span>
				</a>
			</div>
		</div>
	);
}
