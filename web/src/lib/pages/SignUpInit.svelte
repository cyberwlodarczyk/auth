<script lang="ts">
  import { signUpInit } from "../api";
  import { store } from "../store.svelte";
  import { Button, Field, Form, Heading, BottomLink } from "../styled";
  import { defaultFieldState, isFieldEmpty, isEmailInvalid } from "../utils";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let email = $state(defaultFieldState());

  async function onsubmit() {
    if (isFieldEmpty(email) || isEmailInvalid(email)) {
      return;
    }
    await signUpInit(email.value);
    mail = true;
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
