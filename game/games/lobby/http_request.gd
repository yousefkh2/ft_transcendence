extends HTTPRequest

signal login_success
signal login_failed(message: String)
signal lobby_created(code: String)
signal lobby_creation_failed(message: String)
signal lobby_joined(lobby_data: Dictionary)
signal lobby_join_failed(message: String)

@onready var message: RichTextLabel = $"../message_interface/message"

func _translate_error(error_message: String) -> String:
	match error_message:
		"lobby not found":
			return tr("_LOBBY_NOT_FOUND_")
		"already joined this lobby":
			return tr("_LOBBY_ALREADY_JOINED_")
		"invalid credentials":
			return tr("_LOBBY_LOGIN_FAILED_")
		"lobby language does not match your current language":
			return tr("_LOBBY_LANGUAGE_MISMATCH_")
		_:
			return error_message

func _get_active_language_code() -> String:
	match GameState.selected_locale:
		"GERMAN":
			return "de"
		"POLISH":
			return "pl"
		"TURKISH":
			return "tr"
		_:
			return "en"

const AUTH_TIMEOUT_SEC := 5.0

# Only used when running from the Godot editor (no website around the game to log in through).
const DEV_EMAIL := "hallo@hallo.com"
const DEV_PASSWORD := "secret123"

# The JS callback must stay referenced, or the browser loses it.
var _js_message_cb: JavaScriptObject
var _waiting_for_token := false

# In the browser the website already logged the player in: we ask the page around the iframe for
# its session token and then ask the backend who that token belongs to. No password is involved.
func login() -> void:
	if GameState.auth_token != "" and GameState.player_name != "":
		login_success.emit.call_deferred()
		return

	if OS.has_feature("web"):
		_request_token_from_website()
	else:
		_dev_login()

func _request_token_from_website() -> void:
	var window = JavaScriptBridge.get_interface("window")
	if window == null or window.parent == null:
		_fail_login("Website not available.")
		return

	_waiting_for_token = true
	_js_message_cb = JavaScriptBridge.create_callback(_on_js_message)
	window.addEventListener("message", _js_message_cb)
	window.parent.postMessage(JSON.stringify({"type": "lingodash:auth-request"}), window.location.origin)

	get_tree().create_timer(AUTH_TIMEOUT_SEC).timeout.connect(func():
		if _waiting_for_token:
			_stop_listening()
			_fail_login("The website did not send a login token.")
	)

func _on_js_message(args: Array) -> void:
	if not _waiting_for_token or args.is_empty():
		return
	var event = args[0]
	var window = JavaScriptBridge.get_interface("window")
	# Only accept messages from our own website.
	if str(event.origin) != str(window.location.origin):
		return
	if typeof(event.data) != TYPE_STRING:
		return

	var data = JSON.parse_string(event.data)
	if typeof(data) != TYPE_DICTIONARY or data.get("type", "") != "lingodash:auth":
		return

	_stop_listening()
	var token = data.get("token")
	if typeof(token) != TYPE_STRING or token == "":
		_fail_login("Not logged in on the website.")
		return

	GameState.auth_token = token
	_fetch_me()

func _stop_listening() -> void:
	_waiting_for_token = false
	if _js_message_cb != null:
		JavaScriptBridge.get_interface("window").removeEventListener("message", _js_message_cb)
		_js_message_cb = null

# Asks the backend which user the token belongs to (the backend verifies the token's signature).
func _fetch_me() -> void:
	request_completed.connect(_on_me_completed, CONNECT_ONE_SHOT)
	var headers = ["Authorization: Bearer " + GameState.auth_token]
	var error = request("http://localhost:8080/api/users/me", headers, HTTPClient.METHOD_GET)
	if error != OK:
		_fail_login("Profile request failed to send.")

func _on_me_completed(result, response_code, headers, body):
	if response_code != 200:
		GameState.auth_token = ""
		_fail_login("Session rejected: %s" % body.get_string_from_utf8())
		return

	var response = JSON.parse_string(body.get_string_from_utf8())
	if typeof(response) != TYPE_DICTIONARY or str(response.get("username", "")) == "":
		GameState.auth_token = ""
		_fail_login("Invalid profile response.")
		return

	GameState.player_name = response["username"]
	login_success.emit()

func _fail_login(reason: String) -> void:
	push_error("Login failed: " + reason)
	message.text = tr("_LOBBY_LOGIN_FAILED_")
	login_failed.emit(reason)

func _dev_login() -> void:
	request_completed.connect(_on_login_completed, CONNECT_ONE_SHOT)
	var headers = ["Content-Type: application/json"]
	var body = JSON.stringify({"email": DEV_EMAIL, "password": DEV_PASSWORD})
	var error = request(
		"http://localhost:8080/api/v2/auth/login",
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
			message.text = _translate_error(str(err_response.get("message", "")))
		else:
			message.text = tr("_LOBBY_LOGIN_FAILED_")

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
	var logged_in_username = response.get("username", "")
	if logged_in_username == "":
		login_failed.emit("No username in login response.")
		return

	GameState.auth_token = token
	GameState.player_name = logged_in_username
	login_success.emit()

func create_lobby() -> void:
	request_completed.connect(_on_create_lobby_completed, CONNECT_ONE_SHOT)
	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + GameState.auth_token
	]
	var body = JSON.stringify({
		"gameMode": "apartment_setup",
		"lang": GameState.lobby_locale
	})
	var error = request(
		"http://localhost:8080/api/v2/lobbies",
		headers,
		HTTPClient.METHOD_POST,
		body
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
	GameState.lobby_data = response
	GameState.lobby_code = response.get("code", "")
	GameState.game_mode = response.get("gameMode", "")
	GameState.status = response.get("status", "")
	GameState.player_count = int(response.get("playerCount", 0))
	GameState.game_lang = response.get("currLang", "")


	lobby_created.emit(GameState.lobby_code)

func join_lobby(code: String) -> void:
	request_completed.connect(_on_join_lobby_completed, CONNECT_ONE_SHOT)
	var headers = [
		"Content-Type: application/json",
		"Authorization: Bearer " + GameState.auth_token
	]
	var body = JSON.stringify({
		"lang": _get_active_language_code()
	})
	var error = request(
		"http://localhost:8080/api/v2/lobbies/%s/join" % code,
		headers,
		HTTPClient.METHOD_POST,
		body
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
			message.text = _translate_error(str(err_response.get("message", "")))
		else:
			message.text = tr("_LOBBY_LOGIN_FAILED_")
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

	lobby_joined.emit(response)
