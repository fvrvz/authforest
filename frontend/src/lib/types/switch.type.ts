import type { Icon } from "lucide-svelte";
import type { ClassValue } from "svelte/elements";

export type SwitchProp = {
  label?: string;
  isChecked: boolean;
  onChange?: VoidFunction;
  Icon?: typeof Icon;
  size?: "default" | "lg";
};

export const SizeMap: Record<
  NonNullable<SwitchProp["size"]>,
  { knob: ClassValue; iconWrapper: ClassValue; icon: ClassValue }
> = {
  default: {
    knob: "h-6 w-11 after:size-5",
    iconWrapper: "size-5",
    icon: "size-3",
  },
  lg: {
    knob: "h-8 w-[3.8rem] after:size-[1.7rem]",
    iconWrapper: "size-[1.7rem]",
    icon: "size-4",
  },
} as const;
