import { Session } from "next-auth";
import { z } from "zod";
import { fetcher } from "./fetcher";

export const AccountSchema = z.object({
	id: z.number(),
	name: z.string(),
	owner: z.string(),
	balance: z.number(),
	currency: z.string(),
	created_at: z.string(),
});

const AccountsResponseSchema = z.object({
	data: z.array(AccountSchema),
});

export const OneAccountResponseSchema = z.object({
	data: AccountSchema,
});

export const getAccounts = async (user: Session["user"]) => {
	const data = await fetcher("accounts", "GET", user);
	const parsed = AccountsResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};

export const getAccount = async (id: number, user: Session["user"]) => {
	const data = await fetcher(`accounts/${id}`, "GET", user);
	const parsed = OneAccountResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};

export const searchAccounts = async (owner: string, user: Session["user"]) => {
	const data = await fetcher(`accounts/search?owner=${owner}`, "GET", user);
	const parsed = AccountsResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};
