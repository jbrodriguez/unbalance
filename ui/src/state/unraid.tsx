import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Api } from '~/api';
import { Unraid, Operation, History, Plan, Op, Step } from '~/types';
import { convertStatusToStep } from '~/helpers/steps';

interface UnraidStore {
  loaded: boolean;
  // route: string;
  status: Op;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
  step: Step;
  actions: {
    getUnraid: () => Promise<void>;
  };
}

export const useUnraidStore = create<UnraidStore>()(
  immer((set) => {
    const protocol =
      document.location.protocol === 'https:' ? 'wss://' : 'ws://';

    const socket = new WebSocket(`${protocol}${document.location.host}/ws`);

    socket.onopen = function (event) {
      console.log('Socket opened ', event);
    };

    socket.onmessage = function (event) {
      console.log('Socket message ', event);
    };

    socket.onclose = function (event) {
      console.log('Socket closed ', event);
    };

    return {
      loaded: false,
      // route: '/scatter',
      status: Op.Neutral,
      unraid: null,
      operation: null,
      history: null,
      plan: null,
      step: 'idle',
      actions: {
        getUnraid: async () => {
          const array = await Api.getUnraid();
          console.log('useUnraidStore.getUnraid() ', array);
          set((state) => {
            state.loaded = true;
            state.status = array.status;
            state.unraid = array.unraid;
            state.operation = array.operation;
            state.history = array.history;
            state.plan = array.plan;
            state.step = convertStatusToStep(array.status);
          });
        },
        setCurrentStep: (step: Step) => {
          set((state) => {
            state.step = step;
          });
        },
      },
    };
  }),
);

export const useUnraidActions = () => useUnraidStore().actions;

export const useUnraidLoaded = () => useUnraidStore().loaded;
export const useUnraidStatus = () => useUnraidStore().status;
export const useUnraidStep = () => useUnraidStore().step;
export const useUnraidIsBusy = () => useUnraidStore().status !== Op.Neutral;
