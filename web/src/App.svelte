<script lang="ts">
  import { store, NotFound, ResetPassword, SignIn, SignUp, Home } from "./lib";
</script>

{#if store.location === "/"}
  <Home />
{:else if store.location === "/sign-in"}
  <SignIn />
{:else if store.location === "/sign-up"}
  <SignUp />
{:else if store.location === "/reset-password"}
  <ResetPassword />
{:else}
  <NotFound />
{/if}
