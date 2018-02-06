package cmd

import (
	"html/template"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	_ "github.com/go-macaron/cache/nodb"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/base"
	"github.com/smolgu/lib/modules/middleware"
	"github.com/smolgu/lib/modules/setting"
	"github.com/smolgu/lib/routers"
	"github.com/smolgu/lib/routers/account"
	"github.com/smolgu/lib/routers/albums"
	"github.com/smolgu/lib/routers/auth"
	"github.com/smolgu/lib/routers/banner"
	"github.com/smolgu/lib/routers/book"
	"github.com/smolgu/lib/routers/link"
	"github.com/smolgu/lib/routers/menus"
	"github.com/smolgu/lib/routers/page"
	//"github.com/smolgu/lib/routers/ruslan"
	"github.com/smolgu/lib/routers/tag"
	"github.com/smolgu/lib/routers/users"
	"github.com/smolgu/lib/routers/vkr"
)

var CmdWeb = cli.Command{
	Name:  "web",
	Usage: "Start lib web server",
	Description: `lib web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		intFlag("port, p", 3000, "Temporary port number to prevent conflict"),
		stringFlag("config, c", "conf/app.ini", "Configuration file path (default ./conf/app.ini)"),
		stringFlag("mode, m", "dev", "Running mode"),
		stringFlag("storage_dir", "./", "Dirrectory with data"),
	},
}

func newMacaron(ctx *cli.Context) *macaron.Macaron {
	m := macaron.New()
	m.Use(gzip.Gziper())
	m.Use(macaron.Static("public", macaron.StaticOptions{SkipLogging: true, Expires: expires}))
	m.Use(macaron.Static("static", macaron.StaticOptions{SkipLogging: true, Prefix: setting.BuildHash, Expires: expires}))
	m.Use(macaron.Static("static", macaron.StaticOptions{SkipLogging: true, Prefix: "static", Expires: expires}))
	m.Use(macaron.Static(setting.DataDir, macaron.StaticOptions{
		SkipLogging: true,
		Prefix:      ""}))
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "data/sessions",
		CookieName:     "s",
	}))
	m.Use(cache.Cacher(cache.Options{
		Adapter:       "nodb",
		AdapterConfig: "data/cache.db",
	}))
	m.Use(i18n.I18n(i18n.Options{
		Langs: []string{ /*"en-US",*/ "ru-RU"},
		Names: []string{ /*"English",*/ "Русский"},
	}))

	m.Use(macaron.Renderers(macaron.RenderOptions{
		Directory: "templates/default",
		Layout:    "layout",
		Funcs: []template.FuncMap{{
			"raw": func(in string) template.HTML {
				return template.HTML(in)
			},
			"imgUrl": func(path string, opts ...string) string {
				/*if setting.RunMode == "dev" {
					return "http://biblioteka.smolgu.ru" + path
				}*/
				if len(opts) > 0 {
					if opts[0] == "true" {
						path = strings.Replace(path, ".jpg", "_small.jpg", -1)
					}
				}
				return path
			},
			"markdown":    base.RenderMarkdownString,
			"build_hash":  func() string { return setting.BuildHash },
			"str_replace": func(s, old, new string) string { return strings.Replace(s, old, new, -1) },
			"str_split":   strings.Split,
			"str_join":    strings.Join,
			"is_last": func(i int, arr interface{}) bool {
				v := reflect.ValueOf(arr)
				return i+1 == v.Len()
			},
		}}}, "mobile:templates/mobile"))

	m.Use(middleware.Contexter())

	return m
}

func runWeb(ctx *cli.Context) {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	setting.NewContext(ctx.String("mode"), ctx.String("config"))
	m := newMacaron(ctx)

	routers.InitGlobal()

	m.Get("/", routers.Index)

	// used for proxy . todo remove this
	m.Get("/index.shtml", routers.Index)

	m.Get("/pages/*", routers.GetPageByOldId)
	m.Any("/pages/edit/*", auth.MustAdmin, routers.EditPageByOldId)
	m.Get("/resize/*", routers.Resize)

	m.Post("/news/:id/delete", auth.MustAdmin, page.Delete)
	m.Get("/news/:id/restore", auth.MustAdmin, page.Restore)
	m.Get("/news/:id", page.Get)
	m.Any("/news/:id/edit", auth.MustAdmin, page.Edit)

	m.Any("/new", auth.MustAdmin, page.New)

	m.Group("/account", func() {
		//m.Any("/reg", auth.Reg)
		m.Any("/auth", auth.Auth)
		m.Any("/logout", auth.MustSigned, auth.Logout)
		m.Any("/dashboard", auth.MustSigned, account.Dashboard)
	})

	m.Group("/admin/banners", func() {
		m.Get("/", banner.List)
		m.Any("/add", auth.MustAdmin, banner.Add)
		m.Get("/upload", auth.MustAdmin, banner.Upload)
		m.Get("/:id", banner.Show)
		m.Get("/:id/:img", banner.Img)
		m.Get("/:id/hide", banner.Hide)

	})

	m.Get("/tags/:tagName", tag.Show)

	m.Group("/users", func() {
		m.Get("/", users.Users)
		m.Any("/:id/edit", users.Edit)
	}, auth.MustAdmin)

	m.Group("/books", func() {
		m.Get("/:id", book.Show)
		m.Get("/add", book.Add)
		m.Post("/upload", book.Upload)

		m.Get("/page", book.Page)
	})

	m.Group("/links", func() {
		m.Get("/", link.List)
		m.Any("/edit/:id", binding.Bind(models.Link{}), link.Edit)
		m.Get("/delete/:id", link.Delete)
		m.Get("/tags", link.Tags)
		m.Any("/batch", link.Batch)
	})

	/*	m.Group("/katalog", func() {
		m.Get("/search", ruslan.Search)
		m.Get("/book", ruslan.Book)
	})*/

	m.Group("/td", func() {
		m.Get("/", vkr.Index)
		m.Any("/add", vkr.AddDirection)
		m.Get("/:fac", vkr.FacultyDirections)
		m.Get("/:fac/:id/delete", vkr.DeleteDirection)
	}, auth.MustAdmin)

	m.Group("/menus", func() {
		m.Get("/", menus.List)
		m.Post("/create", menus.Create)
		m.Any("/edit", menus.Edit)
		m.Any("/itemcreate", menus.ItemCreate)
		m.Any("/itemedit", menus.ItemEdit)
		m.Post("/setposition", menus.SetPosition)
	}, auth.MustAdmin)

	m.Group("/vkr", func() {
		m.Get("/", vkr.Index)
		m.Any("/add", auth.MustAdmin, vkr.AddVkr)
		m.Get("/:fac", vkr.FacultyDirections)
		m.Get("/:fac/:tdId", vkr.DirectionYears)
		m.Get("/:fac/:tdId/:year", vkr.DirectionYearList)

		m.Get("/download/:id", vkr.Download)
	})

	m.Group("/albums", func() {
		m.Get("/", albums.List)
		m.Get("/:id", albums.Get)
		m.Post("/:id/upload", albums.Upload)
		m.Get("/blob/:id", albums.GetBlob)
		m.Any("/create", albums.Create)
	})

	m.Post("/upload", routers.Upload)

	m.Get("/cmd/index", routers.IndexPages)

	m.Get("/search", routers.Search)
	m.Get("/search/history", routers.SearchHistory)

	m.Get("/bel/searchperson", routers.SearchPerson)
	m.Get("/switch-to-mobile", routers.SwitchToMobile)
	m.Get("/switch-to-full", routers.SwitchToFull)

	m.Any("/files/upload/new", routers.UploadFile)
	m.Get("/files/upload/blob/:id", routers.Blob)

	m.Get("/vktrack", routers.VkTrackerLog)

	m.Get("/*", routers.Route)

	m.NotFound(func() string {
		return "ничего не найдено"
	})

	m.Run(ctx.Int("port"))
}

func expires() string {
	return time.Now().Add(24 * 60 * time.Hour).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
