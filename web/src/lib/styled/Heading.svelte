<script lang="ts">
  import type { ChildrenProps } from "./types";

  interface Props extends ChildrenProps {
    small?: boolean;
  }

  let { children, small }: Props = $props();
</script>

<h1 class={{ small }}>{@render children()}</h1>

<style>
  h1 {
    text-align: center;
    font-size: 2.25rem;
  }

  h1.small {
    font-size: 1.75rem;
  }
</style>
