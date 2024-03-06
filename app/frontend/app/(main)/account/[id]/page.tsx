import { EditableText } from "@/components/ui/editable-text";
import { getAccount, getAccounts } from "@/lib/account";
import { updateAccountName } from "@/lib/actions/account";
import { getCurrentUser } from "@/lib/auth";
import { delay } from "@/lib/utils";
import { revalidatePath } from "next/cache";
import { notFound } from "next/navigation";

export default async function ({ params }: { params: { id: string } }) {
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
		<div>
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
