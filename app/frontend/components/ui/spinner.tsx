import { cn } from "@/lib/utils";
import { Loader } from "lucide-react";

type Props = {
	className?: string;
};

export const Spinner = ({ className }: Props) => {
	return <Loader className={cn("size-4 animate-spin", className)} />;
};
