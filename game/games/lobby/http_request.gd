extends HTTPRequest

signal lobby_created(code: String)
signal lobby_creation_failed(message: String)

var auth_token: String = "ff4a1cc7-adf1-4beb-aefa-26a5ebee687c" # set this after login

func create_lobby():
	request_completed.connect(_on_create_lobby_completed)

	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + auth_token
	]

	var error = request(
		"http://localhost:8080/api/lobbies",
		headers,
		HTTPClient.METHOD_POST,
		"" # or a JSON body if HandleCreateLobby expects one
	)
	if error != OK:
		push_error("An error occurred creating the lobby.")
		lobby_creation_failed.emit("Request failed to send.")

func _on_create_lobby_completed(result, response_code, headers, body):
	if response_code != 200 and response_code != 201:
		var msg = "Create lobby failed: %s" % body.get_string_from_utf8()
		push_error(msg)
		lobby_creation_failed.emit(msg)
		return

	var json = JSON.new()
	var parse_error = json.parse(body.get_string_from_utf8())
	if parse_error != OK:
		push_error("Failed to parse JSON response.")
		lobby_creation_failed.emit("Failed to parse JSON response.")
		return

	var response = json.get_data()
	print(response)

	var lobby_code = response.get("code", "")
	print("Lobby code: ", lobby_code)
	lobby_created.emit(lobby_code)
