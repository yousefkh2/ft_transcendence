extends OptionButton


func _on_item_selected(index: int) -> void:
	var locale := "ENGLISH"
	match index:
		0:
			locale = "ENGLISH"
		1:
			locale = "GERMAN"
		2:
			locale = "POLISH"
		3:
			locale = "TURKISH"
	TranslationServer.set_locale(locale)
	GameState.selected_locale = locale


func _on_ready() -> void:
	match GameState.selected_locale:
		"GERMAN":
			select(1)
		"POLISH":
			select(2)
		"TURKISH":
			select(3)
		_:
			select(0)
	var popup = get_popup()
	popup.add_theme_font_override("font", preload("res://game_files/fonts/GrapeSoda.ttf"))
	popup.add_theme_font_size_override("font_size", 48)
