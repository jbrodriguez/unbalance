import React from 'react';

import { Button } from '~/shared/buttons/button';
import { Icon } from '~/shared/icons/icon';
import { Stepper } from '~/shared/stepper/stepper';
import {
  useUnraidActions,
  useUnraidRoute,
  useUnraidIsBusy,
} from '~/state/unraid';
import { useGatherSelected, useGatherTarget } from '~/state/gather';
import { routeToStep } from '~/helpers/routes';
import { getVariant, getFill } from '~/helpers/styling';
import { useConfigActions, useConfigDryRun } from '~/state/config';
import { Actions } from '~/shared/transfer/actions';

const config = [
  { navTo: 'select', title: 'Select', subtitle: 'Choose source' },
  { navTo: 'plan', title: 'Plan', subtitle: 'Monitor' },
  {
    navTo: 'transfer',
    title: 'Transfer',
    subtitle: 'Choose target & Move',
  },
];

export const Navbar: React.FunctionComponent = () => {
  const route = useUnraidRoute();
  const { transition, gatherMove } = useUnraidActions();
  const target = useGatherTarget();
  const { toggleDryRun } = useConfigActions();
  const dryRun = useConfigDryRun();
  const busy = useUnraidIsBusy();
  const selected = useGatherSelected();

  const onNext = () => transition('next');
  const onPrev = () => transition('prev');
  const onMove = () => gatherMove();
  const onDryRun = () => toggleDryRun();

  const currentStep = routeToStep(route);
  const nextDisabled =
    busy ||
    route === '/gather/transfer/targets' ||
    (route === '/gather/select' && Object.keys(selected).length === 0);

  return (
    <div className="flex flex-row items-center justify-between mb-4">
      <div className="flex justify-start">
        <Button
          label="Prev"
          variant={getVariant(route !== '/gather/select')}
          leftIcon={
            <Icon
              name="prev"
              size={20}
              style={getFill(route !== '/gather/select')}
            />
          }
          disabled={busy || route === '/gather/select'}
          onClick={onPrev}
        />
      </div>

      <div className="flex flex-row flex-1 items-center justify-between">
        <div className="flex flex-row items-center justify-start">
          <span className="mx-2" />
          <Stepper steps={3} currentStep={currentStep} config={config} />
        </div>

        {route === '/gather/transfer/targets' && target !== '' && (
          <div className="flex flex-row items-center justify-end">
            <Button label="MOVE" variant="primary" onClick={onMove} />
            <span className="mx-1">|</span>

            <div className="flex items-center">
              <input
                checked={dryRun}
                id="checked-checkbox"
                type="checkbox"
                className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                onChange={onDryRun}
              />
              <label
                htmlFor="checked-checkbox"
                className="ms-2 text-sm font-medium"
              >
                dry-run
              </label>
            </div>
            <span className="mx-2" />
          </div>
        )}
        {route === '/gather/transfer/operation' && <Actions />}
      </div>

      <div className="flex items-center justify-end">
        <Button
          label="Next"
          variant={getVariant(route !== '/gather/transfer/targets')}
          rightIcon={
            <Icon
              name="next"
              size={20}
              style={getFill(route !== '/gather/transfer/targets')}
            />
          }
          disabled={nextDisabled}
          onClick={onNext}
        />
      </div>
    </div>
  );
};
