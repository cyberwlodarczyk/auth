<script lang="ts">
  import { enterSudoMode } from "../api";
  import { Form, Heading, Button, Field } from "../styled";
  import { defaultFieldState, isFieldEmpty } from "../utils";

  interface Props {
    mail: boolean;
  }

  let { mail = $bindable() }: Props = $props();

  let password = $state(defaultFieldState());

  async function onsubmit() {
    if (isFieldEmpty(password)) {
      return;
    }
    await enterSudoMode(password.value);
    mail = true;
  }
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
