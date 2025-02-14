<script lang="ts">
  import type { HTMLInputTypeAttribute } from "svelte/elements";

  interface Props {
    id: string;
    value: string;
    label: string;
    error?: string;
    type?: HTMLInputTypeAttribute;
  }

  let { id, value = $bindable(), label, error, type }: Props = $props();
  let errorId = `${id}-error`;
</script>

<div>
  <label for={id}>{label}</label>
  <input
    {id}
    {type}
    bind:value
    aria-required={true}
    aria-invalid={error ? true : null}
    aria-describedby={error ? errorId : null}
  />
  {#if error}
    <p id={errorId}>{error}</p>
  {/if}
</div>

<style>
  div {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label {
    font-size: 1rem;
    font-weight: 500;
    margin-left: var(--border-radius);
  }

  input {
    all: unset;
    font-size: 1.15rem;
    border-radius: var(--border-radius);
    width: 225px;
    padding: 0.5rem 1rem;
    transition: background-color 0.2s;
    background-color: var(--surface-1);
  }

  input:focus {
    background-color: var(--surface-2);
  }

  p {
    color: var(--error);
    font-size: 0.875rem;
    margin-left: var(--border-radius);
  }
</style>
