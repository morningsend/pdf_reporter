package report

type Report struct {
	Title         string
	CoverImageURL string
	Theme         ReportTheme
}

type ReportTheme struct {
	Colors Colors
	Font   string
}

type Colors struct {
}

type Typography struct {
	FontFamily string
}

type Charts struct {
}
