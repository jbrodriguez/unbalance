import { State, Op, Branch } from '~/types';

export class Api {
  static host = `${document.location.protocol}//${document.location.host}/api`;

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
      };
    }
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
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(flags),
    };
    try {
      const url = `${Api.host}/config/rsyncArgs`;
      await fetch(url, options);
    } catch (e) {
      console.log('rsyncArgs() error: ', e);
    }
  }

  static async setVerbosity(value: number): Promise<void> {
    const options = {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
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
