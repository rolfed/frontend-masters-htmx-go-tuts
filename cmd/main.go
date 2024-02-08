package main

import (
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Templates struct {
    tmpl *template.Template
}

func newTemplate() *Templates {
    return &Templates{
        tmpl: template.Must(template.ParseGlob("views/*.html")),
    }
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}

type Count struct {
    Count int
}

var id = 0
type Contact struct {
    Name string
    Email string
    Id int
} 

func newContact(name string, email string) Contact {
    id++
    return Contact{
        Name: name,
        Email: email,
        Id: id,
    }
}

type Contacts = []Contact

func (d *Data) indexOf(id int) int {
    for i, contact := range d.Contacts {
        if contact.Id == id {
            return i 
        }
    }
    return -1  
}

func (d *Data) hasEmail(email string) bool {
    hasEmail := false
    for _, contact := range d.Contacts {
        if contact.Email == email {
            hasEmail = true
        }
    }
    return hasEmail 
}

type Data struct {
    Contacts Contacts 
}

func newData() *Data {
    return &Data{
        Contacts: []Contact{
            newContact("jk", "jk"),
            newContact("Jane Doe", "janed@gmail.comg"),
        },
    }
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() *FormData {
    return &FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

type Page struct {
    Data Data
    Form FormData
}

func newPage() Page {
    return Page{
        Data: *newData(),
        Form: *newFormData(),
    }
}

func main() {
    e := echo.New()
    page := newPage()
    e.Renderer = newTemplate()

    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format:  "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human} ${error}\n",
        Output: e.Logger.Output(),
    }))

    e.Static("/css", "css")
    e.Static("/images", "images")

    // Handler
    e.GET("/", func(c echo.Context) error {
        return c.Render(http.StatusOK, "index", page)
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists"
            return c.Render(http.StatusUnprocessableEntity, "form", formData)
        }

        contact := newContact(name, email)
        page.Data.Contacts = append(page.Data.Contacts, contact)

        c.Render(200, "form", newFormData())
        return c.Render(http.StatusOK, "out-band-contact", contact)
    })

    e.DELETE("/contact/:id", func(c echo.Context) error {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            return c.String(http.StatusBadRequest, "Invalid id")
        }

        index := page.Data.indexOf(id)
        if index == -1 {
            return c.String(http.StatusNotFound, "Contact not found")
        }

        page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)

        return c.HTML(http.StatusOK, "")
    })

    // Start server
    e.Logger.Fatal(e.Start(":42069"))
}
