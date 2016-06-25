package nonograms

import (
	"math"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	TrayTypePlay = iota
	TrayTypeHint

	GridSize = 15
	TextureSize = 15
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

func NewTray (tm *TextureManager, color sf.Color, trayType byte) *Tray {
	background, _ := sf.NewRectangleShape()
	background.SetFillColor(color)

	facingSide, _ := sf.NewRectangleShape()
	facingSide.SetFillColor(sf.Color{0, 0, 0, 38})

	shadow, _ := sf.NewRectangleShape()
	shadow.SetFillColor(sf.Color{0, 0, 0, 25})

	highlight, _ := sf.NewVertexArray()
	highlight.PrimitiveType = sf.PrimitiveLines
	highlight.Resize(2)
	highlight.Vertices[0].Color = sf.Color{255, 255, 255, 115}
	highlight.Vertices[1].Color = sf.Color{255, 255, 255, 115}

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
	}
}

func (t *Tray) SetPosition (origin sf.Vector2f) {
	t.Origin = origin
	t.Background.SetPosition(origin)
	t.SetSize(t.Rows, t.Columns)
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

	size_vector := sf.Vector2f{ width, height }

	t.Background.SetSize(size_vector)

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

	quad[0].TexCoords = sf.Vector2f{ x * TextureSize, y * TextureSize }
	quad[1].TexCoords = sf.Vector2f{ (x + 1) * TextureSize, y * TextureSize }
	quad[2].TexCoords = sf.Vector2f{ (x + 1) * TextureSize, (y + 1) * TextureSize }
	quad[3].TexCoords = sf.Vector2f{ x * TextureSize, (y + 1) * TextureSize }
}

func (t *Tray) ToggleHighlight(quad []sf.Vertex) {
	if quad[0].Color == sf.ColorWhite() {
		quad[0].Color = sf.ColorGreen()
		quad[1].Color = sf.ColorGreen()
		quad[2].Color = sf.ColorGreen()
		quad[3].Color = sf.ColorGreen()
	} else {
		quad[0].Color = sf.ColorWhite()
		quad[1].Color = sf.ColorWhite()
		quad[2].Color = sf.ColorWhite()
		quad[3].Color = sf.ColorWhite()
	}
}

func (t *Tray) QuadFromIndexes(index_x, index_y int) ([]sf.Vertex, bool) {
	if index_x > t.Columns -1 || index_x < 0 || index_y > t.Rows -1 || index_y < 0 {
		return nil, false
	}

	index := (index_y * t.Columns * 4) + (index_x * 4)
	return t.Grid.Vertices[index:index+4], true
}

func (t *Tray) IndexesFromCoords(pos sf.Vector2f) (int, int) {
	diff := pos.Minus(t.Origin)

	portions_x := diff.X / ((5 * GridSize) + 3 + 4)
	portions_y := diff.Y / ((5 * GridSize) + 3 + 4)

	coord_xf := (diff.X - (3 * portions_x) - (4 * portions_x)) / 15
	coord_yf := (diff.Y - (3 * portions_y) - (4 * portions_y)) / 15

	return int(math.Floor(float64(coord_xf))), int(math.Floor(float64(coord_yf)))
}

func (t *Tray) Populate () {
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
		BlendMode: sf.BlendAlpha,
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
