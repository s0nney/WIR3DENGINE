package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"WIR3DENGINE/model"

	"WIR3DENGINE/utils"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	pagination "github.com/webstradev/gin-pagination"
	"gopkg.in/gographics/imagick.v2/imagick"
	_ "gopkg.in/gographics/imagick.v2/imagick"
)

func SpinUpRoutes(router *gin.Engine) {
	paginator := pagination.Default()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"headerTags": template.HTML(model.GenerateHead()),
			"foot":       template.HTML(model.GenerateFoot()),
			"title":      "TEMPORAL@WORLD",
			"Header1":    "The TEMPORAL@WORLD is best viewed in a desktop browser.",
		})
	})

	router.GET("/index", func(c *gin.Context) {

		presentRows := 20
		rows, err := utils.Db.Query("SELECT in_name, in_text FROM posts ORDER BY id DESC LIMIT $1", presentRows)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		defer rows.Close()

		var results []map[string]interface{}

		columns, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))

			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			rowData := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				rowData[col] = val
			}

			results = append(results, rowData)
		}

		var resultString string
		for _, result := range results {
			resultString += fmt.Sprintf(
				"\n<p class='post' style='color:%s; font-family: %s; %s'> %s: %s</p>\n",
				model.ColorCalibrator(),
				model.GetFont(8),
				model.GenerateRandomStyles(),
				result["in_name"],
				model.Corrupt(result["in_text"].(string)),
			)
		}

		imageChecksumDb, err := model.GetRandomImagesFromDB(5)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var imagesHTML string

		for _, imageChecksum := range imageChecksumDb {
			imageChecksum := imageChecksum[14:]
			imagesHTML += fmt.Sprintf("\n<img style='position:absolute; z-index: %s; opacity:%s; left:%s; top:%s;' src='images/%v'>\n",
				model.GetZIndex(),
				model.GetOpacity(),
				model.GetRandomImageLeft(),
				model.GetRandomImageTop(),
				imageChecksum)
		}

		token := csrf.GetToken(c)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"headerTags":    template.HTML(model.GenerateHead()),
			"foot":          template.HTML(model.GenerateFoot()),
			"title":         "TEMPORAL@WORLD",
			"posts":         template.HTML(resultString),
			"images":        template.HTML(imagesHTML),
			"requiredInput": "required",
			"CSRFToken":     token, // Anti-CSRF Implementation
		})
	})

	router.GET("/gallery", paginator, func(c *gin.Context) {

		const pageSize = 10

		page := c.GetInt("page")
		if page == 0 {
			c.HTML(http.StatusOK, "err404.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Resource not found.",
			})
			return
		}

		imageChecksumDb, err := model.GetRecentImagesFromDB(100)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Invalid argument.",
			})
			return
		}

		startIndex := (page - 1) * pageSize
		endIndex := startIndex + pageSize
		if startIndex >= len(imageChecksumDb) {
			c.HTML(http.StatusOK, "err404.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Resource not found.",
			})
			return
		}

		if endIndex > len(imageChecksumDb) {
			endIndex = len(imageChecksumDb)
		}
		paginatedImages := imageChecksumDb[startIndex:endIndex]

		var imagesHTML string

		for _, imageChecksum := range paginatedImages {
			imageChecksum := imageChecksum[14:]
			imagesHTML += fmt.Sprintf("\n<img class='pic' src='images/%v'>\n",
				imageChecksum)
		}

		totalPages := int(math.Ceil(float64(len(imageChecksumDb)) / float64(pageSize)))

		token := csrf.GetToken(c)
		c.HTML(http.StatusOK, "gallery.html", gin.H{
			"headerTags": template.HTML(model.GenerateHead()),
			"foot":       template.HTML(model.GenerateFoot()),
			"title":      "TEMPORAL@WORLD",
			"CSRFToken":  token, // Anti-CSRF Implementation
			"images":     template.HTML(imagesHTML),
			"page":       page,
			"total":      totalPages,
			"next":       page + 1,
			"previous":   utils.GetPreviousPage(page),
		})
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{
			"headerTags": template.HTML(model.GenerateHead()),
			"foot":       template.HTML(model.GenerateFoot()),
			"title":      "TEMPORAL@WORLD",
		})
	})

	router.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.html", gin.H{
			"headerTags": template.HTML(model.GenerateHead()),
			"foot":       template.HTML(model.GenerateFoot()),
			"title":      "TEMPORAL@WORLD",
		})
	})

	router.GET("/engine", func(c *gin.Context) {
		c.HTML(http.StatusOK, "engine.html", gin.H{
			"headerTags": template.HTML(model.GenerateHead()),
			"foot":       template.HTML(model.GenerateFoot()),
			"title":      "TEMPORAL@WORLD",
		})
	})

	// Basic API
	router.GET("/posts", func(c *gin.Context) {
		presentRows := 20
		rows, err := utils.Db.Query("SELECT id, in_name, in_text, date_posted FROM posts ORDER BY id DESC LIMIT $1", presentRows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []map[string]interface{}

		for rows.Next() {
			var column1, column2, column3, column4 string
			err := rows.Scan(&column1, &column2, &column3, &column4)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			data := map[string]interface{}{
				"id":          column1,
				"in_name":     column2,
				"in_text":     column3,
				"date_posted": column4[:10],
			}
			result = append(result, data)
		}

		c.JSON(http.StatusOK, result)
	})

	// GET n posts
	router.GET("/posts/:num", func(c *gin.Context) {
		numParam := c.Param("num")
		numRows, err := strconv.Atoi(numParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
			return
		}

		if numRows > 100 {
			numRows = 100
		}

		rows, err := utils.Db.Query("SELECT id, in_name, in_text, date_posted FROM posts ORDER BY id DESC LIMIT $1", numRows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []map[string]interface{}

		for rows.Next() {
			var column1, column2, column3, column4 string
			err := rows.Scan(&column1, &column2, &column3, &column4)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			data := map[string]interface{}{
				"id":          column1,
				"in_name":     column2,
				"in_text":     column3,
				"date_posted": column4[:10],
			}
			result = append(result, data)
		}

		c.JSON(http.StatusOK, result)
	})

	router.POST("/in_text", model.NewRateLimiter().Middleware(), func(c *gin.Context) {

		clientIP := c.ClientIP()

		if model.IsBanned(clientIP) {
			c.HTML(http.StatusBadRequest, "banPage.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
			})
			return
		}

		var formData model.InText

		if err := c.ShouldBind(&formData); err != nil {
			c.HTML(http.StatusBadRequest, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Turn back now.",
			})
			return
		}

		if model.IsWhitespace(formData.TextData) {
			c.HTML(http.StatusBadRequest, "err2.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Write something!",
			})
			return
		}

		if len(formData.TextData) < 1 {
			c.HTML(http.StatusBadRequest, "err2.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Write something!",
			})
			return
		}

		if len(formData.NameData) > 62 {
			c.HTML(http.StatusBadRequest, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Name length exceeds maximum length of 62 characters.",
			})
			return
		}

		if len(formData.TextData) > 512 {
			c.HTML(http.StatusBadRequest, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Message length exceeds maximum length of 512 characters.",
			})
			return
		}

		in_name := strings.TrimSpace(c.PostForm("name"))

		if in_name == "" {
			in_name = "Anonymous"
		}

		sanitized_name := template.HTMLEscapeString(in_name)

		in_text := strings.TrimSpace(c.PostForm("in_text"))
		sanitized_message := template.HTMLEscapeString(in_text)

		_, err := utils.Db.Exec("INSERT INTO posts (in_name, in_text, ip) VALUES ($1, $2, $3)", sanitized_name, sanitized_message, clientIP)
		if err != nil {
			c.HTML(http.StatusBadRequest, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Turn back now.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/index")
	})

	router.POST("/in_image", model.NewRateLimiter().Middleware(), func(c *gin.Context) {

		clientIP := c.ClientIP()

		if model.IsBanned(clientIP) {
			c.HTML(http.StatusBadRequest, "banPage.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
			})
			return
		}

		inImage, err := c.FormFile("in_image")
		if err != nil {
			c.HTML(http.StatusBadRequest, "err.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "No file provided.",
			})
			return
		}

		allowedTypes := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
		}

		ext := filepath.Ext(inImage.Filename)
		if !allowedTypes[ext] {
			c.HTML(http.StatusBadRequest, "fileType.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
			})
			return
		}

		maxSize := int64(2 << 20) // 2 MB
		if inImage.Size > maxSize {
			c.HTML(http.StatusBadRequest, "fileSize.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
			})
			return
		}

		uploadChecksum, err := model.ReadChecksum(inImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		outputFilename := "public/images/" + model.TrimFn(model.FnMixer(inImage.Filename), 8)
		var existingChecksum string
		err = utils.Db.QueryRow("SELECT checksum FROM images WHERE checksum = $1", uploadChecksum).Scan(&existingChecksum)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = utils.Db.Exec("INSERT INTO images (checksum, ip, filename) VALUES ($1, $2, $3)", uploadChecksum, clientIP, outputFilename)
				if err != nil {
					c.HTML(http.StatusBadRequest, "err404.html", gin.H{
						"headerTags": template.HTML(model.GenerateHead()),
						"foot":       template.HTML(model.GenerateFoot()),
						"title":      "TEMPORAL@WORLD",
						"error":      "Not found or something went wrong.",
					})
					return
				}
			} else {
				c.HTML(http.StatusInternalServerError, "err.html", gin.H{
					"headerTags": template.HTML(model.GenerateHead()),
					"foot":       template.HTML(model.GenerateFoot()),
					"title":      "TEMPORAL@WORLD",
					"error":      "Something went wrong.",
				})
				return
			}
		} else {
			c.HTML(http.StatusBadRequest, "duplicate.html", gin.H{
				"headerTags": template.HTML(model.GenerateHead()),
				"foot":       template.HTML(model.GenerateFoot()),
				"title":      "TEMPORAL@WORLD",
				"error":      "Image is too similar to previous images.",
			})
			return
		}

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		savePath := filepath.Join(currentDir, "/tmp")

		if err := os.MkdirAll(savePath, 0755); err != nil {
			log.Fatal(err)
		}

		filePath := filepath.Join(savePath, model.TrimFn(model.FnMixer(inImage.Filename), 8))

		if err := c.SaveUploadedFile(inImage, filePath); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		// image processing here
		src, err := os.Open(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to open the image file"})
			return
		}
		defer src.Close()

		mw := imagick.NewMagickWand()
		defer mw.Destroy()

		if err := mw.ReadImageFile(src); err != nil {
			c.JSON(500, gin.H{"error": "Unable to read the image"})
			return
		}

		imageCol := model.GetRandomSizeValue(200, 300)
		imageRow := model.GetRandomSizeValue(150, 250)

		mw.ResizeImage(imageCol, imageRow, imagick.FILTER_LANCZOS, 1)
		mw.QuantizeImage(50, imagick.COLORSPACE_GRAY, 0, true, true)

		fillColor := imagick.NewPixelWand()

		fillColor.SetColor("rgb(25, 14, 255)")
		mw.ColorizeImage(fillColor, fillColor)

		tintColor := imagick.NewPixelWand()

		tintColor.SetColor("rgb(25, 14, 255)")
		mw.TintImage(tintColor, fillColor)

		brightness := 100.0
		saturation := 70.0
		hue := 100.0

		mw.ModulateImage(brightness, saturation, hue)
		mw.PosterizeImage((20), true)

		if err := mw.WriteImage(outputFilename); err != nil {
			c.JSON(500, gin.H{"error": "Unable to save the processed image"})
			return
		}

		if err := os.Remove(filePath); err != nil {
			c.JSON(500, gin.H{"error": "Unable to delete the original image file"})
			return
		}
		// image processing here

		c.Redirect(http.StatusFound, "/index")
	})
}
