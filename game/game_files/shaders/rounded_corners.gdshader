shader_type canvas_item;

// Rounds the corners of a square texture (used for the player's profile picture).
// radius is a share of the side length: 0.0 = sharp corners, 0.5 = circle.
uniform float radius : hint_range(0.0, 0.5) = 0.2;

void fragment() {
	// Signed distance to a rounded square centred at (0.5, 0.5): negative inside.
	vec2 q = abs(UV - 0.5) - (0.5 - radius);
	float dist = length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - radius;
	// Soften the edge over about one pixel so the corners do not look jagged.
	float edge = fwidth(dist);
	COLOR.a *= 1.0 - smoothstep(-edge, edge, dist);
}
