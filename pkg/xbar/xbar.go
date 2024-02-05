package xbar

import (
	"fmt"
	"strconv"
	"strings"
)

type Xbar struct {
	Line    Line
	SubLine []Xbar
}

type Line struct {
	Title         string
	Href          *string
	Color         *string
	Font          *string
	Size          *int
	Shell         *string
	Terminal      *bool
	Refresh       *bool
	Dropdown      *bool
	Length        *int
	Trim          *bool
	TemplateImage *string
	Image         *string
}

var SeparateLine = Xbar{
	Line: Line{
		Title: "---",
	},
}

func (x *Xbar) Print() {
	fmt.Println(x.Line.ParamsString())

	for _, sub := range x.SubLine {
		fmt.Printf("-- %s\n", sub.Line.ParamsString())
	}
}

func (l *Line) ParamsString() string {
	params := make([]string, 0, 13)
	params = append(params, l.Title)

	if v := l.Href; v != nil {
		params = append(params, "href="+*v)
	}
	if v := l.Color; v != nil {
		params = append(params, "color="+*v)
	}
	if v := l.Font; v != nil {
		params = append(params, "font="+*v)
	}
	if v := l.Size; v != nil {
		params = append(params, "size="+strconv.Itoa(*v))
	}
	if v := l.Shell; v != nil {
		params = append(params, "shell="+*v)
	}
	if v := l.Terminal; v != nil {
		params = append(params, "terminal="+strconv.FormatBool(*v))
	}
	if v := l.Refresh; v != nil {
		params = append(params, "refresh="+strconv.FormatBool(*v))
	}
	if v := l.Dropdown; v != nil {
		params = append(params, "dropdown="+strconv.FormatBool(*v))
	}
	if v := l.Length; v != nil {
		params = append(params, "length="+strconv.Itoa(*v))
	}
	if v := l.Trim; v != nil {
		params = append(params, "trim="+strconv.FormatBool(*v))
	}
	if v := l.TemplateImage; v != nil {
		params = append(params, "templateImage="+*v)
	}
	if v := l.Image; v != nil {
		params = append(params, "image="+*v)
	}

	return strings.Join(params, " | ")
}
