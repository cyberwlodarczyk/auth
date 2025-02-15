<script lang="ts">
  import { Button, Field, Form, Heading, Link, BottomLink } from "../styled";

  const state = $state({
    email: "",
    password: "",
  });

  function onsubmit() {}
</script>

<svelte:head>
  <title>Sign in</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign in</Heading>
    <Field id="email" label="Email" bind:value={state.email} type="email" />
    <Field
      id="password"
      label="Password"
      bind:value={state.password}
      type="password"
    />
    <div style:width="100%">
      <Link href="/reset-password">Forgot password?</Link>
    </div>
    <Button submit>Sign in</Button>
  </Form>
  <BottomLink question="Don't have an account yet?" href="/sign-up">
    Sign up
  </BottomLink>
</main>
