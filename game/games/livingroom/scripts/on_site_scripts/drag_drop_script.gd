class_name drag_drop_cell extends Button

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass

func _get_drag_data(at_position: Vector2) -> Variant:
	if not icon:
		return null
	
	var preview: TextureRect = TextureRect.new()
	preview.texture = icon
	set_drag_preview(preview)
	
	return self
	
func _can_drop_data(at_position: Vector2, data: Variant) -> bool:
	if not data is drag_drop_cell or data == self:
		return false
	
	grab_focus()
	return true
	
func _drop_data(at_position: Vector2, data: Variant) -> void:
	var temp: Texture2D = icon
	icon = data.icon
	data.icon = temp
	
