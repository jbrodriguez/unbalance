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
}
