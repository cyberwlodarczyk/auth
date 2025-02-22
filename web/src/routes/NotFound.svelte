<script lang="ts">
  import { Heading } from "$lib/components";
</script>

<svelte:head>
  <title>Not found</title>
</svelte:head>

<main>
  <p>404</p>
  <Heading>Page not found</Heading>
</main>

<style>
  p {
    text-align: center;
    font-size: 3rem;
    font-weight: 700;
    color: var(--primary-1);
  }
</style>
