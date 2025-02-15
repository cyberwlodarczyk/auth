<script lang="ts">
  import {
    store,
    NotFound,
    ResetPasswordFinish,
    ResetPasswordInit,
    SignIn,
    SignUpFinish,
    SignUpInit,
    Home,
  } from "./lib";

  const segments = $derived(store.location.split("/").filter(Boolean));
</script>

{#if segments.length === 0}
  <Home />
{:else if segments.length === 1}
  {#if segments[0] === "sign-in"}
    <SignIn />
  {:else if segments[0] === "sign-up"}
    <SignUpInit />
  {:else if segments[0] === "reset-password"}
    <ResetPasswordInit />
  {:else}
    <NotFound />
  {/if}
{:else if segments.length === 2}
  {#if segments[0] === "sign-up"}
    <SignUpFinish token={segments[1]} />
  {:else if segments[0] === "reset-password"}
    <ResetPasswordFinish token={segments[1]} />
  {:else}
    <NotFound />
  {/if}
{:else}
  <NotFound />
{/if}
