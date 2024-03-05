import { Session } from "next-auth";
import { z } from "zod";
import { env } from "./env.mjs";
import { fetcher } from "./fetcher";

const AccountSchema = z.object({
	id: z.number(),
	name: z.string(),
	owner: z.string(),
	balance: z.number(),
	currency: z.string(),
	created_at: z.string(),
});

const ResponseSchema = z.object({
	data: z.array(AccountSchema),
});

export const getAccounts = async (user: Session["user"]) => {
	const data = await fetcher("accounts", "GET", user);
	const parsed = ResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};
