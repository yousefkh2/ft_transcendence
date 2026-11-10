extends OptionButton


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	var popup = get_popup()
	popup.add_theme_font_override("font", preload("res://game_files/fonts/GrapeSoda.ttf"))
	popup.add_theme_font_size_override("font_size", 48)

func _on_item_selected(index: int) -> void:
	match index:
		0:
			GameState.lobby_locale = "en"
		1:
			GameState.lobby_locale = "de"
		2:
			GameState.lobby_locale = "pl"
		3:
			GameState.lobby_locale = "tr"


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
