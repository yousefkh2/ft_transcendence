import { createApp, h, ref } from "vue";
import "./style.css";

const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
type Position = {
  x: number;
  y: number;
};

type ObjectPositions = Record<string, Position>;

const objectiveLabels: Record<string, string> = {
    plant_right_of_sofa: "Plant right of sofa",
    lamp_left_of_sofa: "Lamp left of sofa",
    table_in_front_of_sofa: "Table in front of sofa",
  };

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
    const remainingSeconds = ref(0);
    const objectPositions = ref<ObjectPositions>({});
    const selectedObjectID = ref("plant"); // UI starts with plant selected
    const transcriptionStatus = ref("No transcription session");
    const transcriptLines = ref<string[]>([]);


    let socket: WebSocket | null = null;

    let realtimeConnection: RTCPeerConnection | null = null;
    let realtimeDataChannel: RTCDataChannel | null = null;

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
          remainingSeconds.value = message.remainingSeconds || 0;
        }
        const iframe = document.getElementById("godot-game") as HTMLIFrameElement;

        iframe?.contentWindow?.postMessage(
          {
            type: "player_info",
            role: message.role,
            playerId: message.playerId,
            roomCode: message.roomCode,
          },
          "*"
        );
        if (
          message.type === "game.state_updated" ||
          message.type === "game.round_completed" ||
          message.type === "game.round_expired"
        ) {

          completedObjectives.value = message.completedObjectives || [];
          objectPositions.value = message.objectPositions || {};
          roomStatus.value = message.message;
          remainingSeconds.value = message.remainingSeconds || 0;
        }

        if (message.type === "voice.transcript") {
          const speaker = message.role || "unknown";
          const text = message.text || "";

          if (text) {
            transcriptLines.value = [
              ...transcriptLines.value,
              `${speaker}: ${text}`,
            ].slice(-6);
          }
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

    function isOnSite() {
      return role.value === "on_site";
    }

    function isMissionControl() {
      return role.value === "mission_control";
    }

    function completedCount() {
    return completedObjectives.value.length;
    }

    function timerLabel() {
      if (remainingSeconds.value <= 0) {
        return "Waiting";
      }

      return `${remainingSeconds.value}s`;
    }

    async function createTranscriptionSession(): Promise<string | null> {
      transcriptionStatus.value = "Creating transcription session...";

      try {
        const response = await fetch(`${apiUrl}/realtime/transcription-session`, {
          method: "POST",
        });

        if (!response.ok) {
          transcriptionStatus.value = `Session failed: ${response.status}`;
          return null;
        }

        const session = await response.json();
        if (session.value) {
          transcriptionStatus.value = "Realtime transcription session ready";
          return session.value;
        }
        
        transcriptionStatus.value = "Session created without client secret";
        return null;
      } catch (error) {
        console.error(error);
        transcriptionStatus.value = "Session request failed; backend unreachable";
        return null;
      }
    }

    function sendTranscriptToRoom(text: string, isFinal: boolean) {
      const cleanedText = text.trim();

      if (!cleanedText) {
        return;
      }
      
      if (!socket || socket.readyState !== WebSocket.OPEN) {
        transcriptionStatus.value = "Transcript received before room socket connected";
        return;
      }

      socket.send(
        JSON.stringify({
          type: "voice.transcript",
          text: cleanedText,
          isFinal,
        }),
      );
    }

    async function connectRealtimeTranscription() {
      transcriptionStatus.value = "Connecting to OpenAI Realtime...";

      const clientSecret = await createTranscriptionSession();
      if (!clientSecret) {
        return;
      }

      realtimeConnection = new RTCPeerConnection();

      realtimeDataChannel = realtimeConnection.createDataChannel("oai-events");

      realtimeDataChannel.onopen = () => {
        transcriptionStatus.value = "Realtime transcription connected";
      };

      realtimeDataChannel.onmessage = (event) => {
        const realtimeEvent = JSON.parse(event.data);
        console.log("OpenAI realtime event", realtimeEvent);

        if (
          realtimeEvent.type ===
          "conversation.item.input_audio_transcription.completed"
        ) {
          sendTranscriptToRoom(realtimeEvent.transcript || "", true);
        }
      };


      const mediaStream = await navigator.mediaDevices.getUserMedia({
        audio: true,
      });

      for (const track of mediaStream.getAudioTracks()) {
        realtimeConnection.addTrack(track, mediaStream);
      }


      const offer = await realtimeConnection.createOffer();
      await realtimeConnection.setLocalDescription(offer);

      const response = await fetch("https://api.openai.com/v1/realtime/calls", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${clientSecret}`,
          "Content-Type": "application/sdp",
        },
        body: offer.sdp,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("OpenAI realtime connection failed", response.status,
        errorText);
        transcriptionStatus.value = `Realtime connection failed:
        ${response.status}`;
        return;
      }


      const answer = {
        type: "answer" as RTCSdpType,
        sdp: await response.text(),
      };

      await realtimeConnection.setRemoteDescription(answer);
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
            h(
            "button",
            { class: "button", onClick: connectRealtimeTranscription },
            "Connect STT",
            ),
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
              h("strong", "Realtime STT"),
              h("span", transcriptionStatus.value),
            ]),
            h("article", [
              h("strong", "Timer"),
              h("span", timerLabel()),
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
            h("strong", "Progress"),
            h("span", `${completedCount()} / ${Object.keys(objectiveLabels).length} completed`),
          ]),
          ],
        ),
        h("section", { class: "transcript-panel" }, [
          h("h2", "Transcript"),
          h(
            "div",
            { class: "transcript-feed" },
            transcriptLines.value.length > 0
              ? transcriptLines.value.map((line) => h("p", line))
              : [h("p", "No transcript yet.")],
          ),
        ]),
        isMissionControl()
        ? h("section", { class: "mission-panel" }, [
            h("h2", "Mission Control"),
            h(
              "ul",
              { class: "objective-list" },
              Object.entries(objectiveLabels).map(([objectiveID, label]) =>
                h(
                  "li",
                  {
                    class: completedObjectives.value.includes(objectiveID)
                      ? "completed"
                      : "",
                  },
                  label,
                ),
              ),
            ),
          ])
        : null,
        isOnSite()
        ? h("section", { class: "apartment-board" }, [
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
          h("button", { class: "button", onClick: () =>
            moveSelectedObject(0, -1) }, "Up"),
          h("button", { class: "button", onClick: () =>
            moveSelectedObject(-1, 0) }, "Left"),
          h("button", { class: "button", onClick: () =>
            moveSelectedObject(1, 0) }, "Right"),
          h("button", { class: "button", onClick: () =>
            moveSelectedObject(0, 1) }, "Down"),
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
      ])
    : null,
    roomStatus.value === "round completed"
    ? h("section", { class: "round-complete" }, [
        h("strong", "Round Complete"),
        h("span", "All apartment objectives are satisfied."),
      ])
    : null,
    roomStatus.value === "round expired"
    ? h("section", { class: "round-expired" }, [
        h("strong", "Round Expired"),
        h("span", "Time ran out before all apartment objectives were satisfied."),
      ])
    : null,
    h("section", { class: "game-panel" }, [
    h("h2", "Godot Game"),
      h("iframe", {
        src: "/game.html",
        title: "Godot game",
        class: "game-iframe",
        allow: "autoplay; fullscreen; pointer-lock",
      }),
    ]),
      ]);
  },
};

createApp(App).mount("#app");
