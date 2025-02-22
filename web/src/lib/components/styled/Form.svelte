<script lang="ts">
  import type { EventHandler } from "svelte/elements";
  import type { ChildrenProps } from "$lib/components";

  interface Props extends ChildrenProps {
    onsubmit: EventHandler<SubmitEvent, HTMLFormElement>;
  }

  let { children, onsubmit }: Props = $props();
</script>

<form
  onsubmit={(event) => {
    event.preventDefault();
    onsubmit(event);
  }}
  novalidate
>
  {@render children()}
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }
</style>
