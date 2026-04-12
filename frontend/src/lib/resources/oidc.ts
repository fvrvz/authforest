import api from "$lib/api/api";
import { promiseWrapper, type Result } from "$lib/api/promise-wrapper";
import type { ListModel } from "$lib/types/common..type";
import type {
  DashboardStats,
  OAuthClient,
  RegisterClientRequest,
  UpdateClientRequest,
} from "$lib/types/oidc.type";

export class OIDC {
  static async registerClient(body: RegisterClientRequest) {
    const payload = {
      ...body,
      redirect_uris: body.redirect_uris.filter(Boolean),
    };
    return promiseWrapper<OAuthClient>(api.post("/oauth2/register", payload));
  }

  static async listClients(): Promise<Result<OAuthClient[]>> {
    const [err, res] = await promiseWrapper<ListModel<OAuthClient>>(
      api.get("/oauth2/clients"),
    );
    if (err) return [err, null];
    return [null, res.data];
  }

  static async getClient(clientId: string) {
    return promiseWrapper<OAuthClient>(api.get(`/oauth2/clients/${clientId}`));
  }

  static async updateClient(clientId: string, body: UpdateClientRequest) {
    const payload: Record<string, unknown> = { ...body };
    if (body.redirect_uris) {
      payload.redirect_uris = body.redirect_uris.filter(Boolean);
    }
    return promiseWrapper<OAuthClient>(
      api.patch(`/oauth2/clients/${clientId}`, payload),
    );
  }

  static async deleteClient(clientId: string) {
    return promiseWrapper<void>(api.delete(`/oauth2/clients/${clientId}`));
  }

  static async getStats() {
    return promiseWrapper<DashboardStats>(api.get("/oauth2/stats"));
  }
}
