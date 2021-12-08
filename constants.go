package mixpanel

var (
	Headers = struct {
		Accept      string
		ContentType string
	}{
		Accept:      "Accept",
		ContentType: "Content-Type",
	}

	MimeTypes = struct {
		TextPlain          string
		ApplicationJSON    string
		XWWWFormURLEncoded string
	}{
		TextPlain:          "text/plain",
		ApplicationJSON:    "application/json",
		XWWWFormURLEncoded: "application/x-www-form-urlencoded",
	}
)
