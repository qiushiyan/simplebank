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
				})
				.default({ tab: "signin" }),
		}),
		invoice: defineRoute("/invoices/[invoiceId]", {
			params: z.object({
				invoiceId: z.string(),
			}),
		}),
	}));
