"use client";

import { cn } from "@/lib/utils";
import { useOptimistic, useRef, useState } from "react";
import { flushSync } from "react-dom";
import { Button, ButtonVariant } from "./button";
import { Input } from "./input";

type UpdatableText = {
	value: string;
	pending: boolean;
};

export function EditableText({
	fieldName,
	initialValue,
	inputClassName = "",
	inputLabel = "",
	buttonClassName = "",
	buttonLabel = "",
	buttonVariant = "secondary",
	formAction,
}: {
	fieldName: string;
	initialValue: string;
	inputClassName?: string;
	inputLabel?: string;
	buttonClassName?: string;
	buttonLabel?: string;
	buttonVariant?: ButtonVariant;
	formAction: (val: string) => Promise<void>;
}) {
	const [edit, setEdit] = useState(false);
	const inputRef = useRef<HTMLInputElement>(null);
	const buttonRef = useRef<HTMLButtonElement>(null);
	const [text, updateText] = useOptimistic(
		{ value: initialValue, pending: false },
		(state, value: UpdatableText) => {
			return value;
		},
	);

	const escapeEdit = () => {
		flushSync(() => {
			setEdit(false);
		});
		buttonRef.current?.focus();
	};

	return edit ? (
		<form
			action={async (formData) => {
				escapeEdit();
				const newVal = String(formData.get(fieldName));

				if (newVal && newVal.trim() !== "") {
					updateText({ value: newVal, pending: true });
				}
				await formAction(newVal);
			}}
		>
			<Input
				required
				ref={inputRef}
				type="text"
				aria-label={inputLabel}
				name={fieldName}
				defaultValue={initialValue}
				className={inputClassName}
				onKeyDown={(event) => {
					if (event.key === "Escape") {
						escapeEdit();
					}
				}}
				onBlur={(event) => {
					escapeEdit();
				}}
			/>
		</form>
	) : (
		<>
			<Button
				variant={buttonVariant}
				aria-label={buttonLabel}
				type="button"
				ref={buttonRef}
				onClick={() => {
					flushSync(() => {
						setEdit(true);
					});
					inputRef.current?.select();
				}}
				className={cn(buttonClassName, {
					"opacity-50": text.pending,
				})}
			>
				{text.value || <span className="text-slate-400 italic">Edit</span>}
			</Button>
		</>
	);
}
