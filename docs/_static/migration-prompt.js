// Deployment choices are homepage content; copying remains owned by Oatbase.
const promptCard = document.querySelector('[data-agent-prompt]');
if (promptCard) {
  const picker = promptCard.querySelector('[data-prompt-target]');
  const source = promptCard.querySelector('[data-copy-source]');
  const alternative = promptCard.querySelector('[data-prompt-alternative]');
  const prompts = {
    github: source.textContent,
    cloudflare: alternative.querySelector('[data-cloudflare-prompt]').textContent
  };
  picker.addEventListener('change', () => { source.textContent = prompts[picker.value]; });
  picker.hidden = false;
  promptCard.querySelector('[data-prompt-heading]').hidden = true;
  alternative.hidden = true;
}
