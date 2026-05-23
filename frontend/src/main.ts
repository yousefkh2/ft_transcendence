import { createApp, h } from "vue";
import "./style.css";

const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";

const App = {
  setup() {
    const status = [
      ["Frontend", "Vue + Vite dev server"],
      ["Backend", "Go API foundation"],
      ["Database", "PostgreSQL service in Compose"],
    ];

    return () =>
      h("main", { class: "app-shell" }, [
        h("section", { class: "hero" }, [
          h("p", { class: "eyebrow" }, "ft_transcendence"),
          h("h1", "Apartment Setup"),
          h(
            "p",
            { class: "summary" },
            "Dockerized Vue, Go, and PostgreSQL foundation for the two-player language game.",
          ),
          h("div", { class: "actions" }, [
            h("a", { class: "button primary", href: `${apiUrl}/health` }, "Backend Health"),
            h("a", { class: "button", href: `${apiUrl}/health/db` }, "Database Health"),
          ]),
        ]),
        h(
          "section",
          { class: "status-grid", "aria-label": "Project status" },
          status.map(([title, text]) =>
            h("article", [h("strong", title), h("span", text)]),
          ),
        ),
      ]);
  },
};

createApp(App).mount("#app");
