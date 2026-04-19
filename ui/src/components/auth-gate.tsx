import React from 'react';

import logo from '~/assets/unbalance-logo.png';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  useAuthActions,
  useAuthConfigured,
  useAuthEnabled,
  useAuthError,
  useAuthUsername,
} from '~/state/auth';

export function AuthGate() {
  const enabled = useAuthEnabled();
  const configured = useAuthConfigured();
  const authError = useAuthError();
  const storedUsername = useAuthUsername();
  const { login, setup, clearError } = useAuthActions();

  const [username, setUsername] = React.useState(storedUsername || 'admin');
  const [password, setPassword] = React.useState('');
  const [confirmPassword, setConfirmPassword] = React.useState('');
  const [submitting, setSubmitting] = React.useState(false);

  React.useEffect(() => {
    setUsername(storedUsername || 'admin');
  }, [storedUsername]);

  const isSetup = enabled && !configured;

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    clearError();

    if (isSetup && password !== confirmPassword) {
      return;
    }

    setSubmitting(true);

    if (isSetup) {
      await setup(username, password);
    } else {
      await login(username, password);
    }

    setSubmitting(false);
    setPassword('');
    setConfirmPassword('');
  };

  return (
    <div className="min-h-screen bg-neutral-100 dark:bg-slate-900 flex items-center justify-center px-4">
      <div className="w-full max-w-md border bg-white dark:bg-slate-800 shadow-sm p-8 space-y-6">
        <div className="flex items-center justify-center gap-3">
          <img src={logo} alt="unbalanced" className="h-12" />
          <div>
            <h1 className="text-xl font-semibold">unbalanced</h1>
            <p className="text-sm text-slate-500 dark:text-slate-400">
              {isSetup ? 'Create the admin password' : 'Sign in to continue'}
            </p>
          </div>
        </div>

        {isSetup && (
          <p className="text-sm text-slate-600 dark:text-slate-300">
            Authentication is enabled, but no admin password has been created
            yet. The first administrator to open the app must complete setup.
          </p>
        )}

        <form className="space-y-4" onSubmit={onSubmit}>
          <div className="space-y-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              autoComplete="username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              disabled={submitting}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              autoComplete={isSetup ? 'new-password' : 'current-password'}
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              disabled={submitting}
            />
          </div>

          {isSetup && (
            <div className="space-y-2">
              <Label htmlFor="confirm-password">Confirm Password</Label>
              <Input
                id="confirm-password"
                autoComplete="new-password"
                type="password"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
                disabled={submitting}
              />
              {confirmPassword !== '' && confirmPassword !== password && (
                <p className="text-sm text-red-600 dark:text-red-400">
                  Passwords do not match.
                </p>
              )}
            </div>
          )}

          {authError !== '' && (
            <p className="text-sm text-red-600 dark:text-red-400">
              {authError}
            </p>
          )}

          <Button
            className="w-full"
            type="submit"
            disabled={
              submitting ||
              username.trim() === '' ||
              password === '' ||
              (isSetup && (confirmPassword === '' || confirmPassword !== password))
            }
          >
            {isSetup ? 'Create Password' : 'Sign In'}
          </Button>
        </form>
      </div>
    </div>
  );
}
