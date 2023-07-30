<script setup lang="ts">
import { useSessionStore } from "@/stores/session";
import { computed } from "vue";

const sess = useSessionStore();
var fileblob: any;
declare const bootstrap: any;

const loginepassport = () => {
  sess.epassport = sess.epassport.toLowerCase();
  const dataout = JSON.stringify({
    epassport: sess.epassport,
    epassword: sess.epassword,
  });
  fetch("/api/login", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: dataout,
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.status == "ok") {
        if (data.canview == 1) {
          sess.loginx = true;
          if (data.canedit == 1) {
            sess.canedit = true;
          } else {
            sess.canedit = false;
          }
          if (data.canadmin == 1) {
            sess.canadmin = true;
          } else {
            sess.canadmin = false;
          }
        } else {
          sess.epassport = "";
          sess.epassword = "";
          sess.infotxt = "You have no permission";
          new bootstrap.Modal("#infomodal", {
            keyboard: false,
          }).show();
        }
      } else {
        sess.epassport = "";
        sess.epassword = "";
        sess.infotxt = "Login fail";
        new bootstrap.Modal("#infomodal", {
          keyboard: false,
        }).show();
      }
    })
    .catch((err) => {
      sess.epassport = "";
      sess.epassword = "";
      sess.infotxt = "Error! Something Wrong!!!";
      new bootstrap.Modal("#infomodal", {
        keyboard: false,
      }).show();
    });
};

const addtag = (tag: never) => {
  if (sess.tagsnew.indexOf(tag) == -1) {
    sess.tagsnew.push(tag);
  }
};

const addtagsr = (tag: never) => {
  if (sess.tagsr.indexOf(tag) == -1) {
    sess.tagsr.push(tag);
    searchxdoc();
  }
};

const addtagdocupdate = (tag: never) => {
  if (sess.tagdocupdate.indexOf(tag) == -1) {
    sess.tagdocupdate.push(tag);
  }
};

const removetag = (tag: never) => {
  var i = sess.tagsnew.indexOf(tag);
  if (i != -1) {
    sess.tagsnew.splice(i, 1);
  }
};

const removetagsr = (tag: never) => {
  var i = sess.tagsr.indexOf(tag);
  if (i != -1) {
    sess.tagsr.splice(i, 1);
    searchxdoc();
  }
};

const removetagdocupdate = (tag: never) => {
  var i = sess.tagdocupdate.indexOf(tag);
  if (i != -1) {
    sess.tagdocupdate.splice(i, 1);
  }
};

const addnewtag = () => {
  sess.tagnameadd = sess.tagnameadd.trim();
  const dataload = {
    tag: sess.tagnameadd,
  };
  fetch("/api/addnewtag", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      sess.tagnameadd = "";
      getTags();
    });
};

const preparetagdelete = (tag: never) => {
  sess.tagdelete = tag;
};

const deletetag = () => {
  const dataload = {
    tag: sess.tagdelete,
  };
  fetch("/api/deletetag", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      getTags();
    });
};

const preparedeletedoc = (id: any, filename: any) => {
  sess.idfiledelete = id;
  sess.filenamedelete = filename;
};

const deletefile = () => {
  const dataload = {
    id: sess.idfiledelete,
    username: sess.epassport,
  };
  fetch("/api/deletefile", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.status == "true") {
        searchxdoc();
        alert("Delete file successful");
      } else {
        alert("Delete file fail");
      }
    });
};

const logoutepassport = () => {
  sess.epassport = "";
  sess.epassword = "";
  sess.loginx = false;
  sess.canedit = false;
  sess.canadmin = false;
  sess.tagsnew = [];
  sess.searchfileoption = ["public"];
  sess.navbutton = "search";
  searchxdoc();
};

const filech = (event: any) => {
  fileblob = event.target.files[0];
};

const filechupdate = (event: any) => {
  fileblob = event.target.files[0];
};

