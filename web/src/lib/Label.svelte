<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLLabelAttributes } from "svelte/elements";

  interface Props extends HTMLLabelAttributes {
    children: Snippet<[]>;
    for: string;
  }

  let { children, ...props }: Props = $props();
</script>

<label {...props}>{@render children()}</label>

<style>
  label {
    font-size: 1rem;
    font-weight: 500;
  }
</style>
