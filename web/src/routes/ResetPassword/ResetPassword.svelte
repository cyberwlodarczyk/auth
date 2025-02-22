<script lang="ts">
  import { CheckMail } from "$lib/components";
  import { newQueryToken } from "$lib/data";
  import { ResetPasswordFinish, ResetPasswordInit } from ".";

  const token = newQueryToken();

  let mail = $state(false);
</script>

<svelte:head>
  <title>Reset password</title>
</svelte:head>

{#if token !== undefined}
  {#if token === null}
    {null}
    <!-- TODO: token is invalid or expired -->
  {:else}
    <ResetPasswordFinish token={token.raw} />
  {/if}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <ResetPasswordInit bind:mail />
{/if}
