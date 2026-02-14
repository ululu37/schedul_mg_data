package router

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"
	"strings"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ledongthuc/pdf"
	"github.com/nfnt/resize"
)

func extractTextFromPDF(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	r, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		// If GetPlainText fails, try reading page by page
		numPages := r.NumPage()
		for i := 1; i <= numPages; i++ {
			p := r.Page(i)
			if p.V.IsNull() {
				continue
			}
			s, _ := p.GetPlainText(nil)
			buf.WriteString(s)
		}
		return buf.String(), nil
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

var precurriculumUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func downscaleImage(data []byte) ([]byte, string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	// Downscale to max width 1600px, keeping aspect ratio
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	if width > 1600 {
		img = resize.Resize(1600, 0, img, resize.Lanczos3)
	}

	var buf bytes.Buffer
	// Encode as JPEG with 80% quality to save space
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/jpeg", nil
}

func RegisterPreCurriculumRoutes(e *echo.Echo, uc *usecase.PreCurriculum, importUc *usecase.ImportPrecuriculum) {
	g := e.Group("/precurriculum")

	g.POST("/import", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			// Check if it's text instead
			text := c.FormValue("text")
			if text == "" {
				return c.JSON(http.StatusBadRequest, "no file or text provided")
			}
			id, err := importUc.Import(map[string]interface{}{
				"type": "text",
				"text": text,
			}, nil)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, map[string]any{"id": id})
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		contentType := http.DetectContentType(data)
		var input interface{}

		if contentType == "application/pdf" {
			text, err := extractTextFromPDF(data)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to extract PDF text: %v", err))
			}
			input = map[string]interface{}{
				"type": "text",
				"text": text,
			}
		} else if contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/webp" || contentType == "image/gif" {
			// Auto Resize
			newData, newType, err := downscaleImage(data)
			if err == nil {
				data = newData
				contentType = newType
			}
			base64Data := base64.StdEncoding.EncodeToString(data)
			input = map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": fmt.Sprintf("data:%s;base64,%s", contentType, base64Data),
				},
			}
		} else {
			input = map[string]interface{}{
				"type": "text",
				"text": string(data),
			}
		}

		id, err := importUc.Import(input, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"id": id})
	}, middleware.Permit(0))

	g.GET("/import/ws", func(c echo.Context) error {
		ws, err := precurriculumUpgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// For WebSocket import, we expect the client to send the text content or base64 image as the first message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return err
		}

		var input interface{}
		// Basic check if it's likely an image or text
		if strings.HasPrefix(string(msg), "data:image") {
			// Extract base64 part
			parts := strings.Split(string(msg), ",")
			if len(parts) > 1 {
				decoded, err := base64.StdEncoding.DecodeString(parts[1])
				if err == nil {
					newData, _, err := downscaleImage(decoded)
					if err == nil {
						// Re-encode resized image
						newBase64 := base64.StdEncoding.EncodeToString(newData)
						msg = []byte(fmt.Sprintf("data:image/jpeg;base64,%s", newBase64))
					}
				}
			}
			input = map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": string(msg),
				},
			}
		} else {
			input = map[string]interface{}{
				"type": "text",
				"text": string(msg),
			}
		}

		progress := make(chan string)
		done := make(chan uint)
		errChan := make(chan error)

		go func() {
			id, err := importUc.Import(input, progress)
			if err != nil {
				errChan <- err
				return
			}
			done <- id
		}()

		for {
			select {
			case p := <-progress:
				ws.WriteMessage(websocket.TextMessage, []byte(p))
			case id := <-done:
				ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SUCCESS: %d", id)))
				return nil
			case err := <-errChan:
				ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ERROR: %v", err)))
				return nil
			}
		}
	})

	g.POST("/import/webhook", func(c echo.Context) error {
		var req struct {
			CallbackURL string `json:"callback_url"`
			Text        string `json:"text"`
			ImageBase64 string `json:"image_base64"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if req.CallbackURL == "" {
			return c.JSON(http.StatusBadRequest, "callback_url is required")
		}

		var input interface{}
		if req.ImageBase64 != "" {
			// Extract base64 part for resizing if it's a data URL
			imgData := req.ImageBase64
			if strings.HasPrefix(imgData, "data:image") {
				parts := strings.Split(imgData, ",")
				if len(parts) > 1 {
					decoded, err := base64.StdEncoding.DecodeString(parts[1])
					if err == nil {
						newData, _, err := downscaleImage(decoded)
						if err == nil {
							newBase64 := base64.StdEncoding.EncodeToString(newData)
							imgData = fmt.Sprintf("data:image/jpeg;base64,%s", newBase64)
						}
					}
				}
			}
			input = map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": imgData,
				},
			}
		} else {
			input = map[string]interface{}{
				"type": "text",
				"text": req.Text,
			}
		}

		go func() {
			id, err := importUc.Import(input, nil)
			result := map[string]interface{}{
				"status": "success",
				"id":     id,
			}
			if err != nil {
				result["status"] = "error"
				result["error"] = err.Error()
			}

			jsonData, _ := json.Marshal(result)
			http.Post(req.CallbackURL, "application/json", strings.NewReader(string(jsonData)))
		}()

		return c.JSON(http.StatusAccepted, map[string]string{"message": "Import started, result will be sent to callback_url"})
	}, middleware.Permit(0))

	g.POST("", func(c echo.Context) error {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		id, err := uc.Create(req.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"id": id})
	}, middleware.Permit(0))

	g.GET("", func(c echo.Context) error {
		search := c.QueryParam("search")
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
		list, count, err := uc.Listing(search, page, perPage)
		fmt.Println(list)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": list, "count": count})
	}, middleware.Permit(0, 1))

	g.GET("/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		data, err := uc.GetByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, data)
	}, middleware.Permit(0, 1))

	g.PUT("/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var req struct {
			Name string `json:"name"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		updated, err := uc.Update(uint(id), req.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, updated)
	}, middleware.Permit(0))

	g.DELETE("/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if err := uc.Delete(uint(id)); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.POST("/:id/subject", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var req []struct {
			SubjectName string `json:"subject_name"`
			Credit      int    `json:"credit"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fmt.Println("id:", id)
		var newSubjectInCurriculum []entities.SubjectInPreCurriculum
		for _, r := range req {
			newSubjectInCurriculum = append(newSubjectInCurriculum, entities.SubjectInPreCurriculum{
				PreCurriculumID: uint(id),
				Subject:         entities.Subject{Name: r.SubjectName},
				Credit:          r.Credit,
			})
		}

		err := uc.CreateSubject(uint(id), newSubjectInCurriculum)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.DELETE("/subject/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err := uc.RemoveSubject(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.PUT("/subject/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var req struct {
			SubjectName string `json:"subject_name"`
			Credit      int    `json:"credit"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := uc.UpdateSubject(uint(id), req.SubjectName, req.Credit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))
}
