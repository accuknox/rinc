document.addEventListener("DOMContentLoaded", () => {
  if (!window.location.hash) return;
  const hash = window.location.hash.substring(1);
  const targetElement = document.getElementById(hash);
  if (targetElement) {
    targetElement.scrollIntoView({ behavior: "smooth" });
  }
});
