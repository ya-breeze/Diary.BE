package utils

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type imagePrefixRenderer struct {
	html.Renderer
	ImagePrefix string
}

func NewImagePrefixRenderer(imagePrefix string) *imagePrefixRenderer {
	return &imagePrefixRenderer{
		ImagePrefix: imagePrefix,
	}
}

func (r *imagePrefixRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if img, ok := node.(*ast.Image); ok {
		if entering {
			// Change the image URL
			newSrc := fmt.Sprintf("%s%s", r.ImagePrefix, string(img.Destination))
			_, _ = fmt.Fprintf(w, `<br><a href="%s"><img src="%s" alt="%s"`, newSrc, newSrc, img.Title)
		} else {
			// Close the image tag
			_, _ = io.WriteString(w, "/></a><br>")
		}
		return ast.SkipChildren
	}

	// fallback to default
	return r.Renderer.RenderNode(w, node, entering)
}
