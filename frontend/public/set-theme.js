try {
  const preferences = JSON.parse(localStorage.getItem("homebox/preferences/location") || "null");
  const theme = preferences?.theme || "cupcake";
  document.documentElement.setAttribute("data-theme", theme);
  document.documentElement.classList.add("theme-" + theme);
} catch {
  document.documentElement.setAttribute("data-theme", "cupcake");
  document.documentElement.classList.add("theme-cupcake");
}
