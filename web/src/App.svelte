<script lang="ts">
  import {
    store,
    NotFound,
    ResetPassword,
    SignIn,
    SignUp,
    Home,
    Redirect,
  } from "./lib";
</script>

{#if store.user === undefined}
  {null}
{:else if store.location === "/"}
  {#if store.user === null}
    <Redirect to="/sign-in" />
  {:else if store.location === "/"}
    <Home user={store.user} />
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
