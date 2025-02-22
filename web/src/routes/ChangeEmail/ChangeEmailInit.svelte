<script lang="ts">
  import {
    createUserConfirmationToken,
    newFieldState,
    isEmailFieldInvalid,
    isFieldEmpty,
  } from "$lib/data";
  import { Button, Field, Form, Heading } from "$lib/components";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let newEmail = $state(newFieldState());

  const onsubmit = async () => {
    if (isFieldEmpty(newEmail) || isEmailFieldInvalid(newEmail)) {
      return;
    }
    await createUserConfirmationToken({ email: newEmail.value });
    mail = true;
  };
</script>

<main>
  <Form {onsubmit}>
    <Heading>Change email</Heading>
    <Field
      id="new-email"
      label="New email"
      bind:value={newEmail.value}
      bind:error={newEmail.error}
      type="email"
    />
    <Button submit>Change email</Button>
  </Form>
</main>
