package context

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
	"log"
	"github.com/lxn/win"
	"unsafe"
	"strconv"
)

type Menu struct {
	Items []MenuItem
	window *window.Window
}

func (m *Menu) DisplayContextMenu(x, y, w int) *window.Window {
	menu, err := window.New(sciter.SW_POPUP,nil)

	if err != nil {
		log.Fatal(err)
	}

	menu.LoadHtml(getHtml(), "")
	for id, item := range m.Items {
		menu.Call("registerMenuItem", sciter.NewValue(id), sciter.NewValue(item.Text))
	}
	menu.DefineFunction("menuItemClicked", m.menuItemClicked)
	menu.Call(
		"positionMenu",
		sciter.NewValue(x),
		sciter.NewValue(y),
		sciter.NewValue(w),
	)
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

        var menuItems = [];

        function positionMenu(x, y, w)
        {
            var h = ((menuItems.length - 1) * 25);
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
    </script>
</head>
<body>
<ul></ul>
</body>
</html>`
}