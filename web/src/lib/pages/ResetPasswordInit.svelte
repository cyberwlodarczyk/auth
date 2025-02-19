<script lang="ts">
  import { resetPasswordInit } from "../api";
  import { store } from "../store.svelte";
  import { Button, Field, Form, Heading, SecondaryAction } from "../styled";
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
    await resetPasswordInit(email.value);
    mail = true;
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
  <SecondaryAction href="/sign-in">Back to sign in</SecondaryAction>
</main>
