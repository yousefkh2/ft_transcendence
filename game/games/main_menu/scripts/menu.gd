extends Control

@onready var menu: VBoxContainer = $Menu
@onready var settings: Panel = $settings

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	TranslationServer.set_locale("ENGLISH")
	settings.visible = false
	
	

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
	

func _on_settings_button_pressed() -> void:
	settings.visible = true


func _on_back_button_pressed() -> void:
	settings.visible = false
