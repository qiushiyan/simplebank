import { Session } from "next-auth";
import { z } from "zod";
import { BaseUser } from "..";
import { fetcher } from "./fetcher";

const FriendshipSchema = z.object({
	id: z.number(),
	from_account_id: z.number(),
	to_account_id: z.number(),
	pending: z.boolean(),
	accepted: z.boolean(),
	created_at: z.string(),
});

const ResponseSchema = z.object({
	data: z.array(FriendshipSchema),
});

export const getReceived = async (user: BaseUser, accountId: number) => {
	const data = await fetcher(
		`/friend/list?to_account_id=${accountId}`,
		"GET",
		user,
	);
	const parsed = ResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};

export const getSent = async (user: BaseUser, accountId: number) => {
	const data = await fetcher(
		`/friend/list?from_account_id=${accountId}`,
		"GET",
		user,
	);
	const parsed = ResponseSchema.safeParse(data);
	if (!parsed.success) {
		return null;
	}
	return parsed.data;
};
