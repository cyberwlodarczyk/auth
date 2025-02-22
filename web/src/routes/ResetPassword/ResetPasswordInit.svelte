<script lang="ts">
  import {
    createUserPasswordResetToken,
    newFieldState,
    isFieldEmpty,
    isEmailFieldInvalid,
  } from "$lib/data";
  import {
    Button,
    Field,
    Form,
    Heading,
    SecondaryAction,
  } from "$lib/components";
  import { Route } from "$lib/config";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let email = $state(newFieldState());

  const onsubmit = async () => {
    if (isFieldEmpty(email) || isEmailFieldInvalid(email)) {
      return;
    }
    await createUserPasswordResetToken({ email: email.value });
    mail = true;
  };
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
  <SecondaryAction href={Route.SignIn}>Back to sign in</SecondaryAction>
</main>
