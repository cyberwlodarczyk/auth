<script lang="ts">
  import { Button, Field, Form, Heading, Link } from "../styled";

  const state = $state({
    email: "",
    password: "",
  });

  function onsubmit() {}
</script>

<main>
  <Form {onsubmit}>
    <Heading>Sign In</Heading>
    <Field id="email" label="Email" bind:value={state.email} type="email" />
    <Field
      id="password"
      label="Password"
      bind:value={state.password}
      type="password"
    />
    <Button submit>Sign In</Button>
  </Form>
  <Link href="/sign-up">Sign Up</Link>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    align-items: center;
  }
</style>
