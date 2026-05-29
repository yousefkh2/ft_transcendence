import { createApp, h, ref } from "vue";
import "./style.css";

const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
type Position = {
  x: number;
  y: number;
};

type ObjectPositions = Record<string, Position>;

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
    const objectPositions = ref<ObjectPositions>({});
    const selectedObjectID = ref("plant"); // UI starts with plant selected


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
          completedObjectives.value = message.completedObjectives || [];
          objectPositions.value = message.objectPositions || {};
        }
        if (message.type === "game.state_updated") {
          completedObjectives.value = message.completedObjectives || [];
          objectPositions.value = message.objectPositions || {};
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

    // function movePlantRightOfSofa() {
    //   if (!socket || socket.readyState !== WebSocket.OPEN) {
    //     roomStatus.value = "Connect the WebSocket first";
    //     return;
    //   }

    //   socket.send(
    //     JSON.stringify({
    //       type: "game.object_moved",
    //       objectId: "plant",
    //       relation: "right_of",
    //       targetId: "sofa",
    //     }),
    //   );
    // }

    function moveSelectedObject(deltaX: number, deltaY: number) {
      if (!socket || socket.readyState !== WebSocket.OPEN) {
        roomStatus.value = "Connect the WebSocket first";
        return;
      }
      
      const currentPosition = objectPositions.value[selectedObjectID.value];

      if (!currentPosition) {
        roomStatus.value = "Join a room before moving objects";
        return;
      }

      const nextX = currentPosition.x + deltaX;
      const nextY = currentPosition.y + deltaY;

      if (nextX < 0 || nextX > 4 || nextY < 0 || nextY > 4) {
      roomStatus.value = "Object cannot move outside the grid";
      return;
    }

      socket.send(
        JSON.stringify({
          type: "game.object_moved",
          objectId: selectedObjectID.value,
          x: nextX,
          y: nextY,
        }),
      );

    }


    function objectAt(x: number, y: number) {
      return Object.entries(objectPositions.value).find(
        ([, position]) => position.x === x && position.y === y,
      )?.[0];
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
        h("section", { class: "apartment-board" }, [
        h("h2", "Apartment Grid"),
        h("div", { class: "object-controls" }, [
          h("label", [
            "Selected object",
            h(
              "select",
              {
                value: selectedObjectID.value,
                onChange: (event: Event) => {
                  selectedObjectID.value = (event.target as
        HTMLSelectElement).value;
                },
              },
              ["plant", "lamp", "table", "sofa"].map((objectID) =>
                h("option", { value: objectID }, objectID),
              ),
            ),
          ]),
          h("div", { class: "move-controls" }, [
            h("button", { class: "button", onClick: () => moveSelectedObject(0,
        -1) }, "Up"),
            h("button", { class: "button", onClick: () => moveSelectedObject(-1,
        0) }, "Left"),
            h("button", { class: "button", onClick: () => moveSelectedObject(1,
        0) }, "Right"),
            h("button", { class: "button", onClick: () => moveSelectedObject(0,
        1) }, "Down"),
          ]),
        ]),
        h(
          "div",
          { class: "grid-board" },
          Array.from({ length: 25 }, (_, index) => {
            const x = index % 5;
            const y = Math.floor(index / 5);
            const objectID = objectAt(x, y);

            return h("div", { class: "grid-cell" }, objectID || "");
          }),
        ),
      ]),
      ]);
  },
};

createApp(App).mount("#app");
