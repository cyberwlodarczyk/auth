<script lang="ts">
  import { createUserSudoToken, newFieldState, isFieldEmpty } from "$lib/data";
  import { Form, Heading, Button, Field } from "$lib/components";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let password = $state(newFieldState());

  const onsubmit = async () => {
    if (isFieldEmpty(password)) {
      return;
    }
    await createUserSudoToken({ password: password.value });
    mail = true;
  };
</script>

<main>
  <Form {onsubmit}>
    <Heading>Enter sudo mode</Heading>
    <Field
      id="password"
      label="Password"
      bind:value={password.value}
      bind:error={password.error}
      type="password"
    />
    <Button submit>Enter sudo mode</Button>
  </Form>
</main>
