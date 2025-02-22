<script lang="ts">
  import { store } from "$lib/data";
  import type { ChildrenProps } from "$lib/components";
  import type { Route } from "$lib/config";

  interface Props extends ChildrenProps {
    href: Route;
  }

  const { children, href }: Props = $props();
</script>

<a
  {href}
  onclick={(event) => {
    event.preventDefault();
    store.location = href;
  }}>{@render children()}</a
>

<style>
  a {
    text-decoration: none;
    outline: none;
    font-size: 0.875rem;
    color: var(--primary-1);
  }

  a:hover,
  a:focus-visible {
    text-decoration: underline;
  }
</style>
