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
    CheckMail,
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
  {#if segments[0] === "sign-up" && segments[1] === "mail"}
    <CheckMail resend="/sign-up" title="Sign up" />
  {:else if segments[0] === "reset-password" && segments[1] === "mail"}
    <CheckMail resend="/reset-password" title="Reset password" />
  {:else}
    <NotFound />
  {/if}
{:else if segments.length === 3}
  {#if segments[0] === "sign-up" && segments[1] === "finish"}
    <SignUpFinish token={segments[2]} />
  {:else if segments[0] === "reset-password" && segments[1] === "finish"}
    <ResetPasswordFinish token={segments[2]} />
  {:else}
    <NotFound />
  {/if}
{:else}
  <NotFound />
{/if}
