<script lang="ts">
  import { editUserEmail, newQueryToken, store } from "$lib/data";
  import { CheckMail } from "$lib/components";
  import { Route } from "$lib/config";
  import { ChangeEmailInit } from ".";

  const token = newQueryToken();

  $effect(() => {
    if (token) {
      editUserEmail({ token: token.raw }).then(() => {
        store.location = Route.Home;
      });
    }
  });

  let mail = $state(false);
</script>

<svelte:head>
  <title>Change email</title>
</svelte:head>

{#if token !== undefined}
  {#if token === null}
    {null}
    <!-- TODO: token is invalid or expired -->
  {:else}
    {null}
  {/if}
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <ChangeEmailInit bind:mail />
{/if}
