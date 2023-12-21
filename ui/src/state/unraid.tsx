import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { NavigateFunction } from 'react-router-dom';

import { Api } from '~/api';
import { Unraid, Operation, History, Plan, Op, Packet, Topic } from '~/types';
import {
  getNextRoute,
  getRouteFromStatus,
  // getBaseRoute,
} from '~/helpers/routes';
import { useScatterStore } from '~/state/scatter';
// import {
//   CommandScatterPlanStart,
//   EventScatterPlanStarted,
//   EventScatterPlanProgress,
//   EventScatterPlanEnded,
// } from '~/constants';

interface UnraidStore {
  socket: WebSocket;
  navigate: NavigateFunction | null;
  loaded: boolean;
  route: string;
  status: Op;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
  actions: {
    setNavigate: (navigate: NavigateFunction) => void;
    getUnraid: () => Promise<void>;
    syncRoute: (path: string) => void;
    transition: (from: string) => void;
    scatterProgress: (payload: string) => void;
    scatterPlanEnded: (payload: string) => void;
  };
}

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
      navigate: null,
      loaded: false,
      route: '/',
      status: Op.Neutral,
      unraid: null,
      operation: null,
      history: null,
      plan: null,
      actions: {
        setNavigate: (navigate: NavigateFunction) => {
          set((state) => {
            state.navigate = navigate;
          });
        },
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
            // state.route = getRouteFromStatus(array.status);
            // state.step = convertStatusToStep(array.status);
          });

          if (array.status === Op.Neutral) {
            return;
          }

          console.log('navigating to ', getRouteFromStatus(array.status));
          get().navigate?.(getRouteFromStatus(array.status));
        },
        syncRoute: (path: string) => {
          if (!get().loaded) {
            return;
          }

          // set((state) => {
          //   state.route = path;
          //   // state.step = 'select';
          // });
          set({ route: path });

          // const route = get().route;
          // // if (route.slice(0, 5) === path.slice(0, 5)) {
          // //   return;
          // // }
          // if (route === path) {
          //   return;
          // }

          // const next = getNextRoute(path);
          // console.log('next, route, path', next, route, path);
          // if (next === route) {
          //   return;
          // }

          // // don't sync if we're going to the same route
          // console.log(
          //   'base-route, base-path',
          //   getBaseRoute(route),
          //   getBaseRoute(next),
          // );
          // if (getBaseRoute(route) === getBaseRoute(next)) {
          //   return;
          // }

          // console.log('new route ', next);
          // set((state) => {
          //   state.route = next;
          //   // state.step = 'select';
          // });
        },
        transition: (from: string) => {
          const next = getNextRoute(get().route);
          set((state) => {
            // state.status = Op.ScatterPlanning;
            state.route = next;
          });
          // navigate(next);

          if (from === '/scatter/select') {
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
          }
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
export const useUnraidRoute = () => useUnraidStore().route;
export const useUnraidIsBusy = () =>
  ![Op.Neutral, Op.ScatterPlan, Op.GatherPlan].includes(
    useUnraidStore().status,
  );
export const useUnraidDisks = () => useUnraidStore().unraid?.disks ?? [];
