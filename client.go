// Copyright to TechNinja Team
//
//

package main

import (
	"fmt"
	"github.com/icza/gowut/gwu"
	"github.com/go-redis/redis"
	"strconv"
)


const (
	Install = 1
	Rollback = 2
	Update  =  2
        NoAction = 3
)

/*
type Status int
const (
	Success   Status = 1
	Failure   Status = 2
)
*/

const (
	Success = 1
	Failure = 2
)

type SoftwareDb struct{
	Name string
	Version string
	AvailVersion string
	Action int
	Status int
}


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

func DatabaseOperation(client *redis.Client, win gwu.Window) {

	/* Building software database */
	SDB := SoftwareDb{}
	SDB.Status = Success
	SDB.Action = Install
	SDB.Name ="Kubernetes"
	SDB.Version ="1.0.0"
	SDB.AvailVersion ="2.0.0"

	fmt.Println("DataBase Client Function")
	err := client.Set(SDB.Name, SDB.Name, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(SDB.Name).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Software Name: ", val)
	fmt.Println("   ########################################################  ")

	client.HSet(SDB.Name,"Name",SDB.Name)
	client.HSet(SDB.Name,"Version",SDB.Version)
	client.HSet(SDB.Name,"AvailVersion", SDB.AvailVersion)
	client.HSet(SDB.Name,"Action", SDB.Action)
	client.HSet(SDB.Name,"Status", SDB.Status)


    name := client.HGet(SDB.Name,"Name")
    version := client.HGet(SDB.Name,"Version")
	availVersion := client.HGet(SDB.Name,"AvailVersion")
	action := client.HGet(SDB.Name,"Action")
	status := client.HGet(SDB.Name,"Status")
	fmt.Printf("*** Software Details # Name: %s, Version: %s AvailVersion: %s Action:%v, Status: %v",
		name.Val(),version.Val(), availVersion.Val(), action.Val(),status.Val())

	data := client.HGetAll(SDB.Name)
	fmt.Printf("All Value: %+v",data)


    // Display software details at Ninja Client UI
	DisplayAtNinjaClientUI(client,win,SDB.Name)
}

func getDataFromDataBase(key string, client *redis.Client)(out *SoftwareDb){
	fmt.Printf("PRINTING DATABASE")
	outSDB := &SoftwareDb{}
	outSDB.Name = client.HGet(key,"Name").Val()
	outSDB.Version = client.HGet(key,"Version").Val()
	outSDB.AvailVersion = client.HGet(key,"AvailVersion").Val()
	outSDB.Action, _ = strconv.Atoi(client.HGet(key,"Action").Val())
	outSDB.Status ,_= strconv.Atoi(client.HGet(key,"Status").Val())


	fmt.Printf("*** Software Details # Name: %s, Version: %s AvailVersion: %s Action:%v, Status: %v",
		outSDB.Name,outSDB.Version, outSDB.AvailVersion, outSDB.Action,outSDB.Status)

	data := client.HGetAll(key)
	fmt.Printf("All Value: %+v",data)
	return outSDB
}

func DisplayAtNinjaClientUI(client *redis.Client, win gwu.Window, key string){

	sdb := getDataFromDataBase(key,client)
	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(20)

	t := gwu.NewTable()
	t.Style().SetBorder2(10, gwu.BrdStyleSolid, gwu.ClrNavy)
	t.SetAlign(gwu.HARight, gwu.VATop)
	t.Style().SetSize("1000","500")
	t.EnsureSize(5, 5)
	t.RowFmt(0).Style().SetBackground(gwu.ClrNavy)

	t.RowFmt(0).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(1).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(2).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(3).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(4).SetAlign(gwu.HADefault, gwu.VAMiddle)


	img := gwu.NewImage(fmt.Sprintf("Installed Software"), "http://www2.multilizer.com/wp-content/uploads/2014/07/tool.jpg")
	img.Style().SetSize("70","50")
	t.Add(img, 0, 0)

	lb1 := gwu.NewLabel(fmt.Sprintf("Current Version"))
	//lb1.Style().SetBackground("blue")
	lb1.Style().SetColor("white")
	lb1.Style().SetWidth("20")

	lb2 := gwu.NewLabel(fmt.Sprintf("Available Version"))
	lb2.Style().SetColor("white")
	lb2.Style().SetWidth("20")

	lb3 := gwu.NewLabel(fmt.Sprintf("Status"))
	lb3.Style().SetColor("white")
	lb3.Style().SetWidth("20")

	lb4 := gwu.NewLabel(fmt.Sprintf("Action"))
	lb4.Style().SetColor("white")
	lb4.Style().SetWidth("20")


	t.Add(lb1, 0, 1)
	t.Add(lb2, 0, 2)
	t.Add(lb3, 0, 3)
	t.Add(lb4, 0, 4)

	btnsPanel := gwu.NewNaturalPanel()

	for row := 1; row < 2; row++ {
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb.Name)), row, 0)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb.Version)), row, 1)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb.AvailVersion)), row, 2)

		var statusStr string
		if sdb.Status == Success{
			statusStr = "Operation Success"
		}else{
			statusStr = "Operation Failure"
		}
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", statusStr)), row, 3)

		var actionStr string
		if sdb.Action == Install {
			actionStr = "UPDATE"

		}else{
			actionStr = "ROLLBACK"
		}

		butn1 := gwu.NewButton(fmt.Sprintf("%s", actionStr))
		butn1.Style().SetColor("white")
		butn1.Style().SetBackground("green")

		butn1.AddEHandlerFunc(func(e gwu.Event) {
			//newbtn := gwu.NewButton(fmt.Sprintf("Created Environment #%d", btnsPanel.CompsCount()))
			//btnsPanel.Insert(newbtn, 0)
			fmt.Printf("Update/Rollback button pressed!")
		}, gwu.ETypeClick)

		t.Add(butn1,row,4)
	}

	for row := 2; row < 5; row++ {
		for col := 0; col < 5; col++ {
			t.Add(gwu.NewLabel(fmt.Sprintf("Button %d%d", row, col)), row, col)
		}
	}

	p.Add(t)
	p.Add(btnsPanel)
	win.Add(p)
}

