<script lang="ts">
  import type { ChildrenProps } from "../types";
  import { Link } from ".";

  interface Props extends ChildrenProps {
    question?: string;
    href?: string;
    onclick?: () => void;
  }

  let { children, href, onclick, question }: Props = $props();
</script>

<div>
  {#if question}
    <p>{question}</p>
  {/if}
  {#if href}
    <Link {href}>{@render children()}</Link>
  {:else if onclick}
    <button {onclick}>{@render children()}</button>
  {/if}
</div>

<style>
  div {
    margin-top: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  p {
    font-size: 0.875rem;
    color: var(--gray);
  }
</style>
