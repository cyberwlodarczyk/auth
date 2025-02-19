<script lang="ts">
  import { resetPasswordFinish } from "../api";
  import { store } from "../store.svelte";
  import { Button, Field, Form, Heading } from "../styled";
  import {
    defaultFieldState,
    isFieldEmpty,
    isPasswordInvalid,
    arePasswordsDifferent,
  } from "../utils";

  interface Props {
    token: string;
  }

  let { token }: Props = $props();

  let newPassword = $state(defaultFieldState());
  let confirmNewPassword = $state(defaultFieldState());

  async function onsubmit() {
    if (
      isFieldEmpty(newPassword) ||
      isPasswordInvalid(newPassword) ||
      isFieldEmpty(confirmNewPassword) ||
      arePasswordsDifferent(newPassword, confirmNewPassword)
    ) {
      return;
    }
    await resetPasswordFinish(token, newPassword.value);
    store.location = "/";
  }
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
