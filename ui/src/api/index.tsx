import { State, Op, Branch, AuthStatus } from '~/types';

export class Api {
  static host = `${document.location.protocol}//${document.location.host}/api`;
  static csrfToken = '';

  static setCSRFToken(token: string) {
    Api.csrfToken = token;
  }

  static authHeaders() {
    const headers: Record<string, string> = {};
    if (Api.csrfToken !== '') {
      headers['X-CSRF-Token'] = Api.csrfToken;
    }

    return headers;
  }

  static async getConfig() {
    try {
      const response = await fetch(`${Api.host}/config`);
      const config = await response.json();
      return config;
    } catch (e) {
      return {
        version: '0.0.1',
        dryRun: true,
        notifyPlan: 0,
        notifyTransfer: 0,
        reservedAmount: 1,
        reservedUnit: 'GB',
        rsyncArgs: [],
        verbosity: 0,
        refreshRate: 0,
        authEnabled: false,
        authUsername: 'admin',
      };
    }
  }

  static async getAuthStatus(): Promise<AuthStatus> {
    const response = await fetch(`${Api.host}/auth/status`, {
      credentials: 'same-origin',
    });

    return response.json();
  }

  static async login(username: string, password: string): Promise<AuthStatus> {
    const response = await fetch(`${Api.host}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin',
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      throw new Error('Invalid username or password');
    }

    return response.json();
  }

  static async setup(username: string, password: string): Promise<AuthStatus> {
    const response = await fetch(`${Api.host}/auth/setup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin',
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const payload = await response.json().catch(() => null);
      throw new Error(payload?.message || 'Unable to set admin password');
    }

    return response.json();
  }

  static async logout(): Promise<AuthStatus> {
    const response = await fetch(`${Api.host}/auth/logout`, {
      method: 'POST',
      headers: Api.authHeaders(),
      credentials: 'same-origin',
    });

    if (!response.ok) {
      throw new Error('Unable to log out');
    }

    return response.json();
  }

  static async getUnraid(): Promise<State> {
    try {
      const response = await fetch(`${Api.host}/state`);
      const unraid = await response.json();
      return unraid;
    } catch (e) {
      return {
        status: Op.Neutral,
        unraid: null,
        operation: null,
        history: null,
      };
    }
  }

  static async getTree(path: string, id: string): Promise<Branch> {
    const encodedPath = encodeURIComponent(path);
    const encodedId = encodeURIComponent(id);
    try {
      const url = `${Api.host}/tree/${encodedPath}?id=${encodedId}`;
      const response = await fetch(url);
      const branch = await response.json();
      return branch;
    } catch (e) {
      return {
        nodes: {},
        order: [],
      };
    }
  }

  static async locate(path: string): Promise<Array<string>> {
    const encodedPath = encodeURIComponent(path);
    try {
      const url = `${Api.host}/locate/${encodedPath}`;
      const response = await fetch(url);
      const location = await response.json();
      return location;
    } catch (e) {
      return [];
    }
  }

  static async getLog(): Promise<Array<string>> {
    try {
      const url = `${Api.host}/logs`;
      const response = await fetch(url);
      const logs = await response.json();
      return logs;
    } catch (e) {
      return [];
    }
  }

  static async toggleDryRun(): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
    };
    try {
      const url = `${Api.host}/config/dryRun`;
      await fetch(url, options);
    } catch (e) {
      console.log('toggleDryRun() error: ', e);
    }
  }

  static async setNotifyPlan(value: number): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify(value),
    };
    try {
      const url = `${Api.host}/config/notifyPlan`;
      await fetch(url, options);
    } catch (e) {
      console.log('notifyPlan() error: ', e);
    }
  }

  static async setNotifyTransfer(value: number): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify(value),
    };
    try {
      const url = `${Api.host}/config/notifyTransfer`;
      await fetch(url, options);
    } catch (e) {
      console.log('notifyTransfer() error: ', e);
    }
  }

  static async setReservedSpace(amount: number, unit: string): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify({ amount, unit }),
    };
    try {
      const url = `${Api.host}/config/reservedSpace`;
      await fetch(url, options);
    } catch (e) {
      console.log('reservedSpace() error: ', e);
    }
  }

  static async setRsyncArgs(flags: string[]): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify(flags),
    };
    const url = `${Api.host}/config/rsyncArgs`;
    const response = await fetch(url, options);
    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || response.statusText);
    }
  }

  static async setVerbosity(value: number): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify(value),
    };
    try {
      const url = `${Api.host}/config/verbosity`;
      await fetch(url, options);
    } catch (e) {
      console.log('verbosity() error: ', e);
    }
  }

  static async setRefreshRate(value: number): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...Api.authHeaders() },
      body: JSON.stringify(value),
    };
    try {
      const url = `${Api.host}/config/refreshRate`;
      await fetch(url, options);
    } catch (e) {
      console.log('refreshRate() error: ', e);
    }
  }
}
