import { z } from "zod";
import { BaseUser } from "..";
import { fetcher } from "./fetcher";

const ResponseSchema = z.union([
	z.object({
		ok: z.literal(true),
	}),
	z.object({
		error: z.string(),
	}),
]);

export const verifyEmail = async (user: BaseUser, id: string, code: string) => {
	const data = await fetcher("verify-email", "POST", user, { id, code });
	return ResponseSchema.safeParse(data);
};
