"use client";

import { useRouter } from "next/navigation";

type Props = {
	owner: string | undefined;
};

export const AccountSearch = ({ owner }: Props) => {
	const router = useRouter();
	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
				const formData = new FormData(e.target as HTMLFormElement);
				const owner = formData.get("search") as string;
				const url = new URL(window.location.href);
				url.searchParams.set("search_account_owner", owner);
				router.push(url.toString());
			}}
		>
			<input
				type="search"
				name="search"
				className="p-4"
				required
				placeholder="search accounts by owner"
				defaultValue={owner || ""}
			/>
		</form>
	);
};
