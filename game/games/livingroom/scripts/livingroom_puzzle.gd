class_name LivingroomPuzzle extends RefCounted

## Builds the living room puzzle from the lobby code, so both players (on site
## and mission control) get the same target layout without extra networking.
##
## The room has 8 spots in a 4x2 grid, indexed like this:
##   top row:    0 1 2 3
##   bottom row: 4 5 6 7
## A layout is an Array of 8 item ids, with "" for an empty spot.

const COLUMNS := 4
const SPOT_COUNT := 8
const ITEMS: Array[String] = ["lamp", "plant", "couch", "table"]

const ITEM_TEXTURES := {
	"lamp": preload("res://games/livingroom/assets/on_site_asstes/lamp.png"),
	"plant": preload("res://games/livingroom/assets/on_site_asstes/plant.png"),
	"couch": preload("res://games/livingroom/assets/on_site_asstes/sofa.png"),
	"table": preload("res://games/livingroom/assets/on_site_asstes/table.png"),
}

## Translation key prefix per item. Each item has three keys:
##   _LR_ITEM_<X>_       subject form ("THE LAMP", "DIE LAMPE", ...)
##   _LR_ITEM_<X>_SIDE_  form used after "left of" / "right of"
##   _LR_ITEM_<X>_VERT_  form used after "above" / "below"
## (the split is needed for languages with grammatical cases, e.g. Polish).
const ITEM_KEYS := {
	"lamp": "_LR_ITEM_LAMP",
	"plant": "_LR_ITEM_PLANT",
	"couch": "_LR_ITEM_COUCH",
	"table": "_LR_ITEM_TABLE",
}

## Relations between two spots; "a is <relation> b".
const RELATIONS := {
	"left": Vector2i(-1, 0),
	"right": Vector2i(1, 0),
	"above": Vector2i(0, -1),
	"below": Vector2i(0, 1),
}

static func target_layout(lobby_code: String) -> Array:
	var rng := _rng(lobby_code, "target")
	var spots := _shuffled(range(SPOT_COUNT), rng)
	var layout := _empty_layout()
	for i in ITEMS.size():
		layout[spots[i]] = ITEMS[i]
	return layout

## Where the furniture starts for the on-site player. No piece ever starts
## on its target spot, so every piece has to be moved.
static func start_layout(lobby_code: String) -> Array:
	var target := target_layout(lobby_code)
	var rng := _rng(lobby_code, "start")
	while true:
		var spots := _shuffled(range(SPOT_COUNT), rng)
		var layout := _empty_layout()
		var any_in_place := false
		for i in ITEMS.size():
			layout[spots[i]] = ITEMS[i]
			any_in_place = any_in_place or target[spots[i]] == ITEMS[i]
		if not any_in_place:
			return layout
	return []

## One clue per item. The first clue (in solving order) fixes an item to an
## exact spot; every later clue either fixes its item exactly or places it
## next to an item that is already known, so the answer is always unique.
## The clues are returned shuffled, so players still have to work it out.
static func clues(lobby_code: String) -> Array:
	var layout := target_layout(lobby_code)
	var rng := _rng(lobby_code, "clues")
	var order := _shuffled(ITEMS.duplicate(), rng)
	var known: Array[String] = []
	var result: Array[Dictionary] = []

	for item in order:
		var spot := layout.find(item)
		var options: Array[Dictionary] = []
		for other in known:
			var relation := _relation(spot, layout.find(other))
			if relation != "":
				options.append({"type": relation, "item": item, "other": other})
		# Prefer relative clues; fall back to an exact spot.
		if options.is_empty():
			result.append({"type": "spot", "item": item, "spot": spot})
		else:
			result.append(options[rng.randi_range(0, options.size() - 1)])
		known.append(item)

	return _shuffled(result, rng)

static func clue_text(clue: Dictionary) -> String:
	var subject := TranslationServer.translate(ITEM_KEYS[clue.item] + "_")
	match clue.type:
		"spot":
			var key := "_LR_CLUE_TOP_ROW_" if clue.spot < COLUMNS else "_LR_CLUE_BOTTOM_ROW_"
			return TranslationServer.translate(key).format({"a": subject, "n": clue.spot % COLUMNS + 1})
		"left", "right":
			var other := TranslationServer.translate(ITEM_KEYS[clue.other] + "_SIDE_")
			return TranslationServer.translate("_LR_CLUE_%s_" % clue.type.to_upper()).format({"a": subject, "b": other})
		_:
			var other := TranslationServer.translate(ITEM_KEYS[clue.other] + "_VERT_")
			return TranslationServer.translate("_LR_CLUE_%s_" % clue.type.to_upper()).format({"a": subject, "b": other})

## Relation of spot `a` to spot `b` if they are direct neighbours, else "".
static func _relation(a: int, b: int) -> String:
	var diff := Vector2i(a % COLUMNS, a / COLUMNS) - Vector2i(b % COLUMNS, b / COLUMNS)
	for relation in RELATIONS:
		if RELATIONS[relation] == diff:
			return relation
	return ""

static func _rng(lobby_code: String, purpose: String) -> RandomNumberGenerator:
	var rng := RandomNumberGenerator.new()
	rng.seed = ("livingroom:%s:%s" % [lobby_code, purpose]).hash()
	return rng

## Deterministic Fisher-Yates shuffle (Array.shuffle() uses the global RNG).
static func _shuffled(values: Array, rng: RandomNumberGenerator) -> Array:
	var copy := values.duplicate()
	for i in range(copy.size() - 1, 0, -1):
		var j := rng.randi_range(0, i)
		var tmp = copy[i]
		copy[i] = copy[j]
		copy[j] = tmp
	return copy

static func _empty_layout() -> Array:
	var layout := []
	layout.resize(SPOT_COUNT)
	layout.fill("")
	return layout
