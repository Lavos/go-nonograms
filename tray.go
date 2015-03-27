package nonograms

import (
	"log"
	"math"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	TrayTypePlay = iota
	TrayTypeHint

	GridSize = 15
	// TextureSize = 15.0
)

type Tray struct {
	TrayType byte
	TextureManager *TextureManager

	Background *sf.RectangleShape
	FacingSide *sf.RectangleShape
	Shadow *sf.RectangleShape

	Highlight *sf.VertexArray
	Grid *sf.VertexArray

	Origin sf.Vector2f
	Rows int
	Columns int
}

func NewTray (tm *TextureManager, origin sf.Vector2f, color sf.Color, trayType byte) *Tray {
	log.Printf("TYPE: %d", trayType)

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

	return &Tray{
		TrayType: trayType,
		TextureManager: tm,

		Background: background,
		FacingSide: facingSide,
		Shadow: shadow,
		Highlight: highlight,
		Grid: grid,
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
	x := float32(state)
	y := float32(t.TrayType)

	log.Printf("x %d y %d", x, y)

	quad[0].TexCoords = sf.Vector2f{ x * TextureSize, y * TextureSize }
	quad[1].TexCoords = sf.Vector2f{ (x + 1) * TextureSize, y * TextureSize }
	quad[2].TexCoords = sf.Vector2f{ (x + 1) * TextureSize, (y + 1) * TextureSize }
	quad[3].TexCoords = sf.Vector2f{ x * TextureSize, (y + 1) * TextureSize }
}

func (t *Tray) QuadFromCoords(coord_xi, coord_yi int) ([]sf.Vertex, bool) {
	if coord_xi > t.Columns -1 || coord_xi < 0 || coord_yi > t.Rows -1 || coord_yi < 0 {
		return nil, false
	}

	index := (coord_yi * t.Columns * 4) + (coord_xi * 4)
	return t.Grid.Vertices[index:index+4], true
}

func (t *Tray) CoordsFromPosition(x, y int) (int, int) {
	base_x := float32(x) - t.Origin.X
	base_y := float32(y) - t.Origin.Y
	portions_x := base_x / ((5 * GridSize) + 3 + 4)
	portions_y := base_y / ((5 * GridSize) + 3 + 4)

	log.Printf("Portions X: %f, Portions Y: %f", portions_x, portions_y)

	coord_xf := (base_x - (3 * portions_x) - (4 * portions_x)) / 15
	coord_yf := (base_y - (3 * portions_y) - (4 * portions_y)) / 15

	log.Printf("Coord X: %f, Coord Y: %f", coord_xf, coord_yf)

	return int(math.Floor(float64(coord_xf))), int(math.Floor(float64(coord_yf)))
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
			q[1].Color = sf.ColorWhite()
			q[2].Color = sf.ColorWhite()
			q[3].Color = sf.ColorWhite()

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
		Texture: t.TextureManager.Get("tile"),
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