const uploadFile = () => {
  if (document.getElementById("formFile").files.length == 0) {
    alert("การบันทึกไฟล์ไม่สำเร็จ เนื่องจากไม่ได้ระบุไฟล์");
    return;
  }
  if (sess.filename.trim() == "") {
    alert("การบันทึกไฟล์ไม่สำเร็จ กรุณาระบุชื่อไฟล์");
    return;
  }
  var reader = new FileReader();
  reader.readAsDataURL(fileblob);
  reader.onload = () => {
    var fupload = {
      filename: sess.filename,
      tags: sess.tagsnew,
      data: reader.result,
      username: sess.epassport,
      type: sess.fileoption,
      protect: sess.fileoptionprotect,
    };
    fetch("/api/fileupload", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(fupload),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "true") {
          sess.fileoption = "internal";
          sess.fileoptionprotect = "unprotect";
          sess.tagsnew = [];
          sess.filename = "";
          document.getElementById("formFile").value = null;
          alert("การบันทึกสำเร็จ");
          searchxdoc();
        } else {
          alert("การบันทึกมีปัญหาบางประการ");
        }
      })
      .catch((err) => {
        alert("การบันทึกไม่สำเร็จ");
      });
  };
};

const searchxdoc = () => {
  var searchinput = {
    username: sess.epassport,
    keywords: sess.keywords,
    tags: sess.tagsr,
    type: sess.searchfileoption,
  };
  fetch("/api/search", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(searchinput),
  })
    .then((res) => res.json())
    .then((data) => {
      sess.searchlist = data.results;
    });
};

const downloaded = (data: any, filename: any) => {
  const link = document.createElement("a");

  link.setAttribute("href", data);
  link.setAttribute("download", filename);
  link.style.display = "none";

  document.body.appendChild(link);

  link.click();

  document.body.removeChild(link);
};

const downloadfile = (id: any) => {
  var dataload = {
    id: id,
    username: sess.epassport,
  };
  fetch("/api/download", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((data) => data.text())
    .then((data) => {
      downloaded(data, dataload.id);
    });
};

const getTags = () => {
  fetch("/api/tags", {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  })
    .then((res) => res.json())
    .then((data) => {
      sess.tags = data.tags;
      searchxdoc();
    });
};

const getUsers = () => {
  fetch("/api/users", {
    method: "GET",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  })
    .then((res) => res.json())
    .then((data) => {
      sess.userslist = data.users;
    });
};

const updatedoc = () => {
  updateDocMeta();
  updateDocFile();
};

const updateDocMeta = () => {
  if (sess.filenameupdate.trim() == "") {
    return;
  }
  var docmeta = {
    id: sess.fileupdateid,
    username: sess.epassport,
    filename: sess.filenameupdate,
    tags: sess.tagdocupdate,
    type: sess.typeupdate,
    protect: sess.typeupdateprotect,
  };
  fetch("/api/updatemetafile", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(docmeta),
  })
    .then((res) => res.json())
    .then((data) => {
      searchxdoc();
    });
};

const updateDocFile = () => {
  if (document.getElementById("formFileupdate").files.length == 0) {
    return;
  }
  if (sess.filenameupdate.trim() == "") {
    return;
  }
  var reader = new FileReader();
  reader.readAsDataURL(fileblob);
  reader.onload = () => {
    var fupload = {
      id: sess.fileupdateid,
      data: reader.result,
      username: sess.epassport,
    };
    fetch("/api/updatefileupload", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(fupload),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "true") {
          document.getElementById("formFileupdate").value = null;
          alert("การบันทึกสำเร็จ");
        } else {
          alert("การบันทึกมีปัญหาบางประการ");
        }
      })
      .catch((err) => {
        alert("การบันทึกไม่สำเร็จ");
      });
  };
};

