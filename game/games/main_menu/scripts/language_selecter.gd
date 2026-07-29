extends OptionButton


func _on_item_selected(index: int) -> void:
	match index:
		0:
			TranslationServer.set_locale("ENGLISH")
		1:
			TranslationServer.set_locale("GERMAN")
		2:
			TranslationServer.set_locale("POLISH")
		3:
			TranslationServer.set_locale("TURKISH")


func _on_ready() -> void:
	var popup = get_popup()
	popup.add_theme_font_override("font", preload("res://game_files/fonts/GrapeSoda.ttf"))
	popup.add_theme_font_size_override("font_size", 48)
