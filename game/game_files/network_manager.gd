extends Node

signal connected
signal room_joined(role)
signal connection_failed(message)

const WS_URL := "ws://localhost:8080/ws"

var socket := WebSocketPeer.new()

var is_connecting := false
var is_connected := false
var current_role := ""

func connect_to_server():
	#print("CONNECT_TO_SERVER CALLED")
	if is_connected or is_connecting:
		return

	is_connecting = true

	var err = socket.connect_to_url(WS_URL)

	if err != OK:
		connection_failed.emit("Failed to connect")
		is_connecting = false


func join_room(room_code: String):
	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Socket not connected yet")
		return

	var msg = {
		"type": "room.join",
		"roomCode": room_code
	}

	socket.send_text(JSON.stringify(msg))


func _process(_delta):
	socket.poll()

	match socket.get_ready_state():

		WebSocketPeer.STATE_OPEN:
			if not is_connected:
				is_connected = true
				is_connecting = false

				print("Connected to server")
				connected.emit()

			while socket.get_available_packet_count() > 0:
				var text = socket.get_packet().get_string_from_utf8()
				_handle_message(text)

		WebSocketPeer.STATE_CLOSED:
			pass


func _handle_message(text: String):
	if text.strip_edges().is_empty():
		return

	print("SERVER:", text)

	var data = JSON.parse_string(text)

	if typeof(data) != TYPE_DICTIONARY:
		print("Invalid JSON")
		return

	var msg_type = data.get("type", "")

	match msg_type:

		"room.joined":
			current_role = data["role"]

			print("ROLE:", current_role)

			room_joined.emit(current_role)

		"error":
			connection_failed.emit(data["message"])
