<script lang="ts">
  import { signUpFinish } from "../api";
  import { navigate } from "../store.svelte";
  import { Button, Field, Form, Heading } from "../styled";
  import {
    arePasswordsDifferent,
    defaultFieldState,
    isFieldEmpty,
    isPasswordInvalid,
  } from "../utils";

  interface Props {
    token: string;
  }

  let { token }: Props = $props();

  let name = $state(defaultFieldState());
  let password = $state(defaultFieldState());
  let confirmPassword = $state(defaultFieldState());

  async function onsubmit() {
    if (
      isFieldEmpty(name) ||
      isFieldEmpty(password) ||
      isPasswordInvalid(password) ||
      isFieldEmpty(confirmPassword) ||
      arePasswordsDifferent(password, confirmPassword)
    ) {
      return;
    }
    await signUpFinish(token, name.value, password.value);
    navigate("/");
  }
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign up</Heading>
    <Field id="name" label="Name" bind:value={name.value} error={name.error} />
    <Field
      id="password"
      label="Password"
      bind:value={password.value}
      bind:error={password.error}
      type="password"
    />
    <Field
      id="confirm-password"
      label="Confirm password"
      bind:value={confirmPassword.value}
      bind:error={confirmPassword.error}
      type="password"
    />
    <Button submit>Sign up</Button>
  </Form>
</main>
