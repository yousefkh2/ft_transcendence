class_name DragDropCell extends Button

## A furniture slot whose icon can be dragged onto another slot to swap them.
## While dragging, the furniture image follows the mouse cursor (grabbed at the
## point where it was clicked) and the original slot is hidden. The slot under
## the cursor shows a see-through "shadow" of the piece where it would land;
## if that slot is occupied, its current piece fades out to show the swap.

const DRAG_ALPHA := 0.0
const GHOST_ALPHA := 0.45
const ICON_PADDING := 8.0
const REPLACED_ICON_ALPHA := 0.25

var _is_dragging := false
var _drop_ghost: Texture2D = null
var _is_drop_target := false

func _ready() -> void:
	set_process(false)

func _get_drag_data(at_position: Vector2) -> Variant:
	if not icon:
		return null

	# Godot positions the preview's origin at the cursor, so wrap the image in
	# a holder and offset it by the click position to keep the grab point.
	var holder := Control.new()
	holder.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var image := TextureRect.new()
	image.texture = icon
	image.expand_mode = TextureRect.EXPAND_IGNORE_SIZE
	image.stretch_mode = TextureRect.STRETCH_KEEP_ASPECT_CENTERED
	image.size = size
	image.position = -at_position
	image.mouse_filter = Control.MOUSE_FILTER_IGNORE
	holder.add_child(image)

	set_drag_preview(holder)

	_is_dragging = true
	self_modulate.a = DRAG_ALPHA
	return self

func _can_drop_data(_at_position: Vector2, data: Variant) -> bool:
	var valid: bool = data is DragDropCell and data != self
	if valid:
		_set_drop_target(true, data.icon)
	return valid

func _drop_data(_at_position: Vector2, data: Variant) -> void:
	_set_drop_target(false)
	var source: DragDropCell = data
	var temp: Texture2D = icon
	icon = source.icon
	source.icon = temp

func _process(_delta: float) -> void:
	# _can_drop_data only fires while hovering, so clear the highlight once
	# the cursor leaves this slot or the drag is over.
	var hovering := get_global_rect().has_point(get_global_mouse_position())
	if not hovering or not get_viewport().gui_is_dragging():
		_set_drop_target(false)

func _notification(what: int) -> void:
	# Sent to every Control when any drag ends, whether it was dropped or cancelled.
	if what == NOTIFICATION_DRAG_END:
		_set_drop_target(false)
		if _is_dragging:
			_is_dragging = false
			self_modulate.a = 1.0

func _draw() -> void:
	if not _is_drop_target:
		return
	if _drop_ghost:
		draw_texture_rect(_drop_ghost, _fit_rect(_drop_ghost.get_size()), false, Color(1, 1, 1, GHOST_ALPHA))

func _set_drop_target(active: bool, ghost: Texture2D = null) -> void:
	if active == _is_drop_target and ghost == _drop_ghost:
		return
	_is_drop_target = active
	_drop_ghost = ghost if active else null
	set_process(active)
	# Fade the piece that would be swapped out so the shadow reads clearly.
	for state in ["icon_normal_color", "icon_hover_color", "icon_pressed_color", "icon_focus_color"]:
		if active:
			add_theme_color_override(state, Color(1, 1, 1, REPLACED_ICON_ALPHA))
		else:
			remove_theme_color_override(state)
	queue_redraw()

func _fit_rect(texture_size: Vector2) -> Rect2:
	var area := size - Vector2.ONE * ICON_PADDING * 2.0
	var scale_factor: float = min(area.x / texture_size.x, area.y / texture_size.y)
	var fitted := texture_size * scale_factor
	return Rect2((size - fitted) / 2.0, fitted)
