package routes

import (
	"html/template"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"github.com/tdewolff/minify/v2/minify"
)

func RegisterDocs(g *gin.Engine, version string) {
	// Create the style document and cover image
	style := func() string {
		b, err := os.ReadFile("style.css")
		if err != nil {
			log.Fatal("could not read stylesheet")
		}

		min, err := minify.CSS(string(b))
		if err != nil {
			log.Fatal("could not minify CSS file")
		}

		return min
	}()
	createCover := func() {
		jpg, err := gg.LoadJPG("cover.jpg")
		if err != nil {
			log.Fatal("could not open cover image source")
		}

		img := gg.NewContextForImage(jpg)

		notoBytes, err := os.ReadFile("assets/NotoSans-Bold.ttf")
		if err != nil {
			log.Fatal("could not read font file")
		}
		noto, err := truetype.Parse(notoBytes)
		if err != nil {
			log.Fatal("could not parse font as truetype")
		}

		img.SetFontFace(truetype.NewFace(noto, &truetype.Options{Size: 28}))
		img.SetColor(color.NRGBA{35, 202, 255, 255})
		strWidth, _ := img.MeasureString("v" + version)
		img.DrawString("v"+version, float64(jpg.Bounds().Dx()/2)-strWidth/2, 218)

		// Write file
		f, err := os.Create("assets/cover.jpg")
		if err != nil {
			log.Fatal("cannot create cover image")
		}

		if err := jpeg.Encode(f, img.Image(), nil); err != nil {
			log.Fatal("cannot save image")
		}
	}
	createCover()

	// Load our files and create a CSS escaping function
	g.SetFuncMap(template.FuncMap{
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
	})
	g.LoadHTMLFiles("docs.gohtml")

	g.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docs.gohtml", gin.H{
			"style":   style,
			"version": version,
		})
	})

	g.Static("/assets", "./assets")
}
