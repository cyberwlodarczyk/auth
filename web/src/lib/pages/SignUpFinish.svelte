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
    name: "",
    password: "",
    confirmPassword: "",
  });

  function onsubmit() {}
</script>

<svelte:head>
  <title>Sign up</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign up</Heading>
    <Field id="name" label="Name" bind:value={state.name} />
    <Field
      id="password"
      label="Password"
      bind:value={state.password}
      type="password"
    />
    <Field
      id="confirm-password"
      label="Confirm password"
      bind:value={state.confirmPassword}
      type="password"
    />
    <Button submit>Sign up</Button>
  </Form>
</main>
