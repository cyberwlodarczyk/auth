<script lang="ts">
  import { Button, Form, FormField, Heading, Input, Label } from "./lib";

  let email = $state("");
  let password = $state("");

  function onsubmit() {
    console.log(email, password);
  }
</script>

<main>
  <Form {onsubmit}>
    <Heading>Log In</Heading>
    <FormField>
      <Label for="email">Email</Label>
      <Input id="email" type="email" autofocus bind:value={email} />
    </FormField>
    <FormField>
      <Label for="password">Password</Label>
      <Input id="password" type="password" bind:value={password} />
    </FormField>
    <Button type="submit">Log In</Button>
  </Form>
</main>
