<script lang="ts">
  import { Heading, SecondaryAction } from "$lib/components";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();
</script>

<main>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    class="feather feather-at-sign"
    ><circle cx="12" cy="12" r="4"></circle><path
      d="M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.92 7.94"
    ></path></svg
  >
  <Heading small>Check your mailbox for further steps</Heading>
  <SecondaryAction
    description="Didn't receive an email?"
    onclick={() => {
      mail = false;
    }}>Resend</SecondaryAction
  >
</main>

<style>
  main {
    display: flex;
    align-items: center;
    flex-direction: column;
  }

  svg {
    width: 40px;
    height: 40px;
    color: var(--primary-1);
    margin-bottom: 0.75rem;
  }
</style>
