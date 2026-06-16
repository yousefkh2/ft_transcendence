extends Control

@onready var role_label = $role
@onready var http_request = $HTTPRequest

func _ready():
	role_label.text = "Loading role..."

	http_request.request_completed.connect(_on_request_completed)

	var err = http_request.request(
		"http://localhost:8080/api/player"
	)

	if err != OK:
		role_label.text = "Request failed"


func _on_request_completed(
	result,
	response_code,
	headers,
	body
):

	if response_code != 200:
		role_label.text = "Server Error"
		return

	var text = body.get_string_from_utf8()

	print(text)

	var data = JSON.parse_string(text)

	if data:
		role_label.text = "Role: " + str(data.role)
