package nonograms

import (
	"fmt"
	sf "bitbucket.org/krepa098/gosfml2"
)

type TextureManager struct {
	Root string
	Textures map[string]*sf.Texture
}


func NewTextureManger(root string) *TextureManager {
	m := make(map[string]*sf.Texture)

	return &TextureManager{
		Root: root,
		Textures: m,
	}
}

func (t *TextureManager) Get(name string) *sf.Texture {
	texture, ok := t.Textures[name]

	if ok {
		return texture
	}

	texture, _ = sf.NewTextureFromFile(fmt.Sprintf("%s/textures/%s.png", t.Root, name), nil)
	t.Textures[name] = texture

	return texture
}
