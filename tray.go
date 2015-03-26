package nonograms

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	GridSize = 15
	// TextureSize = 15.0
)

type Tray struct {
	Background *sf.RectangleShape
	FacingSide *sf.RectangleShape
	Shadow *sf.RectangleShape

	Highlight *sf.VertexArray
	Grid *sf.VertexArray
	Texture *sf.Texture

	Origin sf.Vector2f
	Rows int
	Columns int
}

func NewTray (origin sf.Vector2f, color sf.Color) *Tray {
	background, _ := sf.NewRectangleShape()
	background.SetFillColor(color)
	background.SetPosition(origin)

	facingSide, _ := sf.NewRectangleShape()
	facingSide.SetFillColor(sf.Color{0, 0, 0, 38})

	shadow, _ := sf.NewRectangleShape()
	shadow.SetFillColor(sf.Color{0, 0, 0, 25})

	highlight, _ := sf.NewVertexArray()
	highlight.PrimitiveType = sf.PrimitiveLines
	highlight.Resize(2)
	highlight.Vertices[0].Color = sf.Color{255, 255, 255, 255}
	highlight.Vertices[1].Color = sf.Color{255, 255, 255, 255}

	grid, _ := sf.NewVertexArray()
	grid.PrimitiveType = sf.PrimitiveQuads

	texture, _ := sf.NewTextureFromFile("../assets/tile.png", nil)

	return &Tray{
		Background: background,
		FacingSide: facingSide,
		Shadow: shadow,
		Highlight: highlight,
		Grid: grid,
		Texture: texture,
		Origin: origin,
	}
}

func (t *Tray) SetSize (rows, columns int) {
	vectors := rows * columns * 4
	t.Rows = rows
	t.Columns = columns

	t.Grid.Resize(vectors)
	t.Populate()

	cf := float32(columns)
	rf := float32(rows)

	var width float32 = 3 * (cf / 5) + (cf - (cf/5)) + (cf * GridSize) + 3
	var height float32 = 3 * (rf / 5) + (rf - (rf/5)) + (rf * GridSize) + 8

	log.Printf("width: %f, height: %f", width, height)

	size_vector := sf.Vector2f{ width, height }

	t.Background.SetSize(size_vector)
	log.Printf("%#v", t.Background.GetSize())

	t.FacingSide.SetSize(sf.Vector2f{ float32(width), 5 })
	t.FacingSide.SetPosition(sf.Vector2f{ t.Origin.X, t.Origin.Y + float32(height - 5) })

	t.Shadow.SetSize(sf.Vector2f{ float32(width), 15 })
	t.Shadow.SetPosition(sf.Vector2f{ t.Origin.X, t.Origin.Y + float32(height) })

	t.Highlight.Vertices[0].Position = sf.Vector2f{ t.Origin.X, t.Origin.Y + float32(height - 6) }
	t.Highlight.Vertices[1].Position = sf.Vector2f{ t.Origin.X + float32(width), t.Origin.Y + float32(height - 6) }
}

func (t *Tray) SetState(quad []sf.Vertex, state byte) {
	base := float32(state)

	quad[0].TexCoords = sf.Vector2f{ base * TextureSize, 0 }
	quad[1].TexCoords = sf.Vector2f{ (base + 1) * TextureSize, 0 }
	quad[2].TexCoords = sf.Vector2f{ (base + 1) * TextureSize, TextureSize }
	quad[3].TexCoords = sf.Vector2f{ base * TextureSize, TextureSize }
}


func (t *Tray) Populate () {
	log.Printf("Vertex Count: %d", t.Grid.GetVertexCount())

	var q []sf.Vertex
	var index int
	var padding_x int
	var padding_y int
	var cellpadding_x int
	var cellpadding_y int

	for y := 0; y < t.Rows; y++ {
		if y % 5 == 0 {
			padding_y += 3
		} else {
			cellpadding_y += 1
		}

		for x := 0; x < t.Columns; x++ {
			index = (x + y * t.Columns) * 4

			q = t.Grid.Vertices[index:index+4]

			if x % 5 == 0 {
				padding_x += 3
			} else {
				cellpadding_x += 1
			}

			q[0].Color = sf.ColorWhite()
			q[1].Color = sf.ColorRed()
			q[2].Color = sf.ColorBlue()
			q[3].Color = sf.ColorYellow()

			q[0].Position = sf.Vector2f{
				t.Origin.X + float32((x * GridSize) + padding_x + cellpadding_x),
				t.Origin.Y + float32((y * GridSize) + padding_y + cellpadding_y),
			}

			q[1].Position = sf.Vector2f{
				t.Origin.X + float32(((x + 1) * GridSize) + padding_x + cellpadding_x),
				t.Origin.Y + float32((y * GridSize) + padding_y + cellpadding_y),
			}

			q[2].Position = sf.Vector2f{
				t.Origin.X + float32(((x + 1) * GridSize) + padding_x + cellpadding_x),
				t.Origin.Y + float32(((y +1) * GridSize) + padding_y + cellpadding_y),
			}

			q[3].Position = sf.Vector2f{
				t.Origin.X + float32((x * GridSize) + padding_x + cellpadding_x),
				t.Origin.Y + float32(((y + 1) * GridSize) + padding_y + cellpadding_y),
			}

			t.SetState(q, ByteEmpty)
		}

		padding_x = 0
		cellpadding_x = 0
	}
}

func (t *Tray) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	r := sf.RenderStates{
		Shader: nil,
		BlendMode: sf.BlendAlpha,
		Transform: sf.TransformIdentity(),
		Texture: t.Texture,
	}

	h := sf.RenderStates{
		Shader: nil,
		BlendMode: sf.BlendMultiply,
		Transform: sf.TransformIdentity(),
		Texture: nil,
	}

	t.Background.Draw(target, renderStates)
	t.FacingSide.Draw(target, renderStates)
	t.Shadow.Draw(target, renderStates)
	t.Highlight.Draw(target, h)
	t.Grid.Draw(target, r)
}

func (t *Tray) Logic() {

}

func (t *Tray) HandleEvent(event sf.Event) {

}
