extends HTTPRequest

signal login_success
signal login_failed(message: String)
signal lobby_created(code: String)
signal lobby_creation_failed(message: String)
signal lobby_joined(lobby_data: Dictionary)
signal lobby_join_failed(message: String)

@onready var message: Label = $"../message_interface/message"

# Fill these in with real credentials (or pass them in from elsewhere)
var username: String = "gabrijel"
var password: String = "secret123"
var email: String = "hallo@hallo.com"

func login() -> void:
	request_completed.connect(_on_login_completed, CONNECT_ONE_SHOT)
	var headers = ["Content-Type: application/json"]
	var body = JSON.stringify({
		"username": username,
		"password": password,
		"email": email
	})
	var error = request(
		"http://localhost:8080/api/auth/login",
		headers,
		HTTPClient.METHOD_POST,
		body
	)
	if error != OK:
		login_failed.emit("Login request failed to send.")

func _on_login_completed(result, response_code, headers, body):
	if response_code != 200:
		var msg = "Login failed: %s" % body.get_string_from_utf8()
		push_error(msg)

		var json = JSON.new()
		if json.parse(body.get_string_from_utf8()) == OK:
			var err_response = json.get_data()
			message.text = err_response.get("message", "")
		else:
			message.text = "Login failed."

		login_failed.emit(msg)
		return

	var json = JSON.new()
	if json.parse(body.get_string_from_utf8()) != OK:
		login_failed.emit("Failed to parse login response.")
		return

	var response = json.get_data()
	var token = response.get("token", "")
	if token == "":
		login_failed.emit("No token in login response.")
		return

	GameState.auth_token = token
	GameState.player_name = username
	GameState.player_pass = password
	GameState.player_email = email
	login_success.emit()

func create_lobby() -> void:
	request_completed.connect(_on_create_lobby_completed, CONNECT_ONE_SHOT)
	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + GameState.auth_token
	]
	var error = request(
		"http://localhost:8080/api/lobbies",
		headers,
		HTTPClient.METHOD_POST,
		""
	)
	if error != OK:
		lobby_creation_failed.emit("Request failed to send.")

func _on_create_lobby_completed(result, response_code, headers, body):
	if response_code != 200 and response_code != 201:
		var msg = "Create lobby failed: %s" % body.get_string_from_utf8()
		push_error(msg)
		lobby_creation_failed.emit(msg)
		return

	var json = JSON.new()
	if json.parse(body.get_string_from_utf8()) != OK:
		lobby_creation_failed.emit("Failed to parse JSON response.")
		return

	var response = json.get_data()
	GameState.player_name = username
	GameState.lobby_data = response
	GameState.lobby_code = response.get("code", "")
	GameState.game_mode = response.get("gameMode", "")
	GameState.status = response.get("status", "")
	GameState.player_count = int(response.get("playerCount", 0))


	lobby_created.emit(GameState.lobby_code)

func join_lobby(code: String) -> void:
	request_completed.connect(_on_join_lobby_completed, CONNECT_ONE_SHOT)
	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + GameState.auth_token
	]
	var error = request(
		"http://localhost:8080/api/lobbies/%s/join" % code,
		headers,
		HTTPClient.METHOD_POST,
		""
	)
	if error != OK:
		lobby_join_failed.emit("Request failed to send.")

func _on_join_lobby_completed(result, response_code, headers, body):
	if response_code != 200 and response_code != 201:
		var msg = "Join lobby failed: %s" % body.get_string_from_utf8()
		push_error(msg)
		
		var json = JSON.new()
		if json.parse(body.get_string_from_utf8()) == OK:
			var err_response = json.get_data()
			message.text = err_response.get("message", "")
		else:
			message.text = "Login failed."
		lobby_join_failed.emit(msg)
		return

	var json = JSON.new()
	if json.parse(body.get_string_from_utf8()) != OK:
		lobby_join_failed.emit("Failed to parse JSON response.")
		return

	var response = json.get_data()
	GameState.lobby_data = response
	GameState.lobby_code = response.get("code", "")
	GameState.game_mode = response.get("gameMode", "")
	GameState.status = response.get("status", "")
	GameState.player_count = int(response.get("playerCount", 0))
	GameState.player_name = username

	lobby_joined.emit(response)
