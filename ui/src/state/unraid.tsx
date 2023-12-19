import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { NavigateFunction } from 'react-router-dom';

import { Api } from '~/api';
import {
  Unraid,
  Operation,
  History,
  Plan,
  Op,
  Step,
  Packet,
  Topic,
} from '~/types';
// import { convertStatusToStep } from '~/helpers/steps';
import { useScatterStore } from '~/state/scatter';
// import {
//   CommandScatterPlanStart,
//   EventScatterPlanStarted,
//   EventScatterPlanProgress,
//   EventScatterPlanEnded,
// } from '~/constants';

interface UnraidStore {
  socket: WebSocket;
  loaded: boolean;
  route: string;
  status: Op;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
  step: Step;
  actions: {
    getUnraid: () => Promise<void>;
    // setCurrentStep: (step: Step) => void;
    syncRouteAndStep: (path: string) => void;
    transition: (navigate: NavigateFunction) => void;
    scatterProgress: (payload: string) => void;
    scatterPlanEnded: (payload: string) => void;
  };
}

const getNextStep = (step: Step) => {
  switch (step) {
    case 'select':
      return 'plan';
    case 'plan':
      return 'transfer';
    case 'transfer':
      return 'idle';
    default:
      return 'idle';
  }
};

const mapEventToAction: { [x: string]: string } = {
  [Topic.EventScatterPlanStarted]: 'scatterProgress',
  [Topic.EventScatterPlanProgress]: 'scatterProgress',
  [Topic.EventScatterPlanEnded]: 'scatterPlanEnded',
  // 'scatter:plan:started': 'scatterProgress',
  // 'scatter:plan:progress': 'scatterProgress',
  // 'scatter:plan:ended': 'scatterProgress',
  // [`${EventScatterPlanProgress}`]: 'scatterProgress',
  // [`${EventScatterPlanEnded}`]: 'scatterProgress',
};

export const useUnraidStore = create<UnraidStore>()(
  immer((set, get) => {
    const protocol =
      document.location.protocol === 'https:' ? 'wss://' : 'ws://';

    const socket = new WebSocket(`${protocol}${document.location.host}/ws`);

    socket.onopen = function (event) {
      console.log('Socket opened ', event);
    };

    socket.onmessage = function (event) {
      // console.log('Socket message ', event);
      const packet: Packet = JSON.parse(event.data);
      const action = mapEventToAction[packet.topic];
      if (!action) {
        return;
      }
      // @ts-expect-error -- TSCONVERSION
      get().actions[action](packet.payload);
    };

    socket.onclose = function (event) {
      console.log('Socket closed ', event);
    };

    return {
      socket,
      loaded: false,
      route: '/',
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
            // state.step = convertStatusToStep(array.status);
          });
        },
        // setCurrentStep: (step: Step) => {
        //   set((state) => {
        //     state.step = step;
        //   });
        // },
        syncRouteAndStep: (path: string) => {
          const route = get().route;
          // if (route.slice(0, 5) === path.slice(0, 5)) {
          //   return;
          // }
          if (route === path) {
            return;
          }

          console.log('syncStep ', route, path);

          set((state) => {
            state.route = path;
            // state.step = 'select';
          });
        },
        transition: (navigate: NavigateFunction) => {
          const next = getNextStep(get().step);
          set((state) => {
            state.status = Op.ScatterPlanning;
            state.step = next;
          });
          navigate(next);

          const scatter = useScatterStore.getState();
          const config = {
            source: scatter.source,
            targets: Object.keys(scatter.targets),
            selected: scatter.selected,
          };

          socket.send(
            JSON.stringify({
              topic: Topic.CommandScatterPlanStart,
              payload: config,
            }),
          );
        },
        scatterProgress: (payload: string) => {
          // console.log('scatterProgress ', payload);
          useScatterStore.getState().actions.addLine(payload);
        },
        scatterPlanEnded: (payload: string) => {
          console.log('scatterPlanEnded ', payload);
          get().actions.getUnraid();
        },
      },
    };
  }),
);

export const useUnraidActions = () => useUnraidStore().actions;

export const useUnraidLoaded = () => useUnraidStore().loaded;
export const useUnraidStatus = () => useUnraidStore().status;
export const useUnraidStep = () => useUnraidStore().step;
export const useUnraidRoute = () => useUnraidStore().route;
export const useUnraidIsBusy = () =>
  ![Op.Neutral, Op.ScatterPlan, Op.GatherPlan].includes(
    useUnraidStore().status,
  );
export const useUnraidDisks = () => useUnraidStore().unraid?.disks ?? [];
