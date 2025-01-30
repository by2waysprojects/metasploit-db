package services

type HTTPPacket struct {
	Seq       string     `json:"seq"`
	Size      string     `json:"size"`
	Protocol  string     `json:"protocol"`
	Body      []Body     `json:"body"`
	Cookie    []Cookie   `json:"cookie"`
	Header    []Header   `json:"header"`
	Uri       []Uri      `json:"uri"`
	Method    []Method   `json:"method"`
	Wildcards []Wildcard `json:"wildcards"`
}

type Body struct {
	Data string `json:"data"`
}

type Cookie struct {
	Data string `json:"data"`
}

type Header struct {
	Data string `json:"data"`
}

type Uri struct {
	Data  string `json:"data"`
	Exact bool   `json:"exact"`
}

type Method struct {
	Data string `json:"data"`
}

type Wildcard struct {
	Data string `json:"data"`
}
