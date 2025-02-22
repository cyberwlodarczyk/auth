<script lang="ts">
  import {
    createUserSessionToken,
    newFieldState,
    isFieldEmpty,
    isEmailFieldInvalid,
  } from "$lib/data";
  import {
    Button,
    Field,
    Form,
    Heading,
    Link,
    SecondaryAction,
  } from "$lib/components";
  import { Route } from "$lib/config";

  let email = $state(newFieldState());
  let password = $state(newFieldState());

  const onsubmit = async () => {
    if (
      isFieldEmpty(email) ||
      isEmailFieldInvalid(email) ||
      isFieldEmpty(password)
    ) {
      return;
    }
    await createUserSessionToken({
      email: email.value,
      password: password.value,
    });
  };
</script>

<svelte:head>
  <title>Sign in</title>
</svelte:head>

<main>
  <Form {onsubmit}>
    <Heading>Sign in</Heading>
    <Field
      id="email"
      label="Email"
      bind:value={email.value}
      bind:error={email.error}
      type="email"
    />
    <Field
      id="password"
      label="Password"
      bind:value={password.value}
      bind:error={password.error}
      type="password"
    />
    <div style:width="100%">
      <Link href={Route.ResetPassword}>Forgot password?</Link>
    </div>
    <Button submit>Sign in</Button>
  </Form>
  <SecondaryAction description="Don't have an account yet?" href={Route.SignUp}>
    Sign up
  </SecondaryAction>
</main>
