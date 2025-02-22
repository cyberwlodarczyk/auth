<script lang="ts">
  import { store } from "$lib/data";
  import { Redirect } from "$lib/components";
  import { Route } from "$lib/config";
  import {
    NotFound,
    ResetPassword,
    SignIn,
    SignUp,
    Home,
    ChangePassword,
    ChangeEmail,
    EnterSudoMode,
  } from "./routes";

  const PUBLIC_ROUTES: string[] = [
    Route.SignUp,
    Route.SignIn,
    Route.ResetPassword,
  ];
  const SESSION_ROUTES: string[] = [
    Route.Home,
    Route.EnterSudoMode,
    Route.ChangePassword,
  ];
  const SUDO_ROUTES: string[] = [Route.ChangeEmail];
</script>

{#if store.user === undefined}
  {null}
{:else if PUBLIC_ROUTES.includes(store.location)}
  {#if store.user !== null}
    <Redirect to={Route.Home} />
  {:else if store.location === Route.SignUp}
    <SignUp />
  {:else if store.location === Route.SignIn}
    <SignIn />
  {:else if store.location === Route.ResetPassword}
    <ResetPassword />
  {/if}
{:else if SESSION_ROUTES.includes(store.location)}
  {#if store.user === null}
    <Redirect to={Route.SignIn} />
  {:else if store.location === Route.Home}
    <Home />
  {:else if store.location === Route.EnterSudoMode}
    <EnterSudoMode />
  {:else if store.location === Route.ChangePassword}
    <ChangePassword />
  {/if}
{:else if SUDO_ROUTES.includes(store.location)}
  {#if store.sudo === null}
    <Redirect to={Route.EnterSudoMode} />
  {:else if store.location === Route.ChangeEmail}
    <ChangeEmail />
  {/if}
{:else}
  <NotFound />
{/if}
