import { createNavigationConfig } from "next-safe-navigation";
import { z } from "zod";

export const { routes, useSafeParams, useSafeSearchParams } =
	createNavigationConfig((defineRoute) => ({
		auth: defineRoute("/auth", {
			search: z
				.object({
					tab: z
						.union([z.literal("signin"), z.literal("signup")])
						.default("signin"),
					error: z.string().optional(),
				})
				.default({ tab: "signin" }),
		}),
		account: defineRoute("/accounts/[id]", {
			params: z.object({
				id: z.string(),
			}),
			search: z.object({
				search_account_owner: z.string().optional(),
			}),
		}),
	}));
