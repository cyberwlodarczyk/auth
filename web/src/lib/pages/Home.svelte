<script lang="ts">
  import { store } from "../store.svelte";
  import { Heading } from "../styled";
</script>

<svelte:head>
  <title>Home</title>
</svelte:head>

<Heading>Hi, {store.user!.name}!</Heading>
{#if store.sudo}
  <p>Sudo</p>
{/if}

<style>
  p {
    text-transform: uppercase;
  }
</style>
