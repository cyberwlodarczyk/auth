<script lang="ts">
  import { store, newQueryToken } from "$lib/data";
  import { CheckMail } from "$lib/components";
  import { Route } from "$lib/config";
  import { EnterSudoModeInit } from ".";

  const token = newQueryToken();

  $effect(() => {
    if (token) {
      store.sudo = token;
      store.location = Route.Home;
    }
  });

  let mail = $state(false);
</script>

<svelte:head>
  <title>Enter sudo mode</title>
</svelte:head>

{#if token !== undefined}
  {#if token === null}
    {null}
    <!-- TODO: token is invalid or expired -->
  {:else}
    {null}
  {/if}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <EnterSudoModeInit bind:mail />
{/if}
