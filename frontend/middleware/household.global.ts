export default defineNuxtRouteMiddleware(to => {
  if (
    to.path === "/reports/label-generator" ||
    /^\/(templates?|maintenance)(\/|$)/.test(to.path) ||
    /^\/collection\/(notifiers|entity-types)(\/|$)/.test(to.path) ||
    /^\/item\/[^/]+\/maintenance$/.test(to.path)
  ) {
    return navigateTo("/home");
  }
});
