import React from 'react';

import { Button } from '~/shared/buttons/button';
import { Icon } from '~/shared/icons/icon';
import { Stepper } from '~/shared/stepper/stepper';
import {
  useUnraidActions,
  useUnraidRoute,
  useUnraidIsBusy,
} from '~/state/unraid';
import { routeToStep } from '~/helpers/routes';
import { getVariant, getFill } from '~/helpers/styling';
import { Topic } from '~/types';
import { useConfigActions, useConfigDryRun } from '~/state/config';
import { useScatterSelected, useScatterTargets } from '~/state/scatter';

const config = [
  { navTo: 'select', title: 'Select', subtitle: 'Choose data' },
  { navTo: 'plan', title: 'Plan', subtitle: 'Monitor' },
  { navTo: 'transfer', title: 'Transfer', subtitle: 'Validate & Run' },
];

export const Navbar: React.FunctionComponent = () => {
  const route = useUnraidRoute();
  const currentStep = routeToStep(route);
  const { transition, scatterOperation } = useUnraidActions();
  const { toggleDryRun } = useConfigActions();
  const dryRun = useConfigDryRun();
  const selected = useScatterSelected();
  const targets = useScatterTargets();
  const busy = useUnraidIsBusy();

  const onNext = () => transition('next');
  const onMove = () => scatterOperation(Topic.CommandScatterMove);
  const onCopy = () => scatterOperation(Topic.CommandScatterCopy);
  const onDryRun = () => toggleDryRun();

  const nextDisabled =
    busy ||
    route === '/scatter/transfer/validation' ||
    (route === '/scatter/select' &&
      (selected.length === 0 || Object.keys(targets).length === 0));

  return (
    <div className="flex flex-row items-center justify-between mb-4">
      <div className="flex justify-start">
        <Button
          label="Prev"
          variant={getVariant(route !== '/scatter/select')}
          leftIcon={
            <Icon
              name="prev"
              size={20}
              style={getFill(route !== '/scatter/select')}
            />
          }
          disabled={busy || route === '/scatter/select'}
        />
      </div>

      <div className="flex flex-row flex-1 items-center justify-between">
        <div className="flex flex-row items-center justify-start">
          <span className="mx-2" />
          <Stepper steps={3} currentStep={currentStep} config={config} />
        </div>

        {route === '/scatter/transfer/validation' && (
          <div className="flex flex-row items-center justify-end">
            <Button label="MOVE" variant="primary" onClick={onMove} />
            <span className="mx-1">|</span>
            <Button label="COPY" variant="primary" onClick={onCopy} />

            <span className="mx-1">|</span>

            <div className="flex items-center">
              <input
                checked={dryRun}
                id="checked-checkbox"
                type="checkbox"
                value=""
                className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                onChange={onDryRun}
              />
              <label
                htmlFor="checked-checkbox"
                className="ms-2 text-sm font-medium text-gray-900 dark:text-gray-300"
              >
                dry-run
              </label>
            </div>
            <span className="mx-2" />
          </div>
        )}
      </div>

      <div className="flex items-center justify-end">
        <Button
          label="Next"
          variant={getVariant(route !== '/scatter/transfer/validation')}
          rightIcon={
            <Icon
              name="next"
              size={20}
              style={getFill(route !== '/scatter/transfer/validation')}
            />
          }
          disabled={nextDisabled}
          onClick={onNext}
        />
      </div>
    </div>
  );
};
