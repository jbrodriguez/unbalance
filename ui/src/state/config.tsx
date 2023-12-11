import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Api } from '~/api';

interface ConfigStore {
  version: string;
  dryRun: boolean;
  notifyPlan: number;
  notifyTransfer: number;
  reservedAmount: bigint;
  reservedUnit: string;
  rsyncArgs: string[];
  verbosity: number;
  checkForUpdate: number;
  refreshRate: number;
  actions: {
    getConfig: () => Promise<void>;
  };
}

export const useConfigStore = create<ConfigStore>()(
  immer((set) => ({
    version: '',
    dryRun: false,
    notifyPlan: 0,
    notifyTransfer: 0,
    reservedAmount: BigInt(0),
    reservedUnit: 'GB',
    rsyncArgs: [],
    verbosity: 0,
    checkForUpdate: 0,
    refreshRate: 0,
    actions: {
      getConfig: async () => {
        const config = await Api.getConfig();
        console.log('useConfigStore.getConfig() ', config);
        set((state) => {
          state.version = config.version;
          state.dryRun = config.dryRun;
          state.notifyPlan = config.notifyPlan;
          state.notifyTransfer = config.notifyTransfer;
          state.reservedAmount = BigInt(config.reservedAmount);
          state.reservedUnit = config.reservedUnit;
          state.rsyncArgs = config.rsyncArgs;
          state.verbosity = config.verbosity;
          state.checkForUpdate = config.checkForUpdate;
          state.refreshRate = config.refreshRate;
        });
      },
    },
  })),
);

// export const useConfigActions = useConfigStore.getState().actions;

export const useConfigActions = () => useConfigStore((state) => state.actions);
