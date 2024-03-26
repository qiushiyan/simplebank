"use client";

import { useRouter } from "next/navigation";
import { useOptimistic, useTransition } from "react";

type Props = {
	owner: string | undefined;
};

export const AccountSearchInput = ({ owner }: Props) => {
	const [optimisticOwner, setOptimisticOwner] = useOptimistic(owner);
	const router = useRouter();
	const [pending, startTransition] = useTransition();
	return (
		<form
			data-pending={pending ? "" : undefined}
			onSubmit={(e) => {
				e.preventDefault();
				const formData = new FormData(e.target as HTMLFormElement);
				const owner = formData.get("search") as string;
				const url = new URL(window.location.href);
				url.searchParams.set("search_account_owner", owner);

				startTransition(() => {
					setOptimisticOwner(owner);
					router.push(url.toString());
				});
			}}
		>
			<input
				type="search"
				name="search"
				className="p-4"
				required
				placeholder="search accounts by owner"
				defaultValue={optimisticOwner}
			/>
		</form>
	);
};
