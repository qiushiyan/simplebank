"use server";
import { OneAccountResponseSchema } from "../account";
import { getCurrentUser } from "../auth";
import { fetcher } from "../fetcher";

export const updateAccountName = async (accountId: number, name: string) => {
	const user = await getCurrentUser();
	if (!user) {
		return null;
	}
	const data = await fetcher(`accounts/${accountId}`, "PATCH", user, { name });
	if (!data) {
		return null;
	}
	const parsed = OneAccountResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}

	return parsed.data;
};
