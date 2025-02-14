<script lang="ts">
  import type { ChildrenProps } from "./types";

  interface Props extends ChildrenProps {
    submit?: boolean;
  }

  let { children, submit }: Props = $props();
</script>

<button type={submit ? "submit" : null}>{@render children()}</button>

<style>
  button {
    font-family: inherit;
    font-weight: 500;
    font-size: 1.15rem;
    border-radius: var(--border-radius);
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
