<script lang="ts">
  import { store } from "../store.svelte";

  interface Props {
    to: string;
  }

  let { to }: Props = $props();

  $effect(() => {
    store.location = to;
  });
</script>
