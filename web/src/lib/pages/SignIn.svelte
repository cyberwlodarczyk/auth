<script lang="ts">
  import { Button, Field, Form, Heading, Link, BottomLink } from "../styled";
  import { defaultFieldState, isFieldEmpty, isEmailInvalid } from "../utils";

  let email = $state(defaultFieldState());
  let password = $state(defaultFieldState());

  function onsubmit() {
    if (
      isFieldEmpty(email) ||
      isEmailInvalid(email) ||
      isFieldEmpty(password)
    ) {
      return;
    }
    console.log(email.value, password.value);
  }
</script>

<svelte:head>
  <title>Sign in</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign in</Heading>
    <Field
      id="email"
      label="Email"
      bind:value={email.value}
      bind:error={email.error}
      type="email"
    />
    <Field
      id="password"
      label="Password"
      bind:value={password.value}
      bind:error={password.error}
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
