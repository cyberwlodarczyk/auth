<script lang="ts">
  import { Button, Field, Form, Heading } from "../styled";

  interface Props {
    token: string;
  }

  let { token }: Props = $props();

  $effect(() => {
    console.log(token);
  });

  const state = $state({
    newPassword: "",
    confirmNewPassword: "",
  });

  function onsubmit() {}
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
      bind:value={state.newPassword}
      type="password"
    />
    <Field
      id="confirm-new-password"
      label="Confirm new password"
      bind:value={state.confirmNewPassword}
      type="password"
    />
    <Button submit>Reset password</Button>
  </Form>
</main>
