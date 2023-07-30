import { ref, computed } from "vue";
import { defineStore } from "pinia";

export const useSessionStore = defineStore("session", () => {
  const epassport = ref("");
  const navbutton = ref("search");
  const fileupdateid = ref("");
  const fileoption = ref("internal");
  const typeupdate = ref("");
  const fileoptionprotect = ref("protect");
  const typeupdateprotect = ref("");
  const searchfileoption = ref(["internal"]);
  const filename = ref("");
  const filenameupdate = ref("");
  const filenamedelete = ref("");
  const idfiledelete = ref("");
  const tags = ref([]);
  const tagsnew = ref([]);
  const tagsr = ref([]);
  const tagdocupdate = ref([]);
  const tagnameadd = ref("");
  const tagdelete = ref("");
  const keywords = ref("");
  const usernameadd = ref("");
  const useraddoption = ref([]);
  const userslist = ref([]);
  const usernamedelete = ref("");
  const usernameupdate = ref("");
  const userupdateeditoption = ref(false);
  const userupdateadminoption = ref(false);
  const searchlist = ref([]);
  const isSeacrhLogin = computed(() => {
    if (loginx.value) {
      return true;
    } else {
      searchfileoption.value = ["public"];
      return false;
    }
  });

  const infotxt = ref("");
  const loginx = ref(false);
  const epassword = ref("");

  const canedit = ref(false);
  const canadmin = ref(false);

  return {
    epassport,
    navbutton,
    fileupdateid,
    canedit,
    canadmin,
    filename,
    filenameupdate,
    filenamedelete,
    idfiledelete,
    tags,
    tagsnew,
    tagsr,
    tagdocupdate,
    tagnameadd,
    keywords,
    searchlist,
    infotxt,
    loginx,
    epassword,
    fileoption,
    typeupdate,
    fileoptionprotect,
    typeupdateprotect,
    searchfileoption,
    isSeacrhLogin,
    tagdelete,
    usernameadd,
    useraddoption,
    userslist,
    usernamedelete,
    usernameupdate,
    userupdateeditoption,
    userupdateadminoption,
  };
});
