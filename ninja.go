// Copyright to TechNinja Team
//
//

package main

import (
	"fmt"
	"github.com/icza/gowut/gwu"
)

type myButtonHandler struct {
	counter int
	text    string
}

func (h *myButtonHandler) HandleEvent(e gwu.Event) {
	if b, isButton := e.Src().(gwu.Button); isButton {
		b.SetText(b.Text() + h.text)
		h.counter++
		b.SetToolTip(fmt.Sprintf("You've clicked %d times!", h.counter))
		e.MarkDirty(b)
	}
}

func main() {
	//  Master window
	masterWin := gwu.NewWindow("web-ui-dashboard", "TECH-NINJA!")
	masterWin.Style().SetFullSize()
	masterWin.SetAlign(gwu.HACenter, gwu.VAMiddle)

	p4 := gwu.NewPanel()
	p4.SetHAlign(gwu.HACenter)
	p4.SetCellPadding(2)

	l1 := gwu.NewLabel("Welcome to TechNinja Dashboard")
	l1.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("300%")
	l1.Style().SetColor("green")
	l1.Style().SetBackground("while")
	p4.Add(l1)
	masterWin.Add(p4)


	// Create and build a window
	win := gwu.NewWindow("create-ui", "TECH-NINJA!")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Adding  project label
	headerLabel := gwu.NewLabel("TECH-NINJA DASHBOARD!")
	headerLabel.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("150%")
	headerLabel.Style().SetColor("white")
	headerLabel.Style().SetBackground("green")
	win.Add(headerLabel)

	btn := gwu.NewButton("Create Installation Environment")
	btn.AddEHandler(&myButtonHandler{text: ":-)"}, gwu.ETypeClick)
	win.Add(btn)
	btnsPanel := gwu.NewNaturalPanel()
	btn.AddEHandlerFunc(func(e gwu.Event) {
		// Create and add a new button...
		newbtn := gwu.NewButton(fmt.Sprintf("Created Environment #%d", btnsPanel.CompsCount()))
		newbtn.AddEHandlerFunc(func(e gwu.Event) {
			btnsPanel.Remove(newbtn) // ...which removes itself when clicked
			e.MarkDirty(btnsPanel)
		}, gwu.ETypeClick)
		btnsPanel.Insert(newbtn, 0)
		e.MarkDirty(btnsPanel)
	}, gwu.ETypeClick)
	win.Add(btnsPanel)


	 // Creating list for registered devices
	devices := []string{}
	devices = append(devices,"item1")

	// ADD ZTP clients here
	// TextBox with echo
	q := gwu.NewHorizontalPanel()
	q.Add(gwu.NewLabel("MACHINE UUID:"))
	tb := gwu.NewTextBox("")
	tb.AddSyncOnETypes(gwu.ETypeKeyUp)
	q.Add(tb)

	btn2 := gwu.NewButton("Register Device")
	btn.AddEHandler(&myButtonHandler{text: "Added mac"}, gwu.ETypeClick)
	btn2.AddEHandlerFunc(func(e gwu.Event) {
		devices = append(devices,tb.Text())
		fmt.Printf("Registered MAC in Installation ENV is :%s",tb.Text())
		}, gwu.ETypeClick)

	q.Add(btn2)
	win.Add(q)

	p := gwu.NewHorizontalPanel()
	p.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("Installation Environment Type:"))
    listBox := gwu.NewListBox([]string{"MASTER", "SATELLITE"})
	listBox.SetMulti(false)
	listBox.SetRows(2)
	p.Add(listBox)
	countLabel := gwu.NewLabel("Default Selection: MASTER")
	listBox.AddEHandlerFunc(func(e gwu.Event) {
		countLabel.SetText(fmt.Sprintf("ENVIRONMENT TYPE: %s", listBox.SelectedValue()))
		e.MarkDirty(countLabel)
	}, gwu.ETypeChange)
	p.Add(countLabel)
	win.Add(p)

	// Self-color changer check box
	greencb := gwu.NewCheckBox("Enable Auto Updates!")
	greencb.AddEHandlerFunc(func(e gwu.Event) {
		if greencb.State() {
			greencb.Style().SetBackground(gwu.ClrGreen)
		} else {
			greencb.Style().SetBackground("")
		}
		e.MarkDirty(greencb)
	}, gwu.ETypeClick)
	win.Add(greencb)

	// TextBox with echo
	p = gwu.NewHorizontalPanel()
	p.Add(gwu.NewLabel("Installation Environment Name:"))
	tb1 := gwu.NewTextBox("")
	tb1.AddSyncOnETypes(gwu.ETypeKeyUp)
	p.Add(tb1)
	win.Add(p)

	// Add login
	/*
	p1 := gwu.NewPanel()
	p1.SetHAlign(gwu.HACenter)
	p1.SetCellPadding(2)

	l := gwu.NewLabel("Login")
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("130%")
	p1.Add(l)
	p1.CellFmt(l).Style().SetBorder2(1, gwu.BrdStyleDashed, gwu.ClrNavy)
	l = gwu.NewLabel("user/pass: admin/a")
	l.Style().SetFontSize("80%").SetFontStyle(gwu.FontStyleItalic)
	p1.Add(l)

	errL := gwu.NewLabel("")
	errL.Style().SetColor(gwu.ClrRed)
	p1.Add(errL)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(2, 2)
	table.Add(gwu.NewLabel("User name:"), 0, 0)
	tb3 := gwu.NewTextBox("")
	tb3.Style().SetWidthPx(160)
	table.Add(tb, 0, 1)
	table.Add(gwu.NewLabel("Password:"), 1, 0)
	pb := gwu.NewPasswBox("")
	pb.Style().SetWidthPx(160)
	table.Add(pb, 1, 1)
	p.Add(table)
	b := gwu.NewButton("OK")
	b.AddEHandlerFunc(func(e gwu.Event) {
		if tb.Text() == "admin" && pb.Text() == "a" {
			e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
			e.ReloadWin("main")
		} else {
			e.SetFocusedComp(tb)
			errL.SetText("Invalid user name or password!")
			e.MarkDirty(errL)
		}
	}, gwu.ETypeClick)
	p.Add(b)
	l = gwu.NewLabel("")
	p.Add(l)
	p.CellFmt(l).Style().SetHeightPx(200)

	win.Add(p)
	win.SetFocusedCompID(tb.ID())
	*/
/*
	p = gwu.NewPanel()
	p.SetLayout(gwu.LayoutHorizontal)
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("Here's an ON/OFF switch which enables/disables the other one:"))
	sw := gwu.NewSwitchButton()
	sw.SetOnOff("ENB", "DISB")
	sw.SetState(true)
	p.Add(sw)
	p.Add(gwu.NewLabel("And the other one:"))
	sw2 := gwu.NewSwitchButton()
	sw2.SetEnabled(true)
	sw2.Style().SetWidthPx(100)
	p.Add(sw2)
	sw.AddEHandlerFunc(func(e gwu.Event) {
		sw2.SetEnabled(sw.State())
		e.MarkDirty(sw2)
	}, gwu.ETypeClick)
	win.Add(p)

*/
	// End
	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("techninja.com", "localhost:8081")
	server.SetText("Starting Tech Ninja!!")
	server.AddWin(win)
	server.AddWin(masterWin)
	//server.Start()
	server.Start("web-ui-dashboard")
}
