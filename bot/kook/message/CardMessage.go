package message

type CardMessage struct {
	Type    string               `json:"type"`
	Theme   string               `json:"theme"`
	Size    string               `json:"size"`
	Color   string               `json:"color"`
	Modules []CardMessageModules `json:"modules"`
}

type CardMessageModules struct {
	Type     string `json:"type"`
	Text     Text   `json:"text"`
	Elements []struct {
		Type  string `json:"type"`
		Src   string `json:"src"`
		Theme string `json:"theme"`
		Value string `json:"value"`
		Text  struct {
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"text"`
	} `json:"elements"`
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  string `json:"size"`
}

type Text struct {
	Type    string  `json:"type"`
	Content string  `json:"content"`
	Cols    int     `json:"cols"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Elements struct {
	Type  string `json:"type"`
	Src   string `json:"src"`
	Theme string `json:"theme"`
	Value string `json:"value"`
	Text  Text   `json:"text"`
}
