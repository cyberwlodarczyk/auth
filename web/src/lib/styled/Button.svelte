<script lang="ts">
  import type { ChildrenProps } from "../types";

  interface Props extends ChildrenProps {
    submit?: boolean;
    text?: boolean;
    onclick?: () => void;
  }

  let { children, submit, text, onclick }: Props = $props();
</script>

<button {onclick} class={{ text }} type={submit ? "submit" : null}
  >{@render children()}</button
>

<style>
  button {
    font-family: inherit;
    font-weight: 500;
    font-size: 1rem;
    border-radius: var(--border-radius);
    padding: 0.5rem 1.25rem;
    width: 17rem;
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

  button.text {
    font-weight: 400;
    background: transparent;
    padding: 0;
    box-shadow: none;
    color: var(--primary-1);
    font-size: 0.875rem;
    width: unset;
  }

  button.text:hover,
  button.text:focus-visible {
    text-decoration: underline;
  }
</style>
