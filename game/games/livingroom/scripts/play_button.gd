extends Button

func _on_pressed() -> void:
	print(NetworkManager.current_role)
	if NetworkManager.current_role == "mission_control":
		get_tree().change_scene_to_file(
			"res://games/livingroom/mission_control.tscn"
		)

	elif NetworkManager.current_role == "on_site":
		get_tree().change_scene_to_file(
			"res://games/livingroom/on_site.tscn"
		)

	else:
		print("Unknown role:", NetworkManager.current_role)
