extends Control

@onready var role_label = $role

func _ready():
	role_label.text = "Waiting for role..."

	if OS.has_feature("web"):
		JavaScriptBridge.eval("""
            window.addEventListener('message', (event) => {
                if (event.data.type === 'player_info') {
                    window.godotPlayerInfo = event.data;
                }
            });
		""")

func _process(delta):

	if not OS.has_feature("web"):
		return

	var json = JavaScriptBridge.eval(
    "JSON.stringify(window.godotPlayerInfo || null)"
	)

	var data = JSON.parse_string(json)

	if data:
		role_label.text = "Role: " + str(data.role)

		JavaScriptBridge.eval(
            "window.godotPlayerInfo = null;"
		)
