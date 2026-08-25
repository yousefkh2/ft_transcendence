extends Control

@onready var lobby_code: Label = $lobby_code_interface/lobby_code

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	lobby_code.text = GameState.lobby_code


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
