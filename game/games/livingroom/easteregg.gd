extends VideoStreamPlayer

func _ready():
	self.visible = false

func _process(delta: float) -> void:
	pass


func _on_button_pressed() -> void:
	print("stream is: ", stream)
	visible = true
	play()
	print("is playing: ", is_playing())
