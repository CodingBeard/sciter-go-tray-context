package context

import (
	"encoding/json"
	"github.com/lxn/win"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
	"strconv"
	"unsafe"
)

type MenuItem struct {
	Text string
	ClickCallback func()
	ShouldShow func() bool
}

type Menu struct {
	Items []MenuItem
	window *window.Window
}

type sciterMenuItem struct {
	Id int
	Text string
}

type sciterXYW struct {
	X int
	Y int
	W int
}

func (m *Menu) DisplayContextMenu(x, y, w int) *window.Window {
	menu, err := window.New(sciter.SW_POPUP,nil)

	if err != nil {
		log.Fatal(err)
	}

	menu.DefineFunction("error", func(args ...*sciter.Value) *sciter.Value {
		log.Println(args[0].String())
		return sciter.NullValue()
	})

	menu.DefineFunction("getMenuItems", func(args ...*sciter.Value) *sciter.Value {
		var items []sciterMenuItem

		for id, item := range m.Items {
			if item.ShouldShow == nil || item.ShouldShow() {
				items = append(items, sciterMenuItem{Id: id, Text: item.Text})
			}
		}

		jsonBytes, e := json.Marshal(items)
		if e != nil {
			log.Println(e)
		}

		return sciter.NewValue(string(jsonBytes))
	})

	menu.DefineFunction("getXYW", func(args ...*sciter.Value) *sciter.Value {
		xyw := sciterXYW{X: x, Y: y, W: w}
		jsonBytes, e := json.Marshal(xyw)
		if e != nil {
			log.Println(e)
		}

		return sciter.NewValue(string(jsonBytes))
	})

	menu.DefineFunction("menuItemClicked", m.menuItemClicked)

	e := menu.LoadHtml(getHtml(), "")
	if e != nil {
		log.Println(e)
	}

	menu.Show()

	hwnd := win.HWND(unsafe.Pointer(menu.GetHwnd()))

	win.ShowWindow(hwnd, win.SW_HIDE)

	flags := win.GetWindowLongPtr(hwnd, win.GWL_EXSTYLE) | win.WS_EX_TOOLWINDOW
	flags -= win.WS_EX_APPWINDOW
	win.SetWindowLongPtr(hwnd, win.GWL_EXSTYLE, flags)

	win.ShowWindow(hwnd, win.SW_SHOW)

	win.SetWindowPos(hwnd, win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE | win.SWP_NOSIZE)

	win.SetFocus(hwnd)

	m.window = menu

	return menu
}

func (m *Menu) menuItemClicked(args ... *sciter.Value) *sciter.Value {
	for id, item := range m.Items {
		itemId, e := strconv.Atoi(args[0].String())

		if e != nil {
			panic(e)
		}

		if id == itemId {
			item.ClickCallback()
			m.window.Eval("view.close(null)")
		}
	}

	return sciter.NewValue(true)
}

func getHtml() string {
	return `<html window-frame="extended">
<head>
    <style>
        html {
            font: system;
            overflow: none;
            background-color: white;
            margin: 0;
            padding: 0;
        }

        body {
            padding: 0;
            margin: 0;
        }

        ul {
            padding: 2dip;
            margin: 0;
        }

        ul li {
            margin: 0;
            list-style: none;
            padding-left: 10dip;
            padding-top: 5dip;
            height: 20dip;
        }

        ul li:hover {
            background: #3499ea;
        }

    </style>
    <script type="text/tiscript">
		VM.unhandledExceptionHandler = function (err) {
			view.error(err);
		};

        var menuItems = [];

        function positionMenu(x, y, w)
        {
            var h = (menuItems.length * 25) + 4;
            w = self.toPixels(w + "dip", #width);
            h = self.toPixels(h + "dip", #height);
            view.move(x, y - h, w, h, true);
        }

        function registerMenuItem(id, text)
        {
            menuItems.push([id, text]);
            $(ul).append("<li id='menuitem" + id + "' data-id='" + id + "'>" + text + "</li>");
            self.select("#menuitem" + id).on("click", function (e) {
                view.menuItemClicked(this.attributes["data-id"]);
				view.close(null);
            });
        }

        view.root.onFocus = function (evt) {
          if( evt.type == Event.FOCUS_OUT ) { view.close(null); }
        }

		var items = parseData(view.getMenuItems());
		for (var item in items) {
			registerMenuItem(item.Id, item.Text);
		}

		var xyw = parseData(view.getXYW());
		positionMenu(xyw.X, xyw.Y, xyw.W);
    </script>
</head>
<body>
<ul></ul>
</body>
</html>`
}
