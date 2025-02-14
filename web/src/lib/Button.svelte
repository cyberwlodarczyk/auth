<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";

  interface Props extends HTMLButtonAttributes {
    children: Snippet<[]>;
  }

  let { children, ...props }: Props = $props();
</script>

<button {...props}>{@render children()}</button>

<style>
  button {
    font-family: inherit;
    font-weight: 500;
    font-size: 1.15rem;
    border-radius: 4px;
    padding: 0.5rem 1.25rem;
    border: none;
    outline: none;
    cursor: pointer;
    transition:
      background-color 0.2s,
      box-shadow 0.2s;
    color: var(--background);
    background-color: var(--primary-1);
    box-shadow: var(--shadow-1);
  }

  button:hover,
  button:focus-visible {
    background-color: var(--primary-2);
    box-shadow: var(--shadow-2);
  }
</style>
