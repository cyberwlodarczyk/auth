<script lang="ts">
  import { Heading } from "../styled";
  import type { User } from "../api";

  interface Props {
    user: User;
  }

  let { user }: Props = $props();
</script>

<svelte:head>
  <title>Home</title>
</svelte:head>

<Heading>Hi, {user.name}!</Heading>
