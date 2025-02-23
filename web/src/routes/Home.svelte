<script lang="ts">
  import { store, type User } from "$lib/data";
  import { Heading, SudoTimer } from "$lib/components";

  interface Props {
    user: User;
  }

  let { user }: Props = $props();
</script>

<svelte:head>
  <title>Home</title>
</svelte:head>

<main>
  <Heading>Hi, {user.name}!</Heading>
  {#if store.sudo}
    <SudoTimer token={store.sudo} />
  {/if}
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }
</style>
