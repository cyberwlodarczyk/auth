<script lang="ts">
  import { store } from "../store.svelte";
  import CheckMail from "./CheckMail.svelte";
  import EnterSudoModeInit from "./EnterSudoModeInit.svelte";

  const params = new URLSearchParams(window.location.search);
  const token = params.get("token");

  $effect(() => {
    if (token) {
      store.sudo = token;
      store.location = "/";
    }
  });

  let mail = $state(false);
</script>

<svelte:head>
  <title>Enter sudo mode</title>
</svelte:head>

{#if token}
  {null}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <EnterSudoModeInit bind:mail />
{/if}
