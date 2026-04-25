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
    <div className="min-h-screen bg-[#131a2d] flex items-center justify-center px-4">
      <div className="w-full max-w-lg border border-slate-300/80 bg-[#243045] shadow-sm px-8 py-9 space-y-7 text-slate-100">
        <div className="flex items-center justify-center gap-4">
          <img src={logo} alt="unbalanced" className="h-14 w-14" />
          <div>
            <h1 className="text-[2rem] leading-none font-semibold tracking-tight text-slate-200">
              unbalanced
            </h1>
            <p className="mt-2 text-lg leading-none text-slate-400">
              {isSetup ? 'Create the admin password' : 'Sign in to continue'}
            </p>
          </div>
        </div>

        {isSetup && (
          <p className="text-sm leading-6 text-slate-300">
            Authentication is enabled, but no admin password has been created
            yet. The first administrator to open the app must complete setup.
          </p>
        )}

        <form className="space-y-5" onSubmit={onSubmit}>
          <div className="space-y-2">
            <Label
              htmlFor="username"
              className="text-lg font-medium text-slate-200"
            >
              Username
            </Label>
            <Input
              id="username"
              autoComplete="username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              disabled={submitting}
              className="h-11 border-transparent bg-[#2b3850] shadow-[inset_0_1px_0_rgba(255,255,255,0.03),0_0_0_1px_rgba(18,24,38,0.35)] text-lg text-slate-100 placeholder:text-slate-500 focus-visible:border-transparent focus-visible:ring-1 focus-visible:ring-blue-400"
            />
          </div>

          <div className="space-y-2">
            <Label
              htmlFor="password"
              className="text-lg font-medium text-slate-200"
            >
              Password
            </Label>
            <Input
              id="password"
              autoComplete={isSetup ? 'new-password' : 'current-password'}
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              disabled={submitting}
              className="h-11 border-transparent bg-[#2b3850] shadow-[inset_0_1px_0_rgba(255,255,255,0.03),0_0_0_1px_rgba(18,24,38,0.35)] text-lg text-slate-100 placeholder:text-slate-500 focus-visible:border-transparent focus-visible:ring-1 focus-visible:ring-blue-400"
            />
          </div>

          {isSetup && (
            <div className="space-y-2">
              <Label
                htmlFor="confirm-password"
                className="text-lg font-medium text-slate-200"
              >
                Confirm Password
              </Label>
              <Input
                id="confirm-password"
                autoComplete="new-password"
                type="password"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
                disabled={submitting}
                className="h-11 border-transparent bg-[#2b3850] shadow-[inset_0_1px_0_rgba(255,255,255,0.03),0_0_0_1px_rgba(18,24,38,0.35)] text-lg text-slate-100 placeholder:text-slate-500 focus-visible:border-transparent focus-visible:ring-1 focus-visible:ring-blue-400"
              />
              {confirmPassword !== '' && confirmPassword !== password && (
                <p className="text-sm text-red-400">
                  Passwords do not match.
                </p>
              )}
            </div>
          )}

          {authError !== '' && (
            <p className="text-sm text-red-400">
              {authError}
            </p>
          )}

          <Button
            className="mt-2 h-12 w-full bg-slate-300 text-slate-900 text-xl hover:bg-slate-200 disabled:bg-slate-400 disabled:text-slate-700"
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
