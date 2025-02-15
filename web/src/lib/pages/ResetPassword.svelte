<script lang="ts">
  import { Button, Field, Form, Heading, Link, BottomLink } from "../styled";

  const state = $state({
    email: "",
  });

  function onsubmit() {}
</script>

<svelte:head>
  <title>Reset password</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Reset password</Heading>
    <Field id="email" label="Email" bind:value={state.email} type="email" />
    <Button submit>Reset password</Button>
  </Form>
  <BottomLink href="/sign-in">Back to sign in</BottomLink>
</main>
