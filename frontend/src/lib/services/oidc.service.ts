import { browser } from "$app/environment";
import { resolve } from "$app/paths";
import { UserManager, WebStorageStateStore } from "oidc-client-ts";

let userManager: UserManager | null = null;

function getUserManager(): UserManager {
  if (!userManager) {
    userManager = new UserManager({
      authority: import.meta.env.VITE_BASE_URL as string,
      client_id: import.meta.env.VITE_OIDC_CLIENT_ID as string,
      redirect_uri: `${window.location.origin}${resolve("/callback")}`,
      response_type: "code",
      scope: "openid profile email",
      automaticSilentRenew: false,
      userStore: new WebStorageStateStore({ store: window.sessionStorage }),
    });
  }
  return userManager;
}

export const oidcService = {
  async getUser() {
    if (!browser) return null;
    return getUserManager().getUser();
  },

  async getAccessToken() {
    if (!browser) return null;
    const user = await getUserManager().getUser();
    if (!user || user.expired) return null;
    return user.access_token;
  },

  async login() {
    return getUserManager().signinRedirect();
  },

  async handleCallback() {
    return getUserManager().signinRedirectCallback();
  },

  async logout() {
    if (!browser) return;
    await getUserManager().removeUser();
  },
};
