<script lang="ts">
  import CheckMail from "./CheckMail.svelte";
  import SignUpFinish from "./SignUpFinish.svelte";
  import SignUpInit from "./SignUpInit.svelte";

  const params = new URLSearchParams(window.location.search);
  const token = params.get("token");

  let mail = $state(false);
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

{#if token}
  <SignUpFinish {token} />
{:else if mail}
  <CheckMail bind:mail />
{:else}
  <SignUpInit bind:mail />
{/if}
