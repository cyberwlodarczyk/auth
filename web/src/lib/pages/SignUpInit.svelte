<script lang="ts">
  import { Button, Field, Form, Heading, BottomLink } from "../styled";
  import { defaultFieldState, isFieldEmpty, isEmailInvalid } from "../utils";

  let email = $state(defaultFieldState());

  function onsubmit() {
    if (isFieldEmpty(email) || isEmailInvalid(email)) {
      return;
    }
    console.log(email.value);
  }
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign up</Heading>
    <Field
      id="email"
      label="Email"
      bind:value={email.value}
      bind:error={email.error}
      type="email"
    />
    <Button submit>Sign up</Button>
  </Form>
  <BottomLink question="Already have an account?" href="/sign-in">
    Sign in
  </BottomLink>
</main>
