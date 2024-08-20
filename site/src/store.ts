import { create } from "zustand";

export const useNamespaceStore = create(() => {
  listOfNamespaces: [];
  getNamespaces: () => {
    const URL = "http://localhost:8080/api/namespaces";
    fetch(URL, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    }).then(response => {
        console.log(response)
    })
  };
});
