<script lang="ts">
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

  $effect(() => {
    console.log(token);
  });

  let newPassword = $state(defaultFieldState());
  let confirmNewPassword = $state(defaultFieldState());

  function onsubmit() {
    if (
      isFieldEmpty(newPassword) ||
      isPasswordInvalid(newPassword) ||
      isFieldEmpty(confirmNewPassword) ||
      arePasswordsDifferent(newPassword, confirmNewPassword)
    ) {
      return;
    }
    console.log(newPassword.value, confirmNewPassword.value);
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