const prepareupdatedoc = (
  id: any,
  filename: any,
  tags: any,
  type: any,
  protect: any
) => {
  sess.fileupdateid = id;
  sess.filenameupdate = filename;
  sess.tagdocupdate = tags;
  sess.typeupdate = type;
  sess.typeupdateprotect = protect;
  document.getElementById("formFileupdate").value = null;
};

const adduser = () => {
  var dataload = {
    username: sess.usernameadd,
    permission: sess.useraddoption,
  };
  fetch("/api/adduser", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      sess.usernameadd = "";
      getUsers();
    });
};

const preparedeleteuser = (username: any) => {
  sess.usernamedelete = username;
};

const deleteuser = () => {
  var dataload = {
    username: sess.usernamedelete,
  };
  fetch("/api/deluser", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      getUsers();
    });
};

const prepareupdateuser = (username: any, edit: any, admin: any) => {
  sess.usernameupdate = username;
  if (edit == 1) {
    sess.userupdateeditoption = true;
  } else {
    sess.userupdateeditoption = false;
  }
  if (admin == 1) {
    sess.userupdateadminoption = true;
  } else {
    sess.userupdateadminoption = false;
  }
};

const updateuser = () => {
  var dataload = {
    username: sess.usernameupdate,
    edit: sess.userupdateeditoption,
    admin: sess.userupdateadminoption,
  };
  fetch("/api/updateuser", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dataload),
  })
    .then((res) => res.json())
    .then((data) => {
      getUsers();
    });
};

const fileqrcode = (id: any) => {
  window.open("/qrcode/" + id, "_blank");
};

getTags();
</script>

