import React from 'react';

import { useLocation, useNavigate, Outlet, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@/components/theme-provider';
import ErrorBoundary from '@/components/error-boundary';
import { AuthGate } from '@/components/auth-gate';
import { Api } from '~/api';

import { Header } from './shared/header/header';
import { Footer } from './shared/footer/footer';
import { useConfigActions, useConfigVersion } from './state/config';
import {
  useAuthActions,
  useAuthenticated,
  useAuthCSRFToken,
  useAuthEnabled,
  useAuthFailed,
  useAuthLoaded,
} from './state/auth';
import { useUnraidActions, useUnraidLoaded } from './state/unraid';

const DEFAULT_AUTHENTICATED_ROUTE = '/scatter/select';

interface LoginLocationState {
  from?: {
    pathname?: string;
    search?: string;
    hash?: string;
  };
}

function getReturnPath(state: unknown) {
  const locationState = state as LoginLocationState | null;
  const from = locationState?.from;

  if (!from?.pathname || from.pathname === '/login') {
    return DEFAULT_AUTHENTICATED_ROUTE;
  }

  return `${from.pathname}${from.search ?? ''}${from.hash ?? ''}`;
}

export function App() {
  const { getConfig } = useConfigActions();
  const { load } = useAuthActions();
  const { setNavigate, getUnraid, syncRoute, connectSocket, disconnectSocket } =
    useUnraidActions();
  const authLoaded = useAuthLoaded();
  const authEnabled = useAuthEnabled();
  const authFailed = useAuthFailed();
  const authenticated = useAuthenticated();
  const csrfToken = useAuthCSRFToken();
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    load();
  }, [load]);

  React.useEffect(() => {
    setNavigate(navigate);
  }, [setNavigate, navigate]);

  React.useEffect(() => {
    syncRoute(location.pathname);
  }, [location, syncRoute]);

  React.useEffect(() => {
    Api.setCSRFToken(csrfToken);
  }, [csrfToken]);

  React.useEffect(() => {
    const shouldBlockForAuth =
      authFailed || (authLoaded && authEnabled && !authenticated);

    if (!authLoaded) {
      return;
    }

    if (shouldBlockForAuth) {
      disconnectSocket();
      return;
    }

    connectSocket();
    getConfig();
    getUnraid();
  }, [
    authFailed,
    authLoaded,
    authEnabled,
    authenticated,
    connectSocket,
    disconnectSocket,
    getConfig,
    getUnraid,
  ]);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <ErrorBoundary>
        <Outlet />
      </ErrorBoundary>
    </ThemeProvider>
  );
}

export function LoginPage() {
  const authLoaded = useAuthLoaded();
  const authFailed = useAuthFailed();
  const authEnabled = useAuthEnabled();
  const authenticated = useAuthenticated();
  const navigate = useNavigate();
  const location = useLocation();
  const returnPath = getReturnPath(location.state);

  if (!authLoaded) {
    return null;
  }

  if (!authFailed && (!authEnabled || authenticated)) {
    return <Navigate to={returnPath} replace />;
  }

  return (
    <AuthGate
      onAuthenticated={() => {
        navigate(returnPath, { replace: true });
      }}
    />
  );
}

export function ProtectedLayout() {
  const isLoaded = useUnraidLoaded();
  const version = useConfigVersion();
  const authLoaded = useAuthLoaded();
  const authEnabled = useAuthEnabled();
  const authFailed = useAuthFailed();
  const authenticated = useAuthenticated();
  const location = useLocation();

  if (!authLoaded) {
    return null;
  }

  if (authFailed || (authEnabled && !authenticated)) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  if (!(isLoaded && version !== '')) {
    return null;
  }

  return (
    <div className="container mx-auto h-screen flex flex-col">
      <header>
        <Header />
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
      <footer>
        <Footer />
      </footer>
    </div>
  );
}
