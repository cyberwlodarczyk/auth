<script lang="ts">
  import { signUpInit } from "../api";
  import { navigate } from "../store.svelte";
  import { Button, Field, Form, Heading, BottomLink } from "../styled";
  import { defaultFieldState, isFieldEmpty, isEmailInvalid } from "../utils";

  let email = $state(defaultFieldState());

  async function onsubmit() {
    if (isFieldEmpty(email) || isEmailInvalid(email)) {
      return;
    }
    await signUpInit(email.value);
    navigate("/sign-up/mail");
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
