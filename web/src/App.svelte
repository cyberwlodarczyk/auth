<script lang="ts">
  import {
    store,
    NotFound,
    ResetPassword,
    SignIn,
    SignUp,
    Home,
    ChangePassword,
    ChangeEmail,
    Redirect,
    EnterSudoMode,
  } from "./lib";
</script>

{#if store.user === undefined}
  {null}
{:else if ["/", "/enter-sudo-mode", "/change-password", "/change-email"].includes(store.location)}
  {#if store.user === null}
    <Redirect to="/sign-in" />
  {:else if store.location === "/"}
    <Home />
  {:else if store.location === "/enter-sudo-mode"}
    <EnterSudoMode />
  {:else if store.location === "/change-password"}
    <ChangePassword />
  {:else if store.location === "/change-email"}
    {#if store.sudo === null}
      <Redirect to="/enter-sudo-mode" />
    {:else}
      <ChangeEmail />
    {/if}
  {/if}
{:else if ["/sign-up", "/sign-in", "/reset-password"].includes(store.location)}
  {#if store.user !== null}
    <Redirect to="/" />
  {:else if store.location === "/sign-up"}
    <SignUp />
  {:else if store.location === "/sign-in"}
    <SignIn />
  {:else if store.location === "/reset-password"}
    <ResetPassword />
  {/if}
{:else}
  <NotFound />
{/if}
