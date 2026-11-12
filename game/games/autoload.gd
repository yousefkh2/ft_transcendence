extends Node

## Emitted when the player's profile picture has been downloaded (or removed).
signal avatar_changed

const API_URL := "http://localhost:8080"
const ROUNDED_CORNERS_SHADER := preload("res://game_files/shaders/rounded_corners.gdshader")
## Corner radius of the profile picture, as a share of its side (0.5 would be a circle).
const AVATAR_CORNER_RADIUS := 0.2

var auth_token: String = ""
var lobby_code: String = ""
var lobby_data: Dictionary = {}
var game_mode: String = ""
var status: String = ""
var player_count: int = 0
var player_name: String = ""  
var player_xp: int = -1  # -1 until the profile is loaded
var game_lang: String = ""
var selected_locale: String = "ENGLISH"
var lobby_locale: String = "en"
var player_role: int = 0
var lobby_message: String = ""

# The player's XP for the profile box, like "1250 XP". Same wording as the website in every language.
func xp_text() -> String:
	if player_xp < 0:
		return tr("_PROFILE_RANK_LOADING_")
	return "%d XP" % player_xp

# --- Profile picture ----------------------------------------------------------

# The picture the website shows for this player; null when they have none (the scene keeps its placeholder).
var avatar_texture: Texture2D
# The avatarUrl from the profile. It changes whenever the player uploads a new picture.
var _avatar_url := ""
var _avatar_request: HTTPRequest

## Called with the avatarUrl from /api/users/me ("" when the player has no picture).
func set_avatar_url(url: String) -> void:
	if url == _avatar_url:
		return
	_avatar_url = url
	_avatar_request.cancel_request()
	if url == "":
		avatar_texture = null
		avatar_changed.emit()
		return
	# The picture endpoint is public, so no token is needed.
	if _avatar_request.request(API_URL + url) != OK:
		push_error("Profile picture request failed to send.")
		_avatar_url = ""

func _on_avatar_completed(result: int, response_code: int, _headers, body: PackedByteArray) -> void:
	if result != HTTPRequest.RESULT_SUCCESS or response_code != 200 or body.size() < 4:
		push_error("Could not download the profile picture (%d)." % response_code)
		_avatar_url = ""
		return

	# The backend stores pictures as PNG or JPEG; tell them apart by their first bytes.
	var image := Image.new()
	var error := ERR_FILE_UNRECOGNIZED
	if body[0] == 0x89 and body[1] == 0x50:
		error = image.load_png_from_buffer(body)
	elif body[0] == 0xFF and body[1] == 0xD8:
		error = image.load_jpg_from_buffer(body)
	if error != OK:
		push_error("Could not read the profile picture.")
		_avatar_url = ""
		return

	# Crop to the centred square so the rounded corners sit on the picture's real corners.
	var side := mini(image.get_width(), image.get_height())
	if image.get_width() != image.get_height():
		image = image.get_region(Rect2i((image.get_width() - side) / 2, (image.get_height() - side) / 2, side, side))

	avatar_texture = ImageTexture.create_from_image(image)
	avatar_changed.emit()

## Puts the player's picture into a profile box as a rounded square, centred in the box.
## Keeps the placeholder when the player has no picture.
func apply_avatar(rect: TextureRect) -> void:
	if avatar_texture == null:
		return
	rect.expand_mode = TextureRect.EXPAND_IGNORE_SIZE
	# Centred (not covered) so the whole square is drawn and the shader's UVs span the picture.
	rect.stretch_mode = TextureRect.STRETCH_KEEP_ASPECT_CENTERED
	rect.texture = avatar_texture
	var rounded := ShaderMaterial.new()
	rounded.shader = ROUNDED_CORNERS_SHADER
	rounded.set_shader_parameter("radius", AVATAR_CORNER_RADIUS)
	rect.material = rounded

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	_avatar_request = HTTPRequest.new()
	add_child(_avatar_request)
	_avatar_request.request_completed.connect(_on_avatar_completed)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