<template>
  <div class="container" data-bs-theme="dark">
    <!-- Header Start -->
    <nav class="navbar navbar-expand-lg bg-body-tertiary" data-bs-theme="dark">
      <div class="container-fluid">
        <a class="navbar-brand" href="#">
          <img src="/logo.png" width="45" />&nbsp;&nbsp;&nbsp;&nbsp;xDocument
          System
        </a>
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0"></ul>
          <button
            v-if="!sess.loginx"
            class="btn btn-outline-primary rounded-pill"
            data-bs-toggle="modal"
            data-bs-target="#loginmodal"
          >
            Login
          </button>
          <button
            v-if="sess.loginx"
            class="btn btn-primary rounded-pill"
            data-bs-toggle="modal"
            data-bs-target="#logoutmodal"
          >
            {{ sess.epassport }}
          </button>
        </div>
      </div>
    </nav>
    <!-- Header End -->
    <br />
    <div>
      <!-- Body Start -->
      <div class="card" style="padding: 10px">
        <div class="card-body">
          <div
            class="btn-group"
            role="group"
            aria-label="Basic radio toggle button group"
          >
            <input
              type="radio"
              class="btn-check"
              name="navradio"
              id="searchbutton"
              autocomplete="off"
              value="search"
              v-model="sess.navbutton"
            />
            <label class="btn btn-outline-primary" for="searchbutton"
              >Search Documents</label
            >

            <input
              type="radio"
              class="btn-check"
              name="navradio"
              id="adddocbutton"
              autocomplete="off"
              v-if="sess.canedit"
              value="add"
              v-model="sess.navbutton"
            />
            <label
              class="btn btn-outline-primary"
              for="adddocbutton"
              v-if="sess.canedit"
              >Add Documents</label
            >

            <input
              type="radio"
              class="btn-check"
              name="navradio"
              id="managetagsbutton"
              autocomplete="off"
              v-if="sess.canedit"
              value="tags"
              @click="getTags"
              v-model="sess.navbutton"
            />
            <label
              class="btn btn-outline-primary"
              for="managetagsbutton"
              v-if="sess.canedit"
              >Manage Tags</label
            >
            <input
              type="radio"
              class="btn-check"
              name="navradio"
              id="adminbutton"
              autocomplete="off"
              v-if="sess.canadmin"
              value="admin"
              @click="getUsers"
              v-model="sess.navbutton"
            />
            <label
              class="btn btn-outline-primary"
              for="adminbutton"
              v-if="sess.canadmin"
              >Admin</label
            >
          </div>
          <br /><br />
          <div
            class="card"
            style="padding: 10px"
            v-if="sess.navbutton == 'search'"
          >
            <h3>Search Documents</h3>
            <div class="card-body">
              <div class="form-floating mb-3">
                <input
                  type="text"
                  class="form-control"
                  id="keywords"
                  placeholder="keywords"
                  v-model="sess.keywords"
                  @change="searchxdoc"
                />
                <label for="keywords">Keywords</label>
              </div>
              <div>
                <li
                  v-for="tag in sess.tagsr"
                  class="badge bg-secondary"
                  @click="removetagsr(tag)"
                >
                  {{ tag }}
                </li>
                <br />
                <div class="dropdown">
                  <a
                    class="btn btn-outline-primary dropdown-toggle rounded-pill"
                    href="#"
                    role="button"
                    data-bs-toggle="dropdown"
                    aria-expanded="false"
                  >
                    Tags
                  </a>

                  <ul class="dropdown-menu">
                    <li v-for="tag in sess.tags">
                      <a
                        class="dropdown-item"
                        href="#"
                        @click="addtagsr(tag)"
                        >{{ tag }}</a
                      >
                    </li>
                  </ul>
                </div>
              </div>
              <br />
              <div
                class="btn-group"
                role="group"
                aria-label="Basic radio toggle button group"
              >
                <input
                  type="checkbox"
                  class="btn-check"
                  name="searchfileradio"
                  id="searchfilepublic"
                  autocomplete="off"
                  value="public"
                  v-model="sess.searchfileoption"
                  @change="searchxdoc"
                />
                <label class="btn btn-outline-primary" for="searchfilepublic"
                  >public</label
                >
                <input
                  type="checkbox"
                  class="btn-check"
                  name="searchfileradio"
                  id="searchfileinternal"
                  autocomplete="off"
                  value="internal"
                  v-model="sess.searchfileoption"
                  @change="searchxdoc"
                  v-if="sess.isSeacrhLogin"
                />
                <label
                  class="btn btn-outline-primary"
                  for="searchfileinternal"
                  v-if="sess.isSeacrhLogin"
                  >internal</label
                >
                <input
                  type="checkbox"
                  class="btn-check"
                  name="searchfileradio"
                  id="searchfileprivate"
                  autocomplete="off"
                  value="private"
                  v-model="sess.searchfileoption"
                  @change="searchxdoc"
                  v-if="sess.isSeacrhLogin"
                />
                <label
                  class="btn btn-outline-primary"
                  for="searchfileprivate"
                  v-if="sess.isSeacrhLogin"
                  >private</label
                >
              </div>
              <br />
            </div>
            <hr />
            <div class="card-body">
              <div v-if="sess.searchlist.length == 0">ไม่มีรายการ</div>
              <table
                class="table table-dark"
                v-if="sess.searchlist.length != 0"
              >
                <thead>
                  <tr>
                    <td><h4>Documents</h4></td>
                    <td><h4>Action</h4></td>
                  </tr>
                </thead>
                <tbody class="table-group-divider">
                  <tr v-for="row in sess.searchlist">
                    <td>
                      {{ row.filename }}
                      <br />
                      <span
                        v-if="row.type == 'public'"
                        class="badge text-bg-success"
                        >{{ row.type }}</span
                      >
                      <span
                        v-if="row.type == 'internal'"
                        class="badge text-bg-warning"
                        >{{ row.type }}</span
                      >
                      <span
                        v-if="row.type == 'private'"
                        class="badge text-bg-danger"
                        >{{ row.type }}</span
                      >
                      <span
                        v-for="tag in row.tags"
                        class="badge bg-secondary"
                        >{{ tag }}</span
                      >
                      <span v-if="sess.loginx">
                        <span
                          v-if="row.protect == 'protect'"
                          class="badge text-bg-info"
                          >protect by {{ row.username }}</span
                        >
                      </span>
                      <br />
                    </td>
                    <td>
                      <button
                        class="btn btn-outline-primary rounded-pill"
                        @click="downloadfile(row.id)"
                      >
                        Download</button
                      ><span v-if="row.type == 'public'"
                        >&nbsp;&nbsp;
                        <button
                          class="btn btn-outline-primary rounded-pill"
                          @click="fileqrcode(row.id)"
                        >
                          QRCode
                        </button></span
                      >
                      <span
                        v-if="
                          sess.isSeacrhLogin &&
                          ((row.protect == 'protect' &&
                            row.username == sess.epassport) ||
                            row.protect == 'unprotect')
                        "
                        >&nbsp;&nbsp;
                        <button
                          class="btn btn-outline-primary rounded-pill"
                          @click="
                            prepareupdatedoc(
                              row.id,
                              row.filename,
                              [...row.tags],
                              row.type,
                              row.protect
                            )
                          "
                          data-bs-toggle="modal"
                          data-bs-target="#updatedocmodal"
                        >
                          Update</button
                        >&nbsp;&nbsp;<button
                          class="btn btn-outline-primary rounded-pill"
                          data-bs-toggle="modal"
                          data-bs-target="#confirmdeletefilemodal"
                          @click="preparedeletedoc(row.id, row.filename)"
                        >
                          Delete
                        </button>
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div
            class="card"
            style="padding: 10px"
            v-if="sess.navbutton == 'add'"
          >
            <h3>Add Documents</h3>
            <div class="card-body">
              <div class="form-floating mb-3">
                <input
                  type="text"
                  class="form-control"
                  id="filename"
                  placeholder=" "
                  v-model="sess.filename"
                />
                <label for="filename">File Name</label><br />
                <input
                  class="form-control"
                  type="file"
                  id="formFile"
                  @change="filech($event)"
                /><br />
                <li
                  v-for="tag in sess.tagsnew"
                  class="badge bg-secondary"
                  @click="removetag(tag)"
                >
                  {{ tag }}
                </li>
                <br />
                <div class="dropdown">
                  <a
                    class="btn btn-outline-primary dropdown-toggle rounded-pill"
                    href="#"
                    role="button"
                    data-bs-toggle="dropdown"
                    aria-expanded="false"
                  >
                    Tags
                  </a>

                  <ul class="dropdown-menu">
                    <li v-for="tag in sess.tags">
                      <a class="dropdown-item" href="#" @click="addtag(tag)">{{
                        tag
                      }}</a>
                    </li>
                  </ul>
                </div>
                <br />
                <div
                  class="btn-group"
                  role="group"
                  aria-label="Basic radio toggle button group"
                >
                  <input
                    type="radio"
                    class="btn-check"
                    name="fileradio"
                    id="filepublic"
                    autocomplete="off"
                    value="public"
                    v-model="sess.fileoption"
                  />
                  <label class="btn btn-outline-primary" for="filepublic"
                    >public</label
                  >
                  <input
                    type="radio"
                    class="btn-check"
                    name="fileradio"
                    id="fileinternal"
                    autocomplete="off"
                    value="internal"
                    v-model="sess.fileoption"
                  />
                  <label class="btn btn-outline-primary" for="fileinternal"
                    >internal</label
                  >
                  <input
                    type="radio"
                    class="btn-check"
                    name="fileradio"
                    id="fileprivate"
                    autocomplete="off"
                    value="private"
                    v-model="sess.fileoption"
                  />
                  <label class="btn btn-outline-primary" for="fileprivate"
                    >private</label
                  >
                </div>
                <br v-if="sess.fileoption != 'private'" /><br
                  v-if="sess.fileoption != 'private'"
                />
                <div
                  class="btn-group"
                  role="group"
                  aria-label="Basic radio toggle button group"
                  v-if="sess.fileoption != 'private'"
                >
                  <input
                    type="radio"
                    class="btn-check"
                    name="fileprotectradio"
                    id="fileprotect"
                    autocomplete="off"
                    value="protect"
                    v-model="sess.fileoptionprotect"
                  />
                  <label class="btn btn-outline-primary" for="fileprotect"
                    >protect</label
                  >
                  <input
                    type="radio"
                    class="btn-check"
                    name="fileprotectradio"
                    id="fileunprotect"
                    autocomplete="off"
                    value="unprotect"
                    v-model="sess.fileoptionprotect"
                  />
                  <label class="btn btn-outline-primary" for="fileunprotect"
                    >unprotect</label
                  >
                </div>
                <br /><br />
                <button
                  type="button"
                  class="btn btn-outline-primary rounded-pill"
                  @click="uploadFile"
                >
                  Upload
                </button>
              </div>
            </div>
          </div>
          <div
            class="card"
            style="padding: 10px"
            v-if="sess.navbutton == 'tags'"
          >
            <h3>Manage Tags</h3>
            <div class="card-body">
              <div class="form-floating mb-3">
                <input
                  type="text"
                  class="form-control"
                  id="tagname"
                  placeholder=" "
                  v-model="sess.tagnameadd"
                />
                <label for="tagname">Tag Name</label><br />
                <button
                  class="btn btn-outline-primary rounded-pill"
                  @click="addnewtag"
                >
                  Add</button
                ><br />
              </div>
              <hr />
              <table class="table table-dark">
                <thead>
                  <tr>
                    <td><h4>Tagsname</h4></td>
                    <td><h4>Action</h4></td>
                  </tr>
                </thead>
                <tbody class="table-group-divider">
                  <tr v-for="tag in sess.tags">
                    <td>{{ tag }}</td>
                    <td>
                      <button
                        class="btn btn-outline-primary rounded-pill"
                        @click="preparetagdelete(tag)"
                        data-bs-toggle="modal"
                        data-bs-target="#confirmmodal"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div
            class="card"
            style="padding: 10px"
            v-if="sess.navbutton == 'admin'"
          >
            <h3>Admin</h3>
            <div class="card-body">
              <div class="form-floating mb-3">
                <input
                  type="text"
                  class="form-control"
                  id="usernameadd"
                  placeholder=" "
                  v-model="sess.usernameadd"
                />
                <label for="usernameadd">Username</label><br />
                <div
                  class="btn-group"
                  role="group"
                  aria-label="Basic radio toggle button group"
                >
                  <input
                    type="checkbox"
                    class="btn-check"
                    name="useraddradio"
                    id="useraddradioedit"
                    autocomplete="off"
                    value="edit"
                    v-model="sess.useraddoption"
                  />
                  <label class="btn btn-outline-primary" for="useraddradioedit"
                    >edit</label
                  >
                  <input
                    type="checkbox"
                    class="btn-check"
                    name="useraddradio"
                    id="useraddradioadmin"
                    autocomplete="off"
                    value="admin"
                    v-model="sess.useraddoption"
                  />
                  <label class="btn btn-outline-primary" for="useraddradioadmin"
                    >admin</label
                  >
                </div>
                <br /><br />
                <button
                  class="btn btn-outline-primary rounded-pill"
                  @click="adduser"
                >
                  Add
                </button>
              </div>
              <hr />
              <table class="table table-dark">
                <thead>
                  <tr>
                    <td><h4>Username</h4></td>
                    <td><h4>Action</h4></td>
                  </tr>
                </thead>
                <tbody class="table-group-divider">
                  <tr v-for="user in sess.userslist">
                    <td>
                      {{ user.username }}<br />
                      <span v-if="user.edit == 1" class="badge bg-secondary"
                        >edit</span
                      >
                      <span v-if="user.admin == 1" class="badge bg-secondary"
                        >admin</span
                      >
                    </td>
                    <td>
                      <button
                        class="btn btn-outline-primary rounded-pill"
                        data-bs-toggle="modal"
                        data-bs-target="#updateusermodal"
                        @click="
                          prepareupdateuser(
                            user.username,
                            user.edit,
                            user.admin
                          )
                        "
                      >
                        Update</button
                      >&nbsp;&nbsp;<button
                        class="btn btn-outline-primary rounded-pill"
                        data-bs-toggle="modal"
                        data-bs-target="#confirmdeleteusermodal"
                        @click="preparedeleteuser(user.username)"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      <!-- Body End -->
    </div>
    <!-- Login Modal Start -->
    <div
      class="modal fade"
      id="loginmodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Login with e-Passport</h1>
          </div>
          <div class="modal-body">
            <div class="form-floating mb-3">
              <input
                type="text"
                class="form-control"
                id="epassport"
                placeholder=" "
                v-model="sess.epassport"
              />
              <label for="epassport">e-Passport</label>
            </div>
            <div class="form-floating">
              <input
                type="password"
                class="form-control"
                id="epassword"
                placeholder=" "
                v-model="sess.epassword"
              />
              <label for="epassword">Password</label>
            </div>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="loginepassport"
            >
              Login with e-Passport
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Login Modal End -->

    <!-- Logout Modal Start -->
    <div
      class="modal fade"
      id="logoutmodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Logout</h1>
          </div>
          <div class="modal-body">Do you really want to logout?</div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="logoutepassport"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Logout Modal End -->

    <!-- Info Modal Start -->
    <div
      class="modal fade"
      id="infomodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Information</h1>
          </div>
          <div class="modal-body">{{ sess.infotxt }}</div>
          <div class="modal-footer"></div>
        </div>
      </div>
    </div>
    <!-- Info Modal End -->
    <!-- Confirm Modal Start -->
    <div
      class="modal fade"
      id="confirmmodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Delete</h1>
          </div>
          <div class="modal-body">
            Are you sure to delete tag?<br /><br />
            <h4>
              <span class="badge text-bg-light">{{ sess.tagdelete }}</span>
            </h4>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="deletetag"
            >
              Sure
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Confirm Modal End -->
    <!-- Update Doc Modal Start -->
    <div
      class="modal fade"
      id="updatedocmodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Update</h1>
          </div>
          <div class="modal-body">
            <div class="form-floating mb-3">
              <input
                type="text"
                class="form-control"
                id="filenameupdate"
                placeholder=" "
                v-model="sess.filenameupdate"
              />
              <label for="filenameupdate">File Name</label><br />
              <input
                class="form-control"
                type="file"
                id="formFileupdate"
                @change="filechupdate($event)"
              /><br />
              <li
                v-for="tag in sess.tagdocupdate"
                class="badge bg-secondary"
                @click="removetagdocupdate(tag)"
              >
                {{ tag }}
              </li>
              <br />
              <div class="dropdown">
                <a
                  class="btn btn-outline-primary dropdown-toggle rounded-pill"
                  href="#"
                  role="button"
                  data-bs-toggle="dropdown"
                  aria-expanded="false"
                >
                  Tags
                </a>

                <ul class="dropdown-menu">
                  <li v-for="tag in sess.tags">
                    <a
                      class="dropdown-item"
                      href="#"
                      @click="addtagdocupdate(tag)"
                      >{{ tag }}</a
                    >
                  </li>
                </ul>
              </div>
              <br />

              <div
                class="btn-group"
                role="group"
                aria-label="Basic radio toggle button group"
              >
                <input
                  type="radio"
                  class="btn-check"
                  name="fileradioupdate"
                  id="filepublicupdate"
                  autocomplete="off"
                  value="public"
                  v-model="sess.typeupdate"
                />
                <label class="btn btn-outline-primary" for="filepublicupdate"
                  >public</label
                >
                <input
                  type="radio"
                  class="btn-check"
                  name="fileradioupdate"
                  id="fileinternalupdate"
                  autocomplete="off"
                  value="internal"
                  v-model="sess.typeupdate"
                />
                <label class="btn btn-outline-primary" for="fileinternalupdate"
                  >internal</label
                >
                <input
                  type="radio"
                  class="btn-check"
                  name="fileradioupdate"
                  id="fileprivateupdate"
                  autocomplete="off"
                  value="private"
                  v-model="sess.typeupdate"
                />
                <label class="btn btn-outline-primary" for="fileprivateupdate"
                  >private</label
                >
              </div>

              <br v-if="sess.typeupdate != 'private'" /><br
                v-if="sess.typeupdate != 'private'"
              />
              <div
                class="btn-group"
                role="group"
                aria-label="Basic radio toggle button group"
                v-if="sess.typeupdate != 'private'"
              >
                <input
                  type="radio"
                  class="btn-check"
                  name="fileprotectradioupdate"
                  id="fileprotectupdate"
                  autocomplete="off"
                  value="protect"
                  v-model="sess.typeupdateprotect"
                />
                <label class="btn btn-outline-primary" for="fileprotectupdate"
                  >protect</label
                >
                <input
                  type="radio"
                  class="btn-check"
                  name="fileprotectradioupdate"
                  id="fileunprotectupdate"
                  autocomplete="off"
                  value="unprotect"
                  v-model="sess.typeupdateprotect"
                />
                <label class="btn btn-outline-primary" for="fileunprotectupdate"
                  >unprotect</label
                >
              </div>
              <br /><br />
            </div>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="updatedoc"
            >
              Update
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Update Doc Modal End -->
    <!-- Confirm Delete File Modal Start -->
    <div
      class="modal fade"
      id="confirmdeletefilemodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Delete</h1>
          </div>
          <div class="modal-body">
            Are you sure to delete File?<br /><br />
            <h4 class="text-warning">
              {{ sess.filenamedelete }}
            </h4>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="deletefile"
            >
              Sure
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Confirm Delete File Modal End -->
    <!-- Confirm Delete User Modal Start -->
    <div
      class="modal fade"
      id="confirmdeleteusermodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Delete</h1>
          </div>
          <div class="modal-body">
            Are you sure to delete User?<br /><br />
            <h4 class="text-warning">
              {{ sess.usernamedelete }}
            </h4>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="deleteuser"
            >
              Sure
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Confirm Delete User Modal End -->
    <!-- Update User Modal Start -->
    <div
      class="modal fade"
      id="updateusermodal"
      data-bs-backdrop="true"
      data-bs-keyboard="false"
      tabindex="-1"
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <img src="/logo.png" width="45" />
            <h1 class="modal-title fs-5">Update</h1>
          </div>
          <div class="modal-body">
            <div class="form-floating mb-3">
              <input
                type="text"
                class="form-control"
                id="usernameupdate"
                placeholder=" "
                readonly
                v-model="sess.usernameupdate"
              />
              <label for="usernameupdate">Username</label><br />
              <div
                class="btn-group"
                role="group"
                aria-label="Basic radio toggle button group"
              >
                <input
                  type="checkbox"
                  class="btn-check"
                  name="userupdateradio"
                  id="userupdateradioedit"
                  autocomplete="off"
                  value="edit"
                  v-model="sess.userupdateeditoption"
                />
                <label class="btn btn-outline-primary" for="userupdateradioedit"
                  >edit</label
                >
                <input
                  type="checkbox"
                  class="btn-check"
                  name="userupdateradio"
                  id="userupdateradioadmin"
                  autocomplete="off"
                  value="admin"
                  v-model="sess.userupdateadminoption"
                />
                <label
                  class="btn btn-outline-primary"
                  for="userupdateradioadmin"
                  >admin</label
                >
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-danger rounded-pill"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-outline-primary rounded-pill"
              data-bs-dismiss="modal"
              @click="updateuser"
            >
              Update
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Update User Modal End -->
  </div>
</template>
