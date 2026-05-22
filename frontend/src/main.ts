import { createApp } from "vue";
import "./style.css";

const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";

const App = {
  template: `
    <main class="app-shell">
      <section class="hero">
        <p class="eyebrow">ft_transcendence</p>
        <h1>Apartment Setup</h1>
        <p class="summary">
          Dockerized Vue, Go, and PostgreSQL foundation for the two-player
          language game.
        </p>

        <div class="actions">
          <a class="button primary" :href="apiUrl + '/health'">Backend Health</a>
          <a class="button" :href="apiUrl + '/health/db'">Database Health</a>
        </div>
      </section>

      <section class="status-grid" aria-label="Project status">
        <article>
          <strong>Frontend</strong>
          <span>Vue + Vite dev server</span>
        </article>
        <article>
          <strong>Backend</strong>
          <span>Go API foundation</span>
        </article>
        <article>
          <strong>Database</strong>
          <span>PostgreSQL service in Compose</span>
        </article>
      </section>
    </main>
  `,
  data() {
    return { apiUrl };
  },
};

createApp(App).mount("#app");
