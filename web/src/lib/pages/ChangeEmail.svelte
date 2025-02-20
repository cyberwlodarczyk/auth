<script lang="ts">
  import { changeEmail } from "../api";
  import { store } from "../store.svelte";
  import ChangeEmailInit from "./ChangeEmailInit.svelte";
  import CheckMail from "./CheckMail.svelte";

  const params = new URLSearchParams(window.location.search);
  const token = params.get("token");

  $effect(() => {
    if (token) {
      changeEmail(token).then(() => {
        store.location = "/";
      });
    }
  });

  let mail = $state(false);
</script>

<svelte:head>
  <title>Change email</title>
</svelte:head>

{#if token}
  {null}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <ChangeEmailInit bind:mail />
{/if}
