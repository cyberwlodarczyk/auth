<script lang="ts">
  import type { Token } from "$lib/data";

  interface Props {
    token: Token;
  }

  let { token }: Props = $props();

  const formatTimeDigits = (n: number) => {
    return n < 10 ? `0${n}` : n.toString();
  };

  const formatTimeLeft = (token: Token) => {
    const time = token.expiresAt - Math.floor(Date.now() / 1000);
    const minutes = Math.floor(time / 60);
    const seconds = time - minutes * 60;
    return `${formatTimeDigits(minutes)}:${formatTimeDigits(seconds)}`;
  };

  let timeLeft = $state(formatTimeLeft(token));

  $effect(() => {
    const id = window.setInterval(() => {
      timeLeft = formatTimeLeft(token);
    }, 1000);
    return () => {
      window.clearInterval(id);
    };
  });
</script>

<p role="alert">SUDO {timeLeft}</p>

<style>
  p {
    color: var(--error);
    background-color: var(--error-transparent);
    border-radius: var(--border-radius);
    padding: 0.375rem 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    letter-spacing: 0.03125rem;
  }
</style>
