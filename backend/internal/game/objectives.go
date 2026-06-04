package game

import (
	"time"
	"transcendence/backend/internal/model"
)

var apartmentObjectives = []model.Objective{
	{
		ID:       "plant_right_of_sofa",
		ObjectID: "plant",
		Relation: "right_of",
		TargetID: "sofa",
	},
	{
		ID:       "lamp_left_of_sofa",
		ObjectID: "lamp",
		Relation: "left_of",
		TargetID: "sofa",
	},
	{
		ID:       "table_in_front_of_sofa",
		ObjectID: "table",
		Relation: "in_front_of",
		TargetID: "sofa",
	},
}

const RoundDuration = 60 * time.Second

func InitialObjectPositions() map[string]model.Position {
	return map[string]model.Position{
		"sofa":  {X: 0, Y: 1},
		"plant": {X: 3, Y: 1},
		"lamp":  {X: 2, Y: 3},
		"table": {X: 4, Y: 4},
	}
}

func MatchingObjective(message model.ClientMessage) (model.Objective, bool) {
	for _, objective := range apartmentObjectives {
		if message.ObjectID == objective.ObjectID &&
			message.Relation == objective.Relation &&
			message.TargetID == objective.TargetID {
			return objective, true
		}
	}

	return model.Objective{}, false
}

func RelationMatches(object model.Position, relation string, target model.Position) bool {
	switch relation {
	case "right_of":
		return object.Y == target.Y && object.X == target.X+1
	case "left_of":
		return object.Y == target.Y && object.X == target.X-1
	case "in_front_of":
		return object.X == target.X && object.Y == target.Y+1
	default:
		return false
	}
}

func CompletedObjectivesFromPositions(positions map[string]model.Position) map[string]bool {
	completed := make(map[string]bool)

	for _, objective := range apartmentObjectives {
		objectPosition, objectExists := positions[objective.ObjectID]
		targetPosition, targetExists := positions[objective.TargetID]
		if !objectExists || !targetExists {
			continue
		}
		if RelationMatches(objectPosition, objective.Relation, targetPosition) {
			completed[objective.ID] = true
		}
	}

	return completed
}

func CompletedObjectiveIDs(completed map[string]bool) []string {
	ids := make([]string, 0, len(completed))
	for _, objective := range apartmentObjectives {
		if completed[objective.ID] {
			ids = append(ids, objective.ID)
		}
	}

	return ids
}

func AllObjectivesCompleted(completed map[string]bool) bool {
	for _, objective := range apartmentObjectives {
		if !completed[objective.ID] {
			return false
		}
	}

	return true
}

func RemainingSeconds(deadline time.Time) int {
	if deadline.IsZero() {
		return 0
	}

	remaining := time.Until(deadline)
	if remaining <= 0 {
		return 0
	}

	return int(remaining.Seconds())
}

func RoundExpired(deadline time.Time) bool {
	if deadline.IsZero() {
		return false
	}

	return time.Now().After(deadline)
}

func CopyObjectPositions(positions map[string]model.Position) map[string]model.Position {
	copied := make(map[string]model.Position, len(positions))
	for objectID, position := range positions {
		copied[objectID] = position
	}

	return copied
}
