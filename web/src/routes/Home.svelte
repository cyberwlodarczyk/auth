<script lang="ts">
  import { store } from "$lib/data";
  import { Heading } from "$lib/components";
</script>

<svelte:head>
  <title>Home</title>
</svelte:head>

<main>
  <Heading>Hi, {store.user!.name}!</Heading>
  {#if store.sudo}
    <p role="alert">Sudo</p>
  {/if}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  p {
    text-transform: uppercase;
    color: var(--error);
    background-color: var(--error-transparent);
    border-radius: var(--border-radius);
    padding: 0.375rem 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    letter-spacing: 0.03125rem;
  }
</style>
