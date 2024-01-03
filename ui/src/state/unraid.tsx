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
  Packet,
  Topic,
  State,
} from '~/types';
import { getRouteFromStatus } from '~/helpers/routes';
import { useScatterStore } from '~/state/scatter';
import { useGatherStore } from '~/state/gather';
import { createMachine, StateMachine } from '~/helpers/sm';

interface UnraidStore {
  socket: WebSocket;
  navigate: NavigateFunction | null;
  machine: StateMachine;
  loaded: boolean;
  route: string;
  status: Op;
  unraid: Unraid | null;
  operation: Operation | null;
  history: History | null;
  plan: Plan | null;
  logs: Array<string>;
  actions: {
    setNavigate: (navigate: NavigateFunction) => void;
    getUnraid: () => Promise<void>;
    refreshUnraid: () => Promise<void>;
    syncRoute: (path: string) => void;
    transition: (event: string) => void;
    scatterPlan: () => void;
    scatterProgress: (payload: string) => void;
    scatterPlanEnded: (payload: Plan) => void;
    scatterOperation: (
      command: Topic.CommandScatterMove | Topic.CommandScatterCopy,
    ) => void;
    scatterValidate: (operation: Operation | undefined) => void;
    transferProgress: (payload: Operation) => void;
    transferEnded: (payload: State) => void;
    gatherPlan: () => void;
    gatherProgress: (payload: string) => void;
    gatherPlanEnded: (payload: Plan) => void;
    gatherMove: () => void;
  };
}

