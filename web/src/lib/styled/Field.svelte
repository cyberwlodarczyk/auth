<script lang="ts">
  import type { HTMLInputTypeAttribute } from "svelte/elements";
  import type { FieldState } from "../utils";

  interface Props extends FieldState {
    id: string;
    label: string;
    type?: HTMLInputTypeAttribute;
  }

  let {
    id,
    value = $bindable(),
    error = $bindable(),
    label,
    type,
  }: Props = $props();
  let errorId = `${id}-error`;
</script>

<div>
  <label for={id}>{label}</label>
  <input
    {id}
    {type}
    bind:value
    oninput={() => {
      error = null;
    }}
    aria-required={true}
    aria-invalid={error ? true : null}
    aria-describedby={error ? errorId : null}
    autocomplete="on"
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
    width: 17rem;
  }

  label {
    font-size: 0.875rem;
    font-weight: 500;
    margin-left: var(--border-radius);
  }

  input {
    border: none;
    outline: none;
    font-size: 1rem;
    font-family: inherit;
    color: inherit;
    width: 100%;
    border-radius: var(--border-radius);
    padding: 0.5rem 1rem;
    transition: background-color 0.2s;
    background-color: var(--primary-transparent-1);
  }

  input:focus {
    background-color: var(--primary-transparent-2);
  }

  p {
    color: var(--error);
    font-size: 0.875rem;
    margin-left: var(--border-radius);
  }
</style>
