<script lang="ts">
  import { Button, Field, Form, Heading, BottomLink } from "../styled";

  const state = $state({
    email: "",
  });

  function onsubmit() {}
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign up</Heading>
    <Field id="email" label="Email" bind:value={state.email} type="email" />
    <Button submit>Sign up</Button>
  </Form>
  <BottomLink question="Already have an account?" href="/sign-in">
    Sign in
  </BottomLink>
</main>
