<script lang="ts">
  import { Button, Field, Form, Heading } from "./lib";

  let email = $state("");
  let password = $state("");

  function onsubmit() {
    console.log(email, password);
  }
</script>

<main>
  <Form {onsubmit}>
    <Heading>Log In</Heading>
    <Field
      id="email"
      label="Email"
      bind:value={email}
      type="email"
      error="Invalid email"
    />
    <Field
      id="password"
      label="Password"
      bind:value={password}
      type="password"
    />
    <Button submit>Log In</Button>
  </Form>
</main>
