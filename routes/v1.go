package routes

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype/truetype"
	"github.com/tdewolff/minify/v2/minify"
	"html/template"
	"image/color"
	"image/jpeg"
	"log"
	"mc-api/cache"
	"mc-api/minecraft"
	"net/http"
	"os"
)

func RegisterV1Routes(g *gin.Engine, ch *cache.Cache, version string) {
	g.GET("/v1/ping/:ip", func(c *gin.Context) {
		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		redisResponse, err := ch.GetServer(c, address.Combined)

		if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			tcpServer, terr := minecraft.PingServer(address)
			if terr != nil {
				fmt.Println(terr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": terr.Error()})
				return
			}

			if err := ch.SetServer(c, address.Combined, tcpServer); err != nil {
				fmt.Println("failed to cache")
			}

			server, err := minecraft.CreateResponseObj(tcpServer, address.IPv4, address.Port, false)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, server)
			return
		} else if err != nil {
			fmt.Println("failed to retrieve cache")
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		server, err := minecraft.CreateResponseObj(redisResponse, address.IPv4, address.Port, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, server)
	})
	g.GET("/v1/icon/:ip", func(c *gin.Context) {
		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cached, err := ch.GetServer(c, address.Combined)
		if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			server, err := minecraft.PingServer(address)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			icon, err := minecraft.GetIcon(server)
			if err != nil {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/png", icon)
			return
		} else if err != nil {
			fmt.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		icon, err := minecraft.GetIcon(cached)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "image/png", icon)
	})

	// Docs routes
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
		noto, err := truetype.Parse(notoBytes)

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

	g.SetFuncMap(template.FuncMap{
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
	})
	g.LoadHTMLFiles("docs.gohtml")

	g.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docs.gohtml", gin.H{
			"style":   style,
			"version": "1.1.1",
		})
	})

	g.Static("/assets", "./assets")

	// fallback redirect
	g.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/docs")
	})
}
