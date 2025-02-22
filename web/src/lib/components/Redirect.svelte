<script lang="ts">
  import { store } from "$lib/data";
  import type { Route } from "$lib/config";

  interface Props {
    to: Route;
  }

  let { to }: Props = $props();

  $effect(() => {
    store.location = to;
  });
</script>
