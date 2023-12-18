import { State, Op, Branch } from '~/types';

export class Api {
  static host = `${document.location.protocol}//${document.location.host}/api`;

  static async getConfig() {
    console.log('Api.getConfig() ', Api.host);
    try {
      const response = await fetch(`${Api.host}/config`);
      const config = await response.json();
      return config;
    } catch (e) {
      console.log('Api.getConfig() error ', e);
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
    console.log('Api.getUnraid() ', Api.host);
    try {
      const response = await fetch(`${Api.host}/state`);
      const unraid = await response.json();
      return unraid;
    } catch (e) {
      console.log('Api.getUnraid() error ', e);
      return {
        status: Op.Neutral,
        unraid: null,
        operation: null,
        history: null,
        plan: null,
      };
    }
  }

  static async getTree(path: string, id: string): Promise<Branch> {
    const encodedPath = encodeURIComponent(path);
    const encodedId = encodeURIComponent(id);
    // console.log('Api.getTree() ', Api.host, path, encodedPath);
    try {
      const url = `${Api.host}/tree/${encodedPath}?id=${encodedId}`;
      console.log('Api.getTree() url ', url);
      const response = await fetch(url);
      const branch = await response.json();
      return branch;
    } catch (e) {
      console.log('Api.getTree() error ', e);
      return {
        nodes: {},
        order: [],
      };
    }
  }
}
