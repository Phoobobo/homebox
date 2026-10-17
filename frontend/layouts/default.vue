<template>
  <div id="app">
    <!--
    Confirmation Modal is a singleton used by all components so we render
    it here to ensure it's always available. Possibly could move this further
    up the tree
    -->
    <ModalConfirm />
    <OutdatedModal v-if="status" :status="status" />
    <EntityCreateModal />
    <TagCreateModal />
    <AppQuickMenuModal :actions="quickMenuActions" />
    <CollectionCreateModal />
    <CollectionJoinModal />
    <CollectionInviteCreateModal />
    <SidebarProvider :default-open="sidebarState">
      <Sidebar collapsible="icon">
        <SidebarHeader class="gap-4 px-4 py-6">
          <NuxtLink class="flex items-center gap-3 group-data-[collapsible=icon]:hidden" to="/home">
            <AppLogo class="size-9 text-primary" />
            <span class="text-xl font-semibold tracking-tight">{{ $t("app.name") }}</span>
          </NuxtLink>
          <p class="text-xs text-muted-foreground group-data-[collapsible=icon]:hidden">
            {{ $t("global.welcome", { username }) }}
          </p>

          <CollectionSelector />

          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <SidebarMenuButton
                class="flex justify-center bg-primary text-primary-foreground hover:bg-primary/90 active:bg-primary/90 active:text-primary-foreground group-data-[collapsible=icon]:justify-start"
                :tooltip="$t('global.create')"
                hotkey="Shortcut: Ctrl+`"
              >
                <MdiPlus />
                <span>
                  {{ $t("global.create") }}
                </span>
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent class="z-40 min-w-[var(--reka-dropdown-menu-trigger-width)]">
              <DropdownMenuItem
                v-for="btn in dropdown"
                :key="btn.id"
                class="group cursor-pointer text-lg"
                @click="
                  () => {
                    if (btn.dialogId === DialogID.CreateEntity) {
                      if (btn.id == 0)
                        // create item
                        openDialog(btn.dialogId, {
                          params: { baseType: 'item' },
                        });
                      else if (btn.id == 1)
                        // create location
                        openDialog(btn.dialogId, {
                          params: { baseType: 'location' },
                        });
                    } else {
                      openDialog(btn.dialogId as NoParamDialogIDs);
                    }
                  }
                "
              >
                {{ btn.name.value }}
                <Shortcut
                  v-if="btn.shortcut"
                  class="invisible ml-auto group-hover:visible"
                  :keys="btn.shortcut.replace('Shift', '⇧').split('+')"
                />
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarMenu>
              <template v-for="n in nav" :key="n.id">
                <SidebarMenuItem v-if="!n.collapsible" :key="n.id">
                  <SidebarMenuLink
                    :href="n.to"
                    :class="{
                      'bg-accent text-accent-foreground': n.active?.value,
                      'text-nowrap': typeof locale === 'string' && locale.startsWith('zh-'),
                    }"
                    :tooltip="n.name.value"
                  >
                    <component :is="n.icon" />
                    <span>{{ n.name.value }}</span>
                  </SidebarMenuLink>
                </SidebarMenuItem>

                <Collapsible v-else default-open class="group/collapsible">
                  <SidebarMenuItem>
                    <SidebarMenuItem class="flex gap-1">
                      <SidebarMenuLink
                        :href="n.to"
                        :class="{
                          'bg-accent text-accent-foreground': n.active?.value,
                          'text-nowrap': typeof locale === 'string' && locale.startsWith('zh-'),
                        }"
                        :tooltip="n.name.value"
                      >
                        <component :is="n.icon" />
                        <span>{{ n.name.value }}</span>
                      </SidebarMenuLink>
                      <CollapsibleTrigger as-child>
                        <SidebarMenuButton class="flex size-12 items-center justify-center">
                          <MdiChevronRight
                            class="transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
                          />
                        </SidebarMenuButton>
                      </CollapsibleTrigger>
                    </SidebarMenuItem>
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        <SidebarMenuSubItem v-for="c in n.collapsible" :key="c.id">
                          <SidebarMenuLink
                            :href="c.to"
                            :class="{
                              'bg-accent text-accent-foreground': c.active?.value,
                              'text-nowrap': typeof locale === 'string' && locale.startsWith('zh-'),
                              'h-min py-0': true,
                            }"
                            :tooltip="c.name.value"
                          >
                            <span>{{ c.name.value }}</span>
                          </SidebarMenuLink>
                        </SidebarMenuSubItem>
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </SidebarMenuItem>
                </Collapsible>
              </template>
            </SidebarMenu>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <SidebarMenuButton
            class="flex justify-center group-data-[collapsible=icon]:justify-start group-data-[collapsible=icon]:bg-destructive group-data-[collapsible=icon]:text-destructive-foreground group-data-[collapsible=icon]:shadow-sm group-data-[collapsible=icon]:hover:bg-destructive/90"
            :tooltip="$t('global.sign_out')"
            data-testid="logout-button"
            @click="logout"
          >
            <MdiLogout />
            <span>
              {{ $t("global.sign_out") }}
            </span>
          </SidebarMenuButton>
        </SidebarFooter>

        <SidebarRail />
      </Sidebar>
      <SidebarInset class="min-h-dvh max-w-full overflow-hidden bg-background-accent">
        <div class="relative flex h-full flex-col justify-center">
          <div v-if="preferences.displayLegacyHeader">
            <AppHeaderDecor class="-mt-10 hidden lg:block" />
            <SidebarTrigger class="absolute left-2 top-2 hidden lg:flex" variant="default" />
            <LanguageQuickSwitch class="absolute right-2 top-2 z-20 hidden lg:flex" />
          </div>
          <!-- IMPORTANT: if you change the height of this div, alter the top value in the item edit page-->
          <div
            class="zainaer-header sticky top-0 z-20 flex h-[var(--header-height-mobile)] translate-y-[-0.5px] flex-col bg-background p-3 sm:h-[var(--header-height)] sm:flex-row"
            :class="{
              'lg:hidden': preferences.displayLegacyHeader,
            }"
          >
            <div class="flex h-1/2 items-center gap-2 sm:h-auto">
              <SidebarTrigger variant="ghost" />
              <NuxtLink to="/home">
                <AppHeaderText class="h-6" />
              </NuxtLink>
              <LanguageQuickSwitch class="ml-auto" />
            </div>
            <div class="sm:grow" />
            <div class="flex h-1/2 grow items-center justify-end gap-2 sm:h-auto">
              <Input
                v-model:model-value="search"
                class="h-9 grow sm:max-w-sm"
                :placeholder="$t('global.search')"
                type="search"
                @keyup.enter="triggerSearch"
              />
              <div>
                <Button size="icon" :aria-label="$t('global.search')" @click="triggerSearch">
                  <MdiMagnify />
                </Button>
              </div>
            </div>
          </div>

          <slot />
          <div
            class="fixed bottom-6 right-5 z-30 flex flex-col items-end gap-3 sm:bottom-8 sm:right-8"
            style="padding-bottom: env(safe-area-inset-bottom)"
          >
            <TextAssistant />
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button
                  class="h-14 gap-2 rounded-full px-5 shadow-lg ring-4 ring-background"
                  :aria-label="$t('global.create')"
                  data-testid="floating-create-button"
                >
                  <MdiPlus class="size-6" />
                  <span>{{ $t("global.create") }}</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" side="top" :side-offset="12" class="w-52 p-2">
                <DropdownMenuItem
                  class="cursor-pointer gap-3 py-3"
                  @click="openDialog(DialogID.CreateEntity, { params: { baseType: 'item' } })"
                >
                  <MdiPackageVariant class="size-5 text-primary" />
                  {{ $t("menu.create_item") }}
                </DropdownMenuItem>
                <DropdownMenuItem
                  class="cursor-pointer gap-3 py-3"
                  @click="openDialog(DialogID.CreateEntity, { params: { baseType: 'location' } })"
                >
                  <MdiMapMarkerOutline class="size-5 text-primary" />
                  {{ $t("menu.create_location") }}
                </DropdownMenuItem>
                <DropdownMenuItem class="cursor-pointer gap-3 py-3" @click="openDialog(DialogID.CreateTag)">
                  <MdiTagMultiple class="size-5 text-primary" />
                  {{ $t("menu.create_tag") }}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <div class="grow" />

          <footer v-if="status" class="zainaer-footer bottom-0 w-full pb-4 text-center">
            <p class="text-center text-sm">
              {{ $t("global.version", { version: status.build.version }) }} ·
              {{ $t("global.build", { build: status.build.commit.slice(0, 7) }) }}
            </p>
          </footer>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>

