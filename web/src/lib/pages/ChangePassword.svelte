<script lang="ts">
  import {
    defaultFieldState,
    isFieldEmpty,
    isPasswordInvalid,
    arePasswordsDifferent,
  } from "../utils";
  import { Form, Heading, Field, Button } from "../styled";
  import { changePassword } from "../api";
  import { store } from "../store.svelte";

  let oldPassword = $state(defaultFieldState());
  let newPassword = $state(defaultFieldState());
  let confirmNewPassword = $state(defaultFieldState());

  async function onsubmit() {
    if (
      isFieldEmpty(oldPassword) ||
      isFieldEmpty(newPassword) ||
      isPasswordInvalid(newPassword) ||
      isFieldEmpty(confirmNewPassword) ||
      arePasswordsDifferent(newPassword, confirmNewPassword)
    ) {
      return;
    }
    await changePassword(oldPassword.value, newPassword.value);
    store.location = "/";
  }
</script>

<svelte:head>
  <title>Change password</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Change password</Heading>
    <Field
      id="old-password"
      label="Old password"
      bind:value={oldPassword.value}
      bind:error={oldPassword.error}
      type="password"
    />
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
    <Button submit>Change password</Button>
  </Form>
</main>
