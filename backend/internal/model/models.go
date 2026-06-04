package model

type ClientMessage struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode,omitempty"`
	ObjectID string `json:"objectId,omitempty"`
	Relation string `json:"relation,omitempty"`
	TargetID string `json:"targetId,omitempty"`
	X        int    `json:"x,omitempty"`
	Y        int    `json:"y,omitempty"`
	Text     string `json:"text,omitempty"`
	IsFinal  bool   `json:"isFinal,omitempty"`
}

type ServerMessage struct {
	Type                string              `json:"type"`
	RoomCode            string              `json:"roomCode,omitempty"`
	PlayerID            string              `json:"playerId,omitempty"`
	Role                string              `json:"role,omitempty"`
	CompletedObjectives []string            `json:"completedObjectives,omitempty"`
	Message             string              `json:"message"`
	ObjectPositions     map[string]Position `json:"objectPositions,omitempty"`
	RemainingSeconds    int                 `json:"remainingSeconds,omitempty"`
	Text                string              `json:"text,omitempty"`
	IsFinal             bool                `json:"isFinal,omitempty"`
}

type Objective struct {
	ID       string
	ObjectID string
	Relation string
	TargetID string
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type TranscriptionResponse struct {
	Text string `json:"text"`
}
