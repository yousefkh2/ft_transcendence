extends VideoStreamPlayer

func _ready():
	self.visible = false

func _process(delta: float) -> void:
	pass


func _on_button_pressed() -> void:
	visible = true
	play()
