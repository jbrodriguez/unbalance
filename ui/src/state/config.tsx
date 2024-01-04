import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Api } from '~/api';

interface ConfigStore {
  version: string;
  dryRun: boolean;
  notifyPlan: number;
  notifyTransfer: number;
  reservedAmount: number;
  reservedUnit: string;
  rsyncArgs: string[];
  verbosity: number;
  checkForUpdate: number;
  refreshRate: number;
  actions: {
    getConfig: () => Promise<void>;
    toggleDryRun: () => Promise<void>;
    setNotifyPlan: (value: number) => Promise<void>;
    setNotifyTransfer: (value: number) => Promise<void>;
    setReservedSpace: (amount: number, unit: string) => Promise<void>;
  };
}

export const useConfigStore = create<ConfigStore>()(
  immer((set) => ({
    version: '',
    dryRun: true,
    notifyPlan: 0,
    notifyTransfer: 0,
    reservedAmount: 512,
    reservedUnit: 'MB',
    rsyncArgs: [],
    verbosity: 0,
    checkForUpdate: 0,
    refreshRate: 0,
    actions: {
      getConfig: async () => {
        const config = await Api.getConfig();
        // console.log('useConfigStore.getConfig() ', config);
        set((state) => {
          state.version = config.version;
          state.dryRun = config.dryRun;
          state.notifyPlan = config.notifyPlan;
          state.notifyTransfer = config.notifyTransfer;
          state.reservedAmount = config.reservedAmount;
          state.reservedUnit = config.reservedUnit;
          state.rsyncArgs = config.rsyncArgs;
          state.verbosity = config.verbosity;
          state.checkForUpdate = config.checkForUpdate;
          state.refreshRate = config.refreshRate;
        });
      },
      toggleDryRun: async () => {
        set((state) => {
          state.dryRun = !state.dryRun;
        });
        await Api.toggleDryRun();
      },
      setNotifyPlan: async (value: number) => {
        set((state) => {
          state.notifyPlan = value;
        });
        await Api.setNotifyPlan(value);
      },
      setNotifyTransfer: async (value: number) => {
        set((state) => {
          state.notifyTransfer = value;
        });
        await Api.setNotifyTransfer(value);
      },
      setReservedSpace: async (amount: number, unit: string) => {
        set((state) => {
          state.reservedAmount = amount;
          state.reservedUnit = unit;
        });
        await Api.setReservedSpace(amount, unit);
      },
    },
  })),
);

// export const useConfigActions = useConfigStore.getState().actions;

export const useConfigActions = () => useConfigStore((state) => state.actions);

export const useConfigVersion = () => useConfigStore((state) => state.version);
export const useConfigDryRun = () => useConfigStore((state) => state.dryRun);
export const useConfigNotifyPlan = () =>
  useConfigStore((state) => state.notifyPlan);
export const useConfigNotifyTransfer = () =>
  useConfigStore((state) => state.notifyTransfer);
export const useConfigReserved = () =>
  useConfigStore((state) => ({
    amount: state.reservedAmount,
    unit: state.reservedUnit,
  }));
