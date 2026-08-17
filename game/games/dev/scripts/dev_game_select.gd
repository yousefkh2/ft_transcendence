extends Control

@onready var dev_game_select: Control = $"."
@onready var background: Panel = $background
@onready var settings: Panel = $settings

func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _back_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")


func _on_mission_control_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/livingroom/mission_control.tscn")


func _on_on_site_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/livingroom/on_site.tscn")
