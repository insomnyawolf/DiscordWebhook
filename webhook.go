package webhook

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Identity Stores Webhook identity
type Identity struct {
	Username   string
	AvatarURL  string
	WebHookURL string
}

//WebHook returns WebHook data structure
func (c *Identity) WebHook() Hook {
	return Hook{
		avatarURL:  c.AvatarURL,
		username:   c.Username,
		webHookURL: c.WebHookURL,
	}
}

//Hook Stores data for the webhook request
type Hook struct {
	username   string //override the default username of the webhook	false
	avatarURL  string //override the default avatar of the webhook	false
	webHookURL string
	Content    string //the message contents (up to 2000 characters)	one of content, file, embeds
	Tts        bool   //true if this is a TTS message	false
	Embedded   EMBEDDED
	//file ???
}

//Close EOL
func (c *Hook) Close() {
	c = nil
}

//SendWebHook Sends request to the api endpoint
func (c *Hook) SendWebHook() string {
	jsonStr := []byte(c.JSONRequest())
	req, err := http.NewRequest("POST", c.webHookURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp.Status
}

//JSONRequest Returns the json that would be sent to api endpoint
func (c *Hook) JSONRequest() string {
	var req string
	if c.username != "" {
		if req != "" {
			req += ","
		}
		req += fmt.Sprintf(`"username":"%v"`, c.username)
	}
	if c.avatarURL != "" {
		if req != "" {
			req += ","
		}
		req += fmt.Sprintf(`"avatar_url":"%v"`, c.avatarURL)
	}
	if c.Content != "" {
		if req != "" {
			req += ","
		}
		req += fmt.Sprintf(`"content":"%v"`, c.Content)
	}
	if c.Tts != false {
		if req != "" {
			req += ","
		}
		req += fmt.Sprintf(`"tts": "%v"`, c.Tts)
	}
	emb := c.Embedded.Parse()
	if emb != "" {
		if req != "" {
			req += ","
		}
		req += emb
	}
	return fmt.Sprintf(`{ %v }`, req)
}

//EMBEDDED Data
type EMBEDDED struct {
	Color       COLOR
	Author      AUTHOR
	Title       TITLE
	URL         URL
	Description DESCRIPTION
	Fields      []FIELDS
	Image       IMAGE
	Thumbnail   THUMBNAIL
	Footer      FOOTER
	Timestamp   TIMESTAMP
}

//COLOR			√
//AUTHOR		√
//TITLE			√
//URL			√
//DESCRIPTION	√
//FIELDS		√
//IMAGE 		√
//THUMBNAIL 	√
//FOOTER 		√
//TIMESTAMP 	√

func (c *EMBEDDED) parseFields() string {
	var Fields string
	for _, v := range c.Fields {
		Field := v.Parse()
		if Field != "" {
			if Fields != "" {
				Fields += ","
			}
			Fields += Field
		}
	}
	if Fields == "" {
		return ""
	}
	return fmt.Sprintf(`"fields": [ %v ]`, Fields)
}

//Parse converts EMBEDDED into the equivalent json
func (c *EMBEDDED) Parse() string {
	var emb string
	Fields := c.parseFields()
	if Fields != "" {
		if emb != "" {
			emb += ","
		}
		emb += Fields
	}
	Desc := c.Description.Parse()
	if Desc != "" {
		if emb != "" {
			emb += ","
		}
		emb += Desc
	}
	Title := c.Title.Parse()
	if Title != "" {
		if emb != "" {
			emb += ","
		}
		emb += Title
	}
	aut := c.Author.Parse()
	if aut != "" {
		if emb != "" {
			emb += ","
		}
		emb += aut
	}
	col := c.Color.Parse()
	if col != "" {
		if emb != "" {
			emb += ","
		}
		emb += col
	}
	URI := c.URL.Parse()
	if URI != "" {
		if emb != "" {
			emb += ","
		}
		emb += URI
	}
	img := c.Image.Parse()
	if img != "" {
		if emb != "" {
			emb += ","
		}
		emb += img
	}
	thumb := c.Thumbnail.Parse()
	if thumb != "" {
		if emb != "" {
			emb += ","
		}
		emb += thumb
	}
	foot := c.Footer.Parse()
	if foot != "" {
		if emb != "" {
			emb += ","
		}
		emb += foot
	}
	time := c.Timestamp.Parse()
	if time != "" {
		if emb != "" {
			emb += ","
		}
		emb += time
	}
	if emb == "" {
		return ""
	}
	return fmt.Sprintf(`"embeds": [{%v}]`, emb)
}

//COLOR Hexdecimal Color
type COLOR struct {
	Color string
}

//Parse Used By EMBEDDED
func (c *COLOR) Parse() string {
	Empty := COLOR{}
	if *c != Empty {
		n, err := strconv.ParseInt(c.Color, 16, 25)
		if err != nil {
			log.Println(err)
		}
		return fmt.Sprintf(`"color": "%v" `, n)
	}
	return ""
}

//AUTHOR  Holds dara for Author data structure
type AUTHOR struct {
	Name    string
	URL     string
	IconURL string
}

//Parse Used By EMBEDDED
func (c *AUTHOR) Parse() string {
	Empty := AUTHOR{}
	var auth string
	if *c != Empty {
		if c.Name != "" {
			auth += fmt.Sprintf(`"name": "%v" `, c.Name)
		}
		if auth != "" {
			if c.URL != "" {
				auth += fmt.Sprintf(`, "url": "%v" `, c.URL)
			}
			if c.IconURL != "" {
				auth += fmt.Sprintf(`, "icon_url": "%v" `, c.IconURL)
			}
		}
		return fmt.Sprintf(`"author": { %v }`, auth)
	}
	return ""
}

//TITLE of EMBEDDED data
type TITLE struct {
	Title string
}

//Parse Used By EMBEDDED
func (c *TITLE) Parse() string {
	Empty := TITLE{}
	if *c != Empty {
		return fmt.Sprintf(`"title": "%v"`, c.Title)
	}
	return ""
}

//URL in the EMBEDDED data
type URL struct {
	Title string
	URL   string
}

//Parse Used By EMBEDDED
func (c *URL) Parse() string {
	Empty := URL{}
	if *c != Empty {
		if c.Title == "" {
			c.Title = c.URL
		}
		return fmt.Sprintf(`"title": "%v", "url": "%v" `, c.Title, c.URL)
	}
	return ""
}

//DESCRIPTION Holds description data
type DESCRIPTION struct {
	Description string
}

//Parse Used By EMBEDDED
func (c *DESCRIPTION) Parse() string {
	Empty := DESCRIPTION{}
	if *c != Empty {
		return fmt.Sprintf(` "description": "%v" `, c.Description)
	}
	return ""
}

//FIELDS Holds fields data
type FIELDS struct {
	Name   string
	Value  string
	Inline bool
}

//Parse Used By EMBEDDED
func (c *FIELDS) Parse() string {
	Empty := FIELDS{}
	var Field string
	if *c != Empty {
		if c.Name != "" {
			Field += fmt.Sprintf(`"name": "%v" `, c.Name)
		}
		if c.Value != "" {
			if Field != "" {
				Field += ","
			}
			Field += fmt.Sprintf(`"value": "%v" `, c.Value)
		}
		if c.Inline != false {
			if Field != "" {
				Field += ","
			}
			Field += fmt.Sprintf(`"inline": "%v" `, c.Inline)
		}
		return fmt.Sprintf("{ %v }", Field)
	}
	return ""
}

//IMAGE of the EMBEDDED
type IMAGE struct {
	URL string
}

//Parse Used By EMBEDDED
func (c *IMAGE) Parse() string {
	Empty := IMAGE{}
	if *c != Empty {
		return fmt.Sprintf(`"image": { "url": "%v" }`, c.URL)
	}
	return ""
}

//THUMBNAIL of the EMBEDDED
type THUMBNAIL struct {
	URL string
}

//Parse Used By EMBEDDED
func (c *THUMBNAIL) Parse() string {
	Empty := THUMBNAIL{}
	if *c != Empty {
		return fmt.Sprintf(`"thumbnail": { "url": "%v" }`, c.URL)
	}
	return ""
}

//FOOTER Data shown at the bottom od the EMBEDDED
type FOOTER struct {
	Text    string
	IconURL string
}

//Parse Used By EMBEDDED
func (c *FOOTER) Parse() string {
	var foot string
	Empty := FOOTER{}
	if *c != Empty {
		if c.Text != "" {
			foot += fmt.Sprintf(`"text": "%v"`, c.Text)
		}
		if c.IconURL != "" {
			if foot != "" {
				foot += ","
			}
			foot += fmt.Sprintf(`"icon_url": "%v"`, c.IconURL)
		}
		return fmt.Sprintf(`"footer": {%v}`, foot)
	}
	return ""
}

//TIMESTAMP Stores date value in the message
type TIMESTAMP struct {
	Time string
}

//Parse Used By EMBEDDED
func (c *TIMESTAMP) Parse() string {
	Empty := TIMESTAMP{}
	if *c != Empty {
		t, err := time.Parse("2006-01-02 03:04:05.0000", c.Time)
		if err != nil {
			log.Println(err)
		}
		return fmt.Sprintf(`"timestamp": "%v"`, t.Format("2006-01-02 03:04:05.0000"))
	}
	return ""
}