<script lang="ts" setup>
  import TextAssistant from "~/components/App/TextAssistant.vue";
  import { useI18n } from "vue-i18n";
  import LanguageQuickSwitch from "~/components/App/LanguageQuickSwitch.vue";
  import { useTagStore } from "~/stores/tags";
  import { useLocationStore } from "~~/stores/locations";
  import { useEntityTypeStore } from "~~/stores/entityTypes";

  import MdiHome from "~icons/mdi/home";
  import MdiFileTree from "~icons/mdi/file-tree";
  import MdiTagMultiple from "~icons/mdi/tag-multiple";
  import MdiMagnify from "~icons/mdi/magnify";
  import MdiAccount from "~icons/mdi/account";
  import MdiCog from "~icons/mdi/cog";
  import MdiPackageVariant from "~icons/mdi/package-variant";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiPlus from "~icons/mdi/plus";
  import MdiLogout from "~icons/mdi/logout";
  import MdiChevronRight from "~icons/mdi/chevron-right";

  import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarHeader,
    SidebarInset,
    SidebarMenu,
    SidebarMenuSub,
    SidebarMenuSubItem,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarMenuLink,
    SidebarProvider,
    SidebarRail,
    SidebarTrigger,
  } from "@/components/ui/sidebar";
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
  } from "@/components/ui/dropdown-menu";
  import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
  import { Shortcut } from "~/components/ui/shortcut";
  import { useDialog } from "~/components/ui/dialog-provider";
  import { Input } from "~/components/ui/input";
  import { Button } from "~/components/ui/button";
  import { DialogID, type NoParamDialogIDs } from "~/components/ui/dialog-provider/utils";
  import ModalConfirm from "~/components/ModalConfirm.vue";
  import OutdatedModal from "~/components/App/OutdatedModal.vue";
  import EntityCreateModal from "~/components/Entity/CreateModal.vue";
  import TagCreateModal from "~/components/Tag/CreateModal.vue";
  import AppQuickMenuModal from "~/components/App/QuickMenuModal.vue";
  import AppLogo from "~/components/App/Logo.vue";
  import AppHeaderDecor from "~/components/App/HeaderDecor.vue";
  import AppHeaderText from "~/components/App/HeaderText.vue";
  import CollectionSelector from "~/components/Collection/Selector.vue";
  import CollectionCreateModal from "~/components/Collection/CreateModal.vue";
  import CollectionJoinModal from "~/components/Collection/JoinModal.vue";
  import CollectionInviteCreateModal from "~/components/Collection/InviteCreateModal.vue";

  const { t, locale } = useI18n();
  const username = computed(() => authCtx.user?.name || "User");

  const { openDialog } = useDialog();

  const preferences = useViewPreferences();

  // get sidebar state from cookies
  const sidebarState = useCookie("sidebar:state", {
    readonly: true,
    decode: value => value !== "false",
  });

  const pubApi = usePublicApi();
  const { data: status } = useAsyncData(async () => {
    const { data } = await pubApi.status();

    return data;
  });

  const search = ref("");

  const triggerSearch = () => {
    if (search.value) {
      navigateTo(`/items?q=${encodeURIComponent(search.value)}`);
      search.value = "";
      // remove focus from input
      if (document.activeElement && "blur" in document.activeElement) {
        (document.activeElement as HTMLElement).blur();
      }
    }
  };

  type DropdownItem = {
    id: number;
    name: ComputedRef<string>;
    shortcut: string;
    dialogId: DialogID;
  };

  const dropdown: DropdownItem[] = [
    {
      id: 0,
      name: computed(() => t("menu.create_item")),
      shortcut: "Shift+1",
      dialogId: DialogID.CreateEntity,
    },
    {
      id: 1,
      name: computed(() => t("menu.create_location")),
      shortcut: "Shift+2",
      dialogId: DialogID.CreateEntity,
    },
    {
      id: 2,
      name: computed(() => t("menu.create_tag")),
      shortcut: "Shift+3",
      dialogId: DialogID.CreateTag,
    },
  ];

  const route = useRoute();
  const router = useRouter();

  const nav: {
    icon: Component;
    active: ComputedRef<boolean>;
    id: number;
    name: ComputedRef<string>;
    to: string;
    collapsible?: {
      active: ComputedRef<boolean>;
      id: number;
      name: ComputedRef<string>;
      to: string;
    }[];
  }[] = [
    {
      icon: MdiHome,
      active: computed(() => route.path === "/home"),
      id: 0,
      name: computed(() => t("menu.home")),
      to: "/home",
    },
    {
      icon: MdiFileTree,
      id: 1,
      active: computed(() => route.path === "/locations"),
      name: computed(() => t("menu.locations")),
      to: "/locations",
    },
    {
      icon: MdiTagMultiple,
      id: 2,
      active: computed(() => route.path === "/tags"),
      name: computed(() => t("global.tags")),
      to: "/tags",
    },
    {
      icon: MdiMagnify,
      id: 3,
      active: computed(() => route.path === "/items"),
      name: computed(() => t("menu.search")),
      to: "/items",
    },
    {
      icon: MdiAccount,
      id: 6,
      active: computed(() => route.path === "/profile"),
      name: computed(() => t("menu.profile")),
      to: "/profile",
    },
    {
      icon: MdiCog,
      id: 7,
      active: computed(() => route.path.includes("/collection")),
      name: computed(() => t("menu.collection")),
      to: "/collection/members",
      collapsible: [
        {
          id: 61,
          active: computed(() => route.path === "/collection/members"),
          name: computed(() => t("collection.tabs.members")),
          to: "/collection/members",
        },
        {
          id: 62,
          active: computed(() => route.path === "/collection/invites"),
          name: computed(() => t("collection.tabs.invites")),
          to: "/collection/invites",
        },
        {
          id: 64,
          active: computed(() => route.path === "/collection/settings"),
          name: computed(() => t("collection.tabs.settings")),
          to: "/collection/settings",
        },
        {
          id: 67,
          active: computed(() => route.path === "/collection/models"),
          name: computed(() => t("model_settings")),
          to: "/collection/models",
        },
        {
          id: 66,
          active: computed(() => route.path === "/collection/tools"),
          name: computed(() => t("collection.tabs.tools")),
          to: "/collection/tools",
        },
      ],
    },
  ];

  const quickMenuActions = reactive([
    ...dropdown.map(v => ({
      text: computed(() => v.name.value),
      dialogId: v.dialogId,
      shortcut: v.shortcut.split("+")[1] as string,
      id: v.id,
      type: "create" as const,
    })),
    ...nav.map(v => ({
      text: computed(() => v.name.value),
      href: v.to,
      type: "navigate" as const,
    })),
  ]);

  const tagStore = useTagStore();
  tagStore.ensureAllTagsFetched();

  const locationStore = useLocationStore();
  locationStore.ensureLocationsFetched();

  const entityTypeStore = useEntityTypeStore();
  entityTypeStore.ensureFetched();

  onMounted(() => {
    locationStore.refreshParents();
    locationStore.refreshTree();

    // Auto-open JoinModal when invitation token is in URL
    const token = route.query.token;
    if (typeof token === "string" && token.length > 0) {
      // Remove token from browser URL
      const url = new URL(window.location.href);
      url.searchParams.delete("token");
      window.history.replaceState(history.state, "", url.toString());

      // Sync router's state to clear route.query.token
      const { token: _, ...cleanQuery } = route.query;
      router.replace({ query: cleanQuery });

      openDialog(DialogID.JoinCollection, {
        params: { inviteCode: token },
      });
    }
  });

  onServerEvent(ServerEvent.TagMutation, () => {
    tagStore.refresh();
  });

  onServerEvent(ServerEvent.EntityMutation, () => {
    locationStore.refreshChildren();
    locationStore.refreshParents();
    locationStore.refreshTree();
  });

  const authCtx = useAuthContext();
  const api = useUserApi();

  async function logout() {
    await authCtx.logout(api);
    navigateTo("/");
  }
</script>
