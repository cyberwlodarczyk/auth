<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLFormAttributes, EventHandler } from "svelte/elements";

  interface Props extends HTMLFormAttributes {
    children: Snippet<[]>;
    onsubmit: EventHandler<SubmitEvent, HTMLFormElement>;
  }

  let { children, onsubmit, ...props }: Props = $props();
</script>

<form
  onsubmit={(event) => {
    event.preventDefault();
    onsubmit(event);
  }}
  {...props}
>
  {@render children()}
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
</style>
