<script lang="ts">
  import CheckMail from "./CheckMail.svelte";
  import ResetPasswordFinish from "./ResetPasswordFinish.svelte";
  import ResetPasswordInit from "./ResetPasswordInit.svelte";

  const params = new URLSearchParams(window.location.search);
  const token = params.get("token");

  let mail = $state(false);
</script>

<svelte:head>
  <title>Reset password</title>
</svelte:head>

{#if token}
  <ResetPasswordFinish {token} />
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <ResetPasswordInit bind:mail />
{/if}
