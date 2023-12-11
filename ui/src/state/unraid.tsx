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
    const socket = new WebSocket('ws://wopr.lan:7090/ws');

    socket.onopen = function (event) {
      console.log('Socket opened ', event);
    };

    socket.onmessage = function (event) {
      console.log('Socket message ', event);
    };

    socket.onclose = function (event) {
      console.log('Socket closed ', event);
    };

    // socket.addEventListener('message', function (event) {
    //   console.log('Message from server ', event.data);
    //   set((state) => {
    //     state.state = event.data;
    //   });
    // });
    // socket.addEventListener('close', function (event) {
    //   console.log('Socket closed ', event);
    // });
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
