extends Control

@onready var lobby_code: Label = $lobby_code_interface/lobby_code
@onready var game_type: Label = $create_lobby_interface/game_type
@onready var player_1_name: Label = $"player_interface/player_1 name"
@onready var player_count: Label = $player_count_interface/player_count

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	lobby_code.text = GameState.lobby_code
	game_type.text = GameState.game_mode
	player_1_name.text = GameState.player_name
	player_count.text = str(GameState.player_count)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
 
