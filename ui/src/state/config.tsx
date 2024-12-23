import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { useShallow } from 'zustand/shallow';

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
  refreshRate: number;
  actions: {
    getConfig: () => Promise<void>;
    toggleDryRun: () => Promise<void>;
    setNotifyPlan: (value: number) => Promise<void>;
    setNotifyTransfer: (value: number) => Promise<void>;
    setReservedSpace: (amount: number, unit: string) => Promise<void>;
    setRsyncArgs: (flags: string[]) => Promise<void>;
    resetRsyncArgs: () => Promise<void>;
    setVerbosity: (value: number) => Promise<void>;
    setRefreshRate: (value: number) => Promise<void>;
  };
}

export const useConfigStore = create<ConfigStore>()(
  immer((set) => ({
    version: '',
    dryRun: true,
    notifyPlan: 0,
    notifyTransfer: 0,
    reservedAmount: 1,
    reservedUnit: 'Gb',
    rsyncArgs: ['-X'],
    verbosity: 0,
    refreshRate: 1000,
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
      setRsyncArgs: async (flags: string[]) => {
        set((state) => {
          state.rsyncArgs = flags;
        });
        await Api.setRsyncArgs(flags);
      },
      resetRsyncArgs: async () => {
        set((state) => {
          state.rsyncArgs = ['-X'];
        });
        await Api.setRsyncArgs(['-X']);
      },
      setVerbosity: async (value: number) => {
        set((state) => {
          state.verbosity = value;
        });
        await Api.setVerbosity(value);
      },
      setRefreshRate: async (value: number) => {
        set((state) => {
          state.refreshRate = value;
        });
        await Api.setRefreshRate(value);
      },
    },
  })),
);

export const useConfigActions = () => useConfigStore((state) => state.actions);

export const useConfigVersion = () => useConfigStore((state) => state.version);
export const useConfigDryRun = () => useConfigStore((state) => state.dryRun);
export const useConfigNotifyPlan = () =>
  useConfigStore((state) => state.notifyPlan);
export const useConfigNotifyTransfer = () =>
  useConfigStore((state) => state.notifyTransfer);
export const useConfigReserved = () =>
  useConfigStore(
    useShallow((state) => ({
      amount: state.reservedAmount,
      unit: state.reservedUnit,
    })),
  );
export const useConfigRsyncArgs = () =>
  useConfigStore((state) => state.rsyncArgs);
export const useConfigVerbosity = () =>
  useConfigStore((state) => state.verbosity);
export const useConfigRefreshRate = () =>
  useConfigStore((state) => state.refreshRate);
