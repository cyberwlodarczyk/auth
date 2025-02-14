<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends HTMLInputAttributes {
    id: string;
    value: string;
  }

  let { value = $bindable(), ...props }: Props = $props();
</script>

<input bind:value {...props} />

<style>
  input {
    all: unset;
    font-size: 1.15rem;
    border-radius: 4px;
    width: 225px;
    padding: 0.5rem 1rem;
    transition: background-color 0.2s;
    background-color: var(--surface-1);
  }

  input:focus {
    background-color: var(--surface-2);
  }
</style>
