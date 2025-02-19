<script lang="ts">
  import type { ChildrenProps } from "../types";
  import { Link, Button } from ".";

  interface Props extends ChildrenProps {
    description?: string;
    href?: string;
    onclick?: () => void;
  }

  let { children, href, onclick, description }: Props = $props();
</script>

<div>
  {#if description}
    <p>{description}</p>
  {/if}
  {#if href}
    <Link {href}>{@render children()}</Link>
  {:else if onclick}
    <Button text {onclick}>{@render children()}</Button>
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
