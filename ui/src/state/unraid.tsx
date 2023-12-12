import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import { Api } from '~/api';

interface UnraidStore {
  array: string;
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
      array: '',
      actions: {
        getUnraid: async () => {
          const array = await Api.getUnraid();
          console.log('useUnraidStore.getUnraid() ', array);
          set((state) => {
            state.array = array;
          });
        },
      },
    };
  }),
);
