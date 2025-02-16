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
  <title>Reset password</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Reset password</Heading>
    <Field
      id="email"
      label="Email"
      bind:value={email.value}
      bind:error={email.error}
      type="email"
    />
    <Button submit>Reset password</Button>
  </Form>
  <BottomLink href="/sign-in">Back to sign in</BottomLink>
</main>
