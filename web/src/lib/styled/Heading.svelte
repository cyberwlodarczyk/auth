<script lang="ts">
  import type { ChildrenProps } from "./types";

  let { children }: ChildrenProps = $props();
</script>

<h1>{@render children()}</h1>

<style>
  h1 {
    text-align: center;
    font-size: 3rem;
  }
</style>
