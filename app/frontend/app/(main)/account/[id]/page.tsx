import { AccountSearch } from "@/components/account/account-search";
import { AccountSearchList } from "@/components/account/account-search-list";
import { BreadCrumb } from "@/components/account/breadcrumb";
import { EditableText } from "@/components/ui/editable-text";
import { Spinner } from "@/components/ui/spinner";
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
			<AccountSearch owner={search_account_owner} />
			<Suspense fallback={<Spinner />}>
				<AccountSearchList user={user} owner={search_account_owner} />
			</Suspense>
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
		</div>
	);
}
