<script lang="ts">
  import {
    store,
    resetUserPassword,
    newFieldState,
    isFieldEmpty,
    isPasswordFieldInvalid,
    arePasswordFieldsDifferent,
  } from "$lib/data";
  import { Button, Field, Form, Heading } from "$lib/components";
  import { Route } from "$lib/config";

  interface Props {
    token: string;
  }

  let { token }: Props = $props();

  let newPassword = $state(newFieldState());
  let confirmNewPassword = $state(newFieldState());

  const onsubmit = async () => {
    if (
      isFieldEmpty(newPassword) ||
      isPasswordFieldInvalid(newPassword) ||
      isFieldEmpty(confirmNewPassword) ||
      arePasswordFieldsDifferent(newPassword, confirmNewPassword)
    ) {
      return;
    }
    await resetUserPassword({ token, password: newPassword.value });
    store.location = Route.Home;
  };
</script>

<svelte:head>
  <title>Reset password</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Reset password</Heading>
    <Field
      id="new-password"
      label="New password"
      bind:value={newPassword.value}
      bind:error={newPassword.error}
      type="password"
    />
    <Field
      id="confirm-new-password"
      label="Confirm new password"
      bind:value={confirmNewPassword.value}
      bind:error={confirmNewPassword.error}
      type="password"
    />
    <Button submit>Reset password</Button>
  </Form>
</main>
