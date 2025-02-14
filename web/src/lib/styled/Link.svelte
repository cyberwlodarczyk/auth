<script lang="ts">
  import { navigate } from "../store.svelte";
  import type { ChildrenProps } from "./types";

  interface Props extends ChildrenProps {
    href: string;
  }

  const { children, href }: Props = $props();
</script>

<a
  {href}
  onclick={(event) => {
    event.preventDefault();
    navigate(href);
  }}>{@render children()}</a
>

<style>
  a {
    text-decoration: none;
    padding: 0.5rem;
    color: var(--primary-1);
  }

  a:hover {
    text-decoration: underline;
  }
</style>
