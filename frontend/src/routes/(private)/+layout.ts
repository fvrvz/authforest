import { resolve } from "$app/paths";
import { oidcService } from "$lib/services/oidc.service";
import { authStore } from "$lib/state/auth.svelte";
import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "../$types";

export const load: LayoutLoad = async () => {
  const user = await oidcService.getUser();

  if (!user || user.expired) {
    throw redirect(302, resolve("/login"));
  }

  authStore.setUser(user);
};
