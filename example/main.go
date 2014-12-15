package main

import (
	"flag"
	"fmt"
	"net/http"

	. "gopkg.in/go-on/lib.v2/types"

	. "gopkg.in/go-on/lib.v2/html"
	"gopkg.in/go-on/lib.v2/internal/bootstrap/bs3"
	"gopkg.in/go-on/cdncache.v1"
)

var (
	port       = flag.Int("port", 8083, "port of the http server")
	mountPoint = flag.String("mountpoint", "", "mount point for the cached files (leave empty to prevent caching)")
)

func main() {
	flag.Parse()
	cdn := cdncache.CDN(*mountPoint)

	http.Handle("/",
		HTML5(
			Lang_("en"),
			HEAD(
				CharsetUtf8(),
				HttpEquiv("X-UA-Compatible", "IE=edge"),
				Viewport("width=device-width, initial-scale=1"),
				TITLE("Starter Template for Bootstrap"),
				CssHref(cdn(bs3.CDN_3_1_1_min)),
				CssHref(cdn(bs3.CDN_3_1_1_theme_min)),
				CssHref(cdn("http://getbootstrap.com/examples/starter-template/starter-template.css")),
				HTMLString(fmt.Sprintf(`
<!--[if lt IE 9]>
  <script src="%s"></script>
  <script src="%s"></script>
<![endif]-->
`,
					cdn("https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"),
					cdn("https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"),
				),
				),
			),
			BODY(
				DIV(bs3.Navbar, bs3.Navbar_inverse, bs3.Navbar_fixed_top,
					Attrs_("role", "navigation"),
					DIV(bs3.Container,
						DIV(bs3.Navbar_header,
							BUTTON(bs3.Navbar_toggle,
								Type_("button"),
								Attrs_("data-toggle", "collapse", "data-target", ".navbar-collapse"),
								SPAN(bs3.Sr_only, "Toogle navigation"),
								SPAN(bs3.Icon_bar),
								SPAN(bs3.Icon_bar),
								SPAN(bs3.Icon_bar),
							),
							AHref("#", bs3.Navbar_brand, "Project name"),
						),
						DIV(bs3.Collapse, bs3.Navbar_collapse,
							UL(bs3.Nav, bs3.Navbar_nav,
								LI(bs3.Active, AHref("#", "Home")),
								LI(AHref("#about", "About")),
								LI(AHref("#contact", "Contact")),
							),
						),
					),
				),
				DIV(bs3.Container,
					DIV(Class("starter-template"),
						H1("Bootstrap starter template"),
						P(bs3.Lead,
							"Use this document as a way to quickly start any new project.",
							BR(),
							"All you get is this text and a mostly barebones HTML document.",
						),
					),
				),

				JsSrc(cdn("//code.jquery.com/jquery-1.11.0.min.js")),
				JsSrc(cdn(bs3.CDN_3_1_1_js_min)),
			),
		),
	)
	fmt.Printf("listening on localhost: %d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
