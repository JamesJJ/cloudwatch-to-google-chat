package main

// https://developers.google.com/hangouts/chat/reference/message-formats/cards

type Message struct {
	Cards []MessageCard `json:"cards"`
}

type MessageCard struct {
	CardHeader CardHeader    `json:"header,omitempty"`
	Sections   []CardSection `json:"sections"`
}

type CardHeader struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}

type CardSection struct {
	Widgets []Widget `json:"widgets"`
}

type Widget struct {
	KeyValue      *WidgetKeyValue `json:"keyValue,omitempty"`
	TextParagraph *TextParagraph  `json:"textParagraph,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text"`
}

type WidgetKeyValue struct {
	TopLabel         string `json:"topLabel"`
	Content          string `json:"content"`
	ContentMultiline bool   `json:"contentMultiline"`
	BottomLabel      string `json:"bottomLabel,omitempty"`
	Icon             string `json:"icon,omitempty"`
	/*
	   "onClick": {
	        "openLink": {
	           "url": "https://example.com/"
	        }
	    },
	   "button": {
	       "textButton": {
	          "text": "VISIT WEBSITE",
	          "onClick": {
	              "openLink": {
	                   "url": "http://site.com"
	               }
	           }
	         }
	    }
	*/
}
