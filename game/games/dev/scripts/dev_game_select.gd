extends Control

@onready var dev_game_select: Control = $"."
@onready var background: Panel = $background
@onready var settings: Panel = $settings

func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass

#DRAG AND DROP
func _on_dad_mission_control_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/livingroom/mission_control.tscn")
	
func _on_dad_on_site_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/livingroom/on_site.tscn")

#WALK AND TALK
func _on_wat_mission_control_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/walk_and_talk/mission_control.tscn")
	
func _on_wat_on_site_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/walk_and_talk/on_site.tscn")



#TAG THE PRICE
func _on_ttp_mission_control_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/tag_the_price/mission_control.tscn")

func _on_ttp_on_site_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/tag_the_price/on_site.tscn")




func _on_back_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
