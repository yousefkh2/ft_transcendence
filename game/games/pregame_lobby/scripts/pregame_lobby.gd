extends Control

@onready var lobby_code: Label = $lobby_code_interface/lobby_code
@onready var game_type: Label = $create_lobby_interface/game_type
@onready var player_1_name: Label = $player_interface/player_1_name
@onready var player_count: Label = $player_count_interface/player_count
@onready var lobby_lang: Label = $create_lobby_interface/lobby_lang

var websocket_url = "ws://localhost:8080/ws"
var message_to_send = "TEST TEST TEST"

@onready var _client : web_socket_client = $web_socket_client

func _ready() -> void:
	lobby_code.text = GameState.lobby_code
	game_type.text = GameState.game_mode
	print(GameState.player_name)
	print(GameState.lobby_data)
	#print(GameState.lobby_lang)
	player_1_name.text = GameState.player_name
	player_count.text = str(GameState.player_count)
	print(GameState.lobby_data)
	print(GameState.game_lang)
	lobby_lang.text = GameState.game_lang
	print("Attemting to connect to server...")
	_connect_to_matchmaking_server()
	

func _connect_to_matchmaking_server():
	var error = _client.connect_to_url(websocket_url)
	if (error != OK):
		print("Error connecting to websocket: %s " % [websocket_url])

func _on_websocket_message_recieved(message):
	print("Message received: %s " % message)

func _on_websocket_client_connection_close():
	var ws = _client.get_socket()
	print("Client disconnected with code %s, reason: %s" % [ws.get_close_code(), ws.get_close_reason()])

func _on_websocket_client_connected_to_server():
	print("Client connected to server")
	

func _on_lobby_button_pressed() -> void:
	get_tree().change_scene_to_file("res://games/lobby/lobby.tscn")
 

func _process(delta: float) -> void:
	pass


func _on_send_websocket_message_pressed() -> void:
	print("send out message from client to server")
	_client.send("from cleint to server")
	
	
