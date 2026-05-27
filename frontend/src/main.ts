import { createApp, h, ref } from "vue";
import "./style.css";

const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";

const App = {
  setup() {
    const status = [
      ["Frontend", "Vue + Vite dev server"],
      ["Backend", "Go API foundation"],
      ["Database", "PostgreSQL service in Compose"],
    ];

    const connectionStatus = ref("Disconnected");
    const roomCode = ref("ABCD");
    const roomStatus = ref("No room request sent");
    const role = ref("No role assigned");
    const playerID = ref("No player ID assigned");
    const completedObjectives = ref<string[]>([]);

    let socket: WebSocket | null = null;

    function connect() {
      connectionStatus.value = "Connecting...";

      socket = new WebSocket("ws://localhost:8080/ws");

      socket.onopen = () => {
        connectionStatus.value = "Connected";
      };

      socket.onmessage = (event) => {
        const message = JSON.parse(event.data);

        if (message.type === "room.joined") {
          roomStatus.value = `Joined room ${message.roomCode}`;
          playerID.value = message.playerId;
          role.value = message.role;
        }
        if (message.type === "game.state_updated") {
          completedObjectives.value = message.completedObjectives;
          roomStatus.value = message.message;
        }
        if (message.type === "error") {
          roomStatus.value = message.message;
        }

      };

      socket.onclose = () => {
        connectionStatus.value = "Disconnected";
      };

      socket.onerror = () => {
        connectionStatus.value = "Connection failed";
      };
    }

    function joinRoom() {
      if (!socket || socket.readyState !== WebSocket.OPEN) {
        roomStatus.value = "Connect the WebSocket first";
        return;
      }

      socket.send(
        JSON.stringify({
          type: "room.join",
          roomCode: roomCode.value,
        }),
      );
    }

    function movePlantRightOfSofa() {
      if (!socket || socket.readyState !== WebSocket.OPEN) {
        roomStatus.value = "Connect the WebSocket first";
        return;
      }

      socket.send(
        JSON.stringify({
          type: "game.object_moved",
          objectId: "plant",
          relation: "right_of",
          targetId: "sofa",
        })
      )
    }

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
            h("button", { class: "button", onClick: connect }, "Connect WebSocket"),
            h("input", {
              value: roomCode.value,
              maxlength: 8,
              onInput: (event: Event) => {
                roomCode.value = (event.target as HTMLInputElement).value;
              },
            }),
            h("button", { class: "button", onClick: joinRoom }, "Send Join Request"),
            h("button", { class: "button", onClick: movePlantRightOfSofa }, "Move Plant Right of Sofa"),
            h("article", [h("strong", "Realtime"), h("span", connectionStatus.value)]),
          ]),
        ]),
        h(
          "section",
          { class: "status-grid", "aria-label": "Project status" },
          [
            ...status.map(([title, text]) =>
              h("article", [h("strong", title), h("span", text)]),
            ),
            h("article", [
              h("strong", "Room Event"),
              h("span", roomStatus.value),
            ]),
            h("article", [
              h("strong", "Role"),
              h("span", role.value),
            ]),
            h("article", [
              h("strong", "Player ID"),
              h("span", playerID.value),
            ]),
            h("article", [
              h("strong", "Completed Objective"),
              h("span", completedObjectives.value.join(", ") || "None"),
            ]),
          ],
        ),
      ]);
  },
};

createApp(App).mount("#app");
