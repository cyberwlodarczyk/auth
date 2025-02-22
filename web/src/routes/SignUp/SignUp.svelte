<script lang="ts">
  import { CheckMail } from "$lib/components";
  import { newQueryToken } from "$lib/data";
  import { SignUpFinish, SignUpInit } from ".";

  const token = newQueryToken();

  let mail = $state(false);
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

{#if token !== undefined}
  {#if token === null}
    {null}
    <!-- TODO: token is invalid or expired -->
  {:else}
    <SignUpFinish token={token.raw} />
  {/if}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <SignUpInit bind:mail />
{/if}
