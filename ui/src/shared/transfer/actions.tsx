import React from 'react';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useToast } from '@/components/ui/use-toast';

import { useUnraidActions, useUnraidError } from '~/state/unraid';
import { useConfigActions, useConfigRefreshRate } from '~/state/config';
import { Modal } from '~/shared/modal/modal';

const refreshDiscreteToSec: Record<number, string> = {
  1000: '1 sec',
  5000: '5 sec',
  15000: '15 sec',
  30000: '30 sec',
};

export const Actions: React.FunctionComponent = () => {
  const { setRefreshRate } = useConfigActions();
  const { stop } = useUnraidActions();
  const refresh = useConfigRefreshRate();
  const error = useUnraidError();
  const { toast } = useToast();

  const [pressed, setPressed] = React.useState(false);
  const [modal, setModal] = React.useState<boolean>(false);
  React.useEffect(() => {
    if (error !== '') {
      toast({
        title: 'Operation Error',
        description: error,
        variant: 'destructive',
      });
    }
  }, [toast, error]);

  const display = refreshDiscreteToSec[refresh];

  const onStop = () => {
    setPressed(true);
    stop();
  };
  const onChangeRefreshRate = (value: string) => setRefreshRate(+value);
  const onConfirm = () => setModal(true);
  const onClose = () => setModal(false);
  const onYes = () => {
    setModal(false);
    onStop();
  };

  return (
    <div className="flex flex-row items-center justify-end pr-2">
      <Modal
        isOpen={modal}
        hasCloseBtn={false}
        onClose={onClose}
        // style="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full"
      >
        <div className="relative p-4 w-full max-w-md max-h-full bg-white dark:bg-gray-900">
          <div className="relative rounded-lg shadow">
            <button
              type="button"
              className="absolute top-3 end-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
              onClick={onClose}
            >
              <svg
                className="w-3 h-3"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 14 14"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
                />
              </svg>
              <span className="sr-only">Close modal</span>
            </button>
            <div className="p-4 md:p-5 text-center">
              <svg
                className="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 20 20"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 11V6m0 8h.01M19 10a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                />
              </svg>
              <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
                Are you sure you want to proceed with this action?
              </h3>
              <button
                type="button"
                className="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-red-300 dark:focus:ring-red-800 font-medium rounded-lg text-sm inline-flex items-center px-5 py-2.5 text-center me-2"
                onClick={onYes}
              >
                Yes, I'm sure
              </button>
              <button
                type="button"
                className="text-gray-500 bg-white hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 rounded-lg border border-gray-200 text-sm font-medium px-5 py-2.5 hover:text-gray-900 focus:z-10 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-500 dark:hover:text-white dark:hover:bg-gray-600 dark:focus:ring-gray-600"
                onClick={onClose}
              >
                No, cancel
              </button>
            </div>
          </div>
        </div>
      </Modal>
      <Button variant="destructive" onClick={onConfirm} disabled={pressed}>
        STOP
      </Button>
      <span className="px-1">|</span>
      <span>Refresh Rate: </span>
      <span className="pr-2" />
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button>{display}</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="">
          <DropdownMenuLabel>Unit</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuRadioGroup
            value={refresh.toString()}
            onValueChange={onChangeRefreshRate}
          >
            <DropdownMenuRadioItem value="1000">1 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="5000">5 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="15000">15 sec</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="30000">30 sec</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
