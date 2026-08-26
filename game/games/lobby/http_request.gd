extends HTTPRequest

signal login_success
signal login_failed(message: String)
signal lobby_created(code: String)
signal lobby_creation_failed(message: String)

var auth_token: String = ""

# Fill these in with real credentials (or pass them in from elsewhere)
var username: String = "daniel"
var password: String = "secret123"
var email: String = "test@test.com"

func login_and_create_lobby() -> void:
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
		login_failed.emit(msg)
		return

	var json = JSON.new()
	if json.parse(body.get_string_from_utf8()) != OK:
		login_failed.emit("Failed to parse login response.")
		return

	var response = json.get_data()
	print(response) # check the actual field name here!

	auth_token = response.get("token", "")
	GameState.auth_token = auth_token
	if auth_token == "":
		login_failed.emit("No token in login response.")
		return

	login_success.emit()
	create_lobby() # chain straight into lobby creation

func create_lobby():
	request_completed.connect(_on_create_lobby_completed, CONNECT_ONE_SHOT)
	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + auth_token
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
	print(response)

	var lobby_code = response.get("code", "")
	var player_name = username
	var player_count = response.get("playerCount", "")
	var status = response.get("status", "")
	var game_mode = response.get("gameMode", "")
	
	GameState.lobby_code = lobby_code
	GameState.player_name = player_name
	GameState.player_count = player_count
	GameState.status = status
	GameState.game_mode = game_mode
	
	GameState.lobby_data = response
	print("Lobby code: ", lobby_code)
	lobby_created.emit(lobby_code)
