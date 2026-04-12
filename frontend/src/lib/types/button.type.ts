import type { Snippet } from "svelte";
import type { HTMLButtonAttributes } from "svelte/elements";

export type ButtonProps = {
  children: Snippet;
} & HTMLButtonAttributes;
