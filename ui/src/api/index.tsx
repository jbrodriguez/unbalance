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
        dryRun: false,
        notifyPlan: 0,
        notifyTransfer: 0,
        reservedAmount: BigInt(0),
        reservedUnit: 'GB',
        rsyncArgs: [],
        verbosity: 0,
        checkForUpdate: 0,
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
        // plan: null,
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
}
