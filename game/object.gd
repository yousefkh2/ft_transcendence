extends Sprite2D


var dragging = false
var of = Vector2(0,0)


var snap = 64

# Called every frame. 'deltax' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	if dragging:
		#var newPos = get_global_mouse_position() - of
		position = get_global_mouse_position() - of
		print(position)


func _on_button_button_down() -> void:
	dragging = true
	of = get_global_mouse_position() - global_position


func _on_button_button_up() -> void:
	dragging = false
