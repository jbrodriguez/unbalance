import React from 'react';

import { useLocation, useNavigate, Outlet } from 'react-router-dom';
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
import {
  useUnraidActions,
  useUnraidLoaded,
} from './state/unraid';

export function App() {
  const { getConfig } = useConfigActions();
  const { load } = useAuthActions();
  const { setNavigate, getUnraid, syncRoute, connectSocket, disconnectSocket } =
    useUnraidActions();
  const isLoaded = useUnraidLoaded();
  const version = useConfigVersion();
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
    // console.log('sync location >>>>>>>>>>>>> ', location);
    syncRoute(location.pathname);
  }, [location, syncRoute]);

  React.useEffect(() => {
    Api.setCSRFToken(csrfToken);
  }, [csrfToken]);

  React.useEffect(() => {
    const shouldBlockForAuth =
      authFailed || (authLoaded && !authenticated && (authEnabled || csrfToken === ''));

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
    authLoaded,
    authEnabled,
    authenticated,
    connectSocket,
    disconnectSocket,
    getConfig,
    getUnraid,
  ]);

  if (!authLoaded) {
    return null;
  }

  const shouldShowAuthGate =
    authFailed || (authLoaded && !authenticated && (authEnabled || csrfToken === ''));

  if (shouldShowAuthGate) {
    return (
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <AuthGate />
      </ThemeProvider>
    );
  }

  if (!(isLoaded && version !== '')) {
    return null;
  }

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <ErrorBoundary>
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
      </ErrorBoundary>
    </ThemeProvider>
  );
}
