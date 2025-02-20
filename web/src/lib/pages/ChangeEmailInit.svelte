<script lang="ts">
  import { signUpInit } from "../api";
  import { Button, Field, Form, Heading } from "../styled";
  import { defaultFieldState, isEmailInvalid, isFieldEmpty } from "../utils";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let newEmail = $state(defaultFieldState());

  async function onsubmit() {
    if (isFieldEmpty(newEmail) || isEmailInvalid(newEmail)) {
      return;
    }
    await signUpInit(newEmail.value);
    mail = true;
  }
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