const mapEventToAction: { [x: string]: string } = {
  [Topic.EventScatterPlanStarted]: 'scatterProgress',
  [Topic.EventScatterPlanProgress]: 'scatterProgress',
  [Topic.EventScatterPlanEnded]: 'scatterPlanEnded',
  [Topic.EventTransferStarted]: 'transferProgress',
  [Topic.EventTransferProgress]: 'transferProgress',
  [Topic.EventTransferEnded]: 'transferEnded',
  [Topic.EventGatherPlanStarted]: 'gatherProgress',
  [Topic.EventGatherPlanProgress]: 'gatherProgress',
  [Topic.EventGatherPlanEnded]: 'gatherPlanEnded',
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

    const machine = {
      initialState: '/',
      '/scatter/select': {
        next: {
          target: '/scatter/plan',
          action() {
            console.log(
              'transition action for "next" in "/scatter/select" state',
            );
            get().actions.scatterPlan();
          },
        },
      },
      '/scatter/plan': {
        next: {
          target: '/scatter/transfer/validation',
          action() {
            console.log(
              'transition action for "next" in "/scatter/plan" state',
            );
            // get().actions.refreshUnraid();
          },
        },
      },
      '/scatter/transfer/validation': {
        next: {
          target: '/scatter/transfer/operation',
          action() {
            console.log(
              'transition action for "next" in "/scatter/transfer/validation" state',
            );
          },
        },
      },
      '/gather/select': {
        next: {
          target: '/gather/plan',
          action() {
            console.log(
              'transition action for "next" in "/gather/select" state',
            );
            get().actions.gatherPlan();
          },
        },
      },
      '/gather/plan': {
        next: {
          target: '/gather/transfer/targets',
          action() {
            console.log('transition action for "next" in "/gather/plan" state');
            // get().actions.refreshUnraid();
          },
        },
      },
      '/gather/transfer/targets': {
        next: {
          target: '/gather/transfer/operation',
          action() {
            console.log(
              'transition action for "next" in "/gather/transfer/targets" state',
            );
          },
        },
      },
    };

    return {
      socket,
      navigate: null,
      machine: createMachine(machine),
      loaded: false,
      route: '/',
      status: Op.Neutral,
      unraid: null,
      operation: null,
      history: null,
      plan: null,
      logs: [],
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
          });

          if (array.status === Op.Neutral) {
            return;
          }

          console.log('navigating to ', getRouteFromStatus(array.status));
          get().navigate?.(getRouteFromStatus(array.status));
        },
        refreshUnraid: async () => {
          const array = await Api.getUnraid();

          console.log('refreshUnraid ... ', array);
          set((state) => {
            state.status = array.status;
            state.unraid = array.unraid;
            state.operation = array.operation;
            state.history = array.history;
            // state.plan = array.plan;
          });
        },
        syncRoute: (path: string) => {
          set({ route: path });
        },
        transition: (event: string) => {
          const machine = get().machine;
          const route = machine.transition(get().route, event);
          // console.log('unraid.transition ', get().route, event, route);
          get().navigate?.(route);
        },
        scatterPlan: () => {
          console.log('running scatter plan');
          set((state) => {
            state.status = Op.ScatterPlan;
            state.logs = [];
            state.plan = null;
          });

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
          // useScatterStore.getState().actions.addLine(payload);
          set((state) => {
            state.logs.push(payload);
          });
        },
        scatterPlanEnded: (payload: Plan) => {
          console.log('scatterPlanEnded ', payload);
          set((state) => {
            state.status = Op.Neutral;
            state.plan = payload;
          });
          // get().actions.getUnraid();
        },
        scatterOperation: (
          command: Topic.CommandScatterMove | Topic.CommandScatterCopy,
        ) => {
          const machine = get().machine;
          const route = machine.transition(get().route, 'next');
          // console.log('unraid.transition ', get().route, event, route);
          const socket = get().socket;
          const plan = get().plan;
          socket.send(
            JSON.stringify({
              topic: command,
              payload: plan,
            }),
          );
          get().navigate?.(route);
        },
        scatterValidate: (operation: Operation | undefined) => {
          if (!operation) {
            return;
          }

          set((state) => {
            state.plan = null;
            state.operation = null;
            state.logs = [];
          });
          const socket = get().socket;
          socket.send(
            JSON.stringify({
              topic: Topic.CommandScatterValidate,
              payload: operation,
            }),
          );
          get().navigate?.('/scatter/transfer/operation');
        },
        transferProgress: (payload: Operation) => {
          // console.log('transferProgress ', payload);
          set((state) => {
            state.operation = payload;
          });
        },
        transferEnded: (payload: State) => {
          // console.log('transferProgress ', payload);
          set((state) => {
            state.status = payload.status;
            state.unraid = payload.unraid;
            state.operation = payload.operation;
            state.history = payload.history;
          });
        },
        gatherPlan: () => {
          console.log('running gather plan');
          set((state) => {
            state.status = Op.GatherPlan;
            state.logs = [];
            state.plan = null;
          });

          const gather = useGatherStore.getState();
          const config = {
            selected: Object.values(gather.selected),
          };

          socket.send(
            JSON.stringify({
              topic: Topic.CommandGatherPlanStart,
              payload: config,
            }),
          );
        },
        gatherProgress: (payload: string) => {
          // console.log('scatterProgress ', payload);
          // useGatherStore.getState().actions.addLine(payload);
          set((state) => {
            state.logs.push(payload);
          });
        },
        gatherPlanEnded: (payload: Plan) => {
          console.log('gatherPlanEnded ', payload);
          set((state) => {
            state.status = Op.Neutral;
            state.plan = payload;
          });
          // get().actions.getUnraid();
        },
        gatherMove: () => {
          const machine = get().machine;
          const route = machine.transition(get().route, 'next');
          // console.log('unraid.transition ', get().route, event, route);
          const socket = get().socket;
          const plan = get().plan;

          if (!plan) {
            return;
          }

          const target = useGatherStore.getState().target;

          socket.send(
            JSON.stringify({
              topic: Topic.CommandGatherMove,
              payload: { ...plan, target },
            }),
          );
          get().navigate?.(route);
        },
      },
    };
  }),
);

export const useUnraidActions = () => useUnraidStore((state) => state.actions);

export const useUnraidLoaded = () => useUnraidStore((state) => state.loaded);
export const useUnraidStatus = () => useUnraidStore((state) => state.status);
export const useUnraidRoute = () => useUnraidStore((state) => state.route);
export const useUnraidIsBusy = () =>
  useUnraidStore((state) => state.status !== Op.Neutral);
export const useUnraidDisks = () =>
  useUnraidStore((state) => state.unraid?.disks ?? []);
export const useUnraidPlan = () => useUnraidStore((state) => state.plan);
export const useUnraidOperation = () =>
  useUnraidStore((state) => state.operation);
export const useUnraidLogs = () => useUnraidStore((state) => state.logs);
export const useUnraidHistory = () =>
  useUnraidStore((state) =>
    state.history
      ? // @ts-expect-error -- TSCONVERSION
        state.history.order.map((id) => state.history.items[id])
      : [],
  );
