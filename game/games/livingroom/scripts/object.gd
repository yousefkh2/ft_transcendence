extends Sprite2D

var dragging = false
var drag_offset = Vector2.ZERO

@export var slots_parent: NodePath

var slots = []
var current_slot = null

func _ready():
	for child in get_node(slots_parent).get_children():
		slots.append(child)

func _process(delta):
	if dragging:
		global_position = get_global_mouse_position() - drag_offset

func _on_button_button_down():
	dragging = true

	# Free the current slot while dragging
	if current_slot:
		current_slot.occupied = false
		current_slot = null

	drag_offset = get_global_mouse_position() - global_position

func _on_button_button_up():
	dragging = false
	snap_to_closest_slot()

func snap_to_closest_slot():
	var closest = null
	var closest_distance = INF

	for slot in slots:
		if slot.occupied:
			continue

		var dist = global_position.distance_squared_to(slot.global_position)

		if dist < closest_distance:
			closest_distance = dist
			closest = slot

	if closest:
		global_position = closest.global_position
		closest.occupied = true
		current_slot = closest
