import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

interface ScatterStore {
  source: string;
  folders: Array<string>;
  targets: Array<string>;
  actions: {
    setSource: (source: string) => void;
  };
}

export const useConfigStore = create<ScatterStore>()(
  immer((set) => ({
    source: '',
    folders: [],
    targets: [],

    actions: {
      setSource: (source: string) => {
        set((state) => {
          state.source = source;
        });
      },
    },
  })),
);

export const useScatterActions = () => useConfigStore((state) => state.actions);

export const useScatterSource = () => useConfigStore((state) => state.source);
