import { AccountSearch } from "@/components/account/account-search";
import { BreadCrumb } from "@/components/account/breadcrumb";
import { Friendship } from "@/components/friendship/friendship";
import { EditableText } from "@/components/ui/editable-text";
import { getAccount } from "@/lib/account";
import { updateAccountName } from "@/lib/actions/account";
import { getCurrentUser } from "@/lib/auth";
import { routes } from "@/lib/navigataion";
import { revalidatePath } from "next/cache";
import { notFound } from "next/navigation";
import { Suspense } from "react";

type Props = {
	params: { id: string };
	searchParams?: { [key: string]: string | string[] | undefined };
};

export default async function ({ params, searchParams }: Props) {
	const { search_account_owner } =
		routes.account.$parseSearchParams(searchParams);
	const user = await getCurrentUser();
	if (!user) {
		return null;
	}

	const result = await getAccount(Number(params.id), user);
	if (!result) {
		return notFound();
	}

	const account = result.data;

	return (
		<div className="space-y-4">
			<BreadCrumb id={account.id} name={account.name} />
			<AccountSearch owner={search_account_owner} user={user} />
			<EditableText
				initialValue={account.name}
				fieldName="account"
				inputLabel="account name"
				buttonLabel="edit account name"
				inputClassName=""
				buttonClassName="justify-start"
				formAction={async (val) => {
					"use server";

					await updateAccountName(account.id, val);
					revalidatePath(".");
				}}
			/>
			<Suspense fallback={<Friendship.Skeleton />}>
				<Friendship accountId={account.id} user={user} />
			</Suspense>
		</div>
	);
}
