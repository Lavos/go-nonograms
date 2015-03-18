package entities

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	GridSize = 16
	Offset = 11
)

type TileMap struct {
	VertexArray *sf.VertexArray
	Texture *sf.Texture

	Width int
	Height int
}

func NewTileMap () *TileMap {
	vm, _ := sf.NewVertexArray()
	vm.PrimitiveType = sf.PrimitiveQuads

	texture, _ := sf.NewTextureFromFile("../assets/tile.png", nil)

	return &TileMap{
		VertexArray: vm,
		Texture: texture,
	}
}

func (t *TileMap) SetSize (height, width int) {
	log.Printf("Height: %d Width: %d", height, width)

	destination_size := width * height * 4
	t.Width = width
	t.Height = height

	t.VertexArray.Resize(destination_size)
	t.Populate()
}

func (t *TileMap) QuadFromCoords(x, y int) ([]sf.Vertex, bool) {
	if x > t.Width -1 || x < 0 || y > t.Height -1 || y < 0 {
		return nil, false
	}

	index := (y * t.Width * 4) + (x * 4)
	return t.VertexArray.Vertices[index:index+4], true
}

func (t *TileMap) CoordsFromPosition(x, y int) (int, int) {
	log.Printf("Mouse Position, x: %d y: %d", x, y)

	coord_x := (x - (Offset * GridSize)) / GridSize
	coord_y := (y - (Offset * GridSize)) / GridSize

	log.Printf("Grid, x: %d y: %d", coord_x, coord_y)

	return coord_x, coord_y
}

func (t *TileMap) SetState(quad []sf.Vertex, state byte) {
	base := int(state)

	quad[0].TexCoords = sf.Vector2f{ GridToPixelsf(base), GridToPixelsf(0) }
	quad[1].TexCoords = sf.Vector2f{ GridToPixelsf(base + 1), GridToPixelsf(0) }
	quad[2].TexCoords = sf.Vector2f{ GridToPixelsf(base + 1), GridToPixelsf(1) }
	quad[3].TexCoords = sf.Vector2f{ GridToPixelsf(base), GridToPixelsf(1) }
}

func (t *TileMap) Populate () {
	log.Printf("Vertex Count: %d", t.VertexArray.GetVertexCount())

	var q []sf.Vertex
	var index int

	for y := 0; y < t.Height; y++ {
		// log.Printf("y: %d", y)

		for x := 0; x < t.Width; x++ {
			// log.Printf("x: %d", x)

			index = (x + y * t.Width) * 4

			// log.Printf("Index: %d", index)

			q = t.VertexArray.Vertices[index:index+4]

			q[0].Color = sf.ColorWhite()
			q[1].Color = sf.ColorWhite()
			q[2].Color = sf.ColorWhite()
			q[3].Color = sf.ColorWhite()

			q[0].Position = sf.Vector2f{ GridToPixelsf(Offset + x), GridToPixelsf(Offset + y) }
			q[1].Position = sf.Vector2f{ GridToPixelsf(Offset + x+1), GridToPixelsf(Offset + y) }
			q[2].Position = sf.Vector2f{ GridToPixelsf(Offset + x+1), GridToPixelsf(Offset + y+1) }
			q[3].Position = sf.Vector2f{ GridToPixelsf(Offset + x), GridToPixelsf(Offset + y+1) }

			t.SetState(q, ByteEmpty)
		}
	}
}

func (t *TileMap) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	r := sf.RenderStates{
		Shader: nil,
		BlendMode: sf.BlendAlpha,
		Transform: sf.TransformIdentity(),
		Texture: t.Texture,
	}

	t.VertexArray.Draw(target, r)
}