func main() {
	//  Master window
	masterWin := gwu.NewWindow("web-ui-dashboard", "TECH-NINJA CLIENT GUI !")
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
	win := gwu.NewWindow("display-ui", "TECH-NINJA CLIENT GUI!")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)
	
	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)
        l2 := gwu.NewLabel("PNP Software Catalog")
	l2.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("300%")
	l2.Style().SetColor("green")
	l2.Style().SetBackground("while")
	p.Add(l2)
	win.Add(p)

	// Database object creation
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6389",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	DatabaseOperation(client,win)

/*
	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)
        t := gwu.NewTable()
	t.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrGrey)
	t.SetAlign(gwu.HARight, gwu.VATop)
	t.EnsureSize(5, 5)
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			t.Add(gwu.NewButton(fmt.Sprintf("Button %d%d", row, col)), row, col)
		}
	}
	t.SetColSpan(2, 1, 2)
	t.SetRowSpan(3, 1, 2)
	t.CellFmt(2, 2).Style().SetSizePx(150, 80)
	t.CellFmt(2, 2).SetAlign(gwu.HARight, gwu.VABottom)
	t.RowFmt(2).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.CompAt(2, 1).Style().SetFullSize()
	t.CompAt(4, 2).Style().SetFullWidth()
	t.RowFmt(0).Style().SetBackground(gwu.ClrRed)
	t.RowFmt(1).Style().SetBackground(gwu.ClrGreen)
	t.RowFmt(2).Style().SetBackground(gwu.ClrBlue)
	t.RowFmt(3).Style().SetBackground(gwu.ClrGrey)
	t.RowFmt(4).Style().SetBackground(gwu.ClrTeal)
	p.Add(t)
	win.Add(p)
*/	
	server := gwu.NewServer("techninja.com", "localhost:8081")
	server.SetText("Starting Tech Ninja!!")
	server.AddWin(win)
	server.AddWin(masterWin)

        //server.Start()
	server.Start("display-ui")
}
