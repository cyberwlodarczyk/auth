<script lang="ts">
  import {
    createUser,
    store,
    arePasswordFieldsDifferent,
    newFieldState,
    isFieldEmpty,
    isPasswordFieldInvalid,
  } from "$lib/data";
  import { Button, Field, Form, Heading } from "$lib/components";
  import { Route } from "$lib/config";

  interface Props {
    token: string;
  }

  let { token }: Props = $props();

  let name = $state(newFieldState());
  let password = $state(newFieldState());
  let confirmPassword = $state(newFieldState());

  const onsubmit = async () => {
    if (
      isFieldEmpty(name) ||
      isFieldEmpty(password) ||
      isPasswordFieldInvalid(password) ||
      isFieldEmpty(confirmPassword) ||
      arePasswordFieldsDifferent(password, confirmPassword)
    ) {
      return;
    }
    await createUser({ token, name: name.value, password: password.value });
    store.location = Route.Home;
  };
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign up</Heading>
    <Field id="name" label="Name" bind:value={name.value} error={name.error} />
    <Field
      id="password"
      label="Password"
      bind:value={password.value}
      bind:error={password.error}
      type="password"
    />
    <Field
      id="confirm-password"
      label="Confirm password"
      bind:value={confirmPassword.value}
      bind:error={confirmPassword.error}
      type="password"
    />
    <Button submit>Sign up</Button>
  </Form>
</main>
