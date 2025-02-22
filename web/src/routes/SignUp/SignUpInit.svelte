<script lang="ts">
  import {
    createUserConfirmationToken,
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
    await createUserConfirmationToken({ email: email.value });
    mail = true;
  };
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
  <SecondaryAction description="Already have an account?" href={Route.SignIn}>
    Sign in
  </SecondaryAction>
</main>
