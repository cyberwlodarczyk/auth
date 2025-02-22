<script lang="ts">
  import {
    newFieldState,
    isFieldEmpty,
    isPasswordFieldInvalid,
    arePasswordFieldsDifferent,
    editUserPassword,
    store,
  } from "$lib/data";
  import { Form, Heading, Field, Button } from "$lib/components";
  import { Route } from "$lib/config";

  let oldPassword = $state(newFieldState());
  let newPassword = $state(newFieldState());
  let confirmNewPassword = $state(newFieldState());

  const onsubmit = async () => {
    if (
      isFieldEmpty(oldPassword) ||
      isFieldEmpty(newPassword) ||
      isPasswordFieldInvalid(newPassword) ||
      isFieldEmpty(confirmNewPassword) ||
      arePasswordFieldsDifferent(newPassword, confirmNewPassword)
    ) {
      return;
    }
    await editUserPassword({
      password: oldPassword.value,
      newPassword: newPassword.value,
    });
    store.location = Route.Home;
  };
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
