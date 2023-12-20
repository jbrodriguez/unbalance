import React from 'react';

import { Button } from '~/shared/buttons/button';
import { Icon } from '~/shared/icons/icon';
import { Stepper } from '~/shared/stepper/stepper';
import { useUnraidRoute } from '~/state/unraid';
import { routeToStep } from '~/helpers/routes';
import { getVariant, getFill } from '~/helpers/styling';

const config = [
  { navTo: 'select', title: 'Select', subtitle: 'Choose source' },
  { navTo: 'plan', title: 'Plan', subtitle: 'Choose destination' },
  {
    navTo: 'transfer',
    title: 'Transfer',
    subtitle: 'Verify and run operation',
  },
];

export const Navbar: React.FC = () => {
  const route = useUnraidRoute();
  const currentStep = routeToStep(route);

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
              fill={getFill(route !== '/gather/select')}
            />
          }
          disabled={route === '/gather/select'}
        />
      </div>

      <div className="flex flex-row flex-1 items-center justify-between">
        <div className="flex flex-row items-center justify-start">
          <span className="mx-2" />
          <Stepper steps={3} currentStep={currentStep} config={config} />
        </div>

        {route === '/gather/transfer' && (
          <div className="flex flex-row items-center justify-end">
            <Button label="MOVE" variant="primary" />
            <span className="mx-1">|</span>

            <div className="flex items-center">
              <input
                checked
                id="checked-checkbox"
                type="checkbox"
                value=""
                className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
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
          variant={getVariant(route !== '/gather/transfer')}
          rightIcon={
            <Icon
              name="next"
              size={20}
              fill={getFill(route !== '/gather/transfer')}
            />
          }
          disabled={route === '/gather/transfer'}
        />
      </div>
    </div>
  );
};
