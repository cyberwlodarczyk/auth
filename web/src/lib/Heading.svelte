<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAttributes } from "svelte/elements";

  interface Props extends HTMLAttributes<HTMLHeadingElement> {
    children: Snippet<[]>;
  }

  let { children, ...props }: Props = $props();
</script>

<h1 {...props}>{@render children()}</h1>

<style>
  h1 {
    text-align: center;
    font-size: 3rem;
  }
</style>
