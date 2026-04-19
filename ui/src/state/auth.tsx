import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Api } from '~/api';

interface AuthStore {
  loaded: boolean;
  enabled: boolean;
  configured: boolean;
  authenticated: boolean;
  username: string;
  error: string;
  actions: {
    load: () => Promise<void>;
    login: (username: string, password: string) => Promise<boolean>;
    setup: (username: string, password: string) => Promise<boolean>;
    logout: () => Promise<void>;
    clearError: () => void;
  };
}

export const useAuthStore = create<AuthStore>()(
  immer((set) => ({
    loaded: false,
    enabled: false,
    configured: false,
    authenticated: false,
    username: 'admin',
    error: '',
    actions: {
      load: async () => {
        try {
          const status = await Api.getAuthStatus();
          set((state) => {
            state.loaded = true;
            state.enabled = status.enabled;
            state.configured = status.configured;
            state.authenticated = status.authenticated;
            state.username = status.username || 'admin';
            state.error = '';
          });
        } catch (e) {
          set((state) => {
            state.loaded = true;
            state.error = 'Unable to load authentication state';
          });
        }
      },
      login: async (username: string, password: string) => {
        try {
          const status = await Api.login(username, password);
          set((state) => {
            state.enabled = status.enabled;
            state.configured = status.configured;
            state.authenticated = status.authenticated;
            state.username = status.username || username;
            state.error = '';
          });
          return true;
        } catch (e) {
          set((state) => {
            state.error =
              e instanceof Error ? e.message : 'Unable to log in';
          });
          return false;
        }
      },
      setup: async (username: string, password: string) => {
        try {
          const status = await Api.setup(username, password);
          set((state) => {
            state.enabled = status.enabled;
            state.configured = status.configured;
            state.authenticated = status.authenticated;
            state.username = status.username || username;
            state.error = '';
          });
          return true;
        } catch (e) {
          set((state) => {
            state.error =
              e instanceof Error ? e.message : 'Unable to save password';
          });
          return false;
        }
      },
      logout: async () => {
        await Api.logout();
        set((state) => {
          state.authenticated = false;
          state.error = '';
        });
      },
      clearError: () => {
        set((state) => {
          state.error = '';
        });
      },
    },
  })),
);

export const useAuthActions = () => useAuthStore((state) => state.actions);
export const useAuthLoaded = () => useAuthStore((state) => state.loaded);
export const useAuthEnabled = () => useAuthStore((state) => state.enabled);
export const useAuthConfigured = () =>
  useAuthStore((state) => state.configured);
export const useAuthenticated = () =>
  useAuthStore((state) => state.authenticated);
export const useAuthUsername = () => useAuthStore((state) => state.username);
export const useAuthError = () => useAuthStore((state) => state.error);
