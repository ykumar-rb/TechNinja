// Copyright to TechNinja Team
//
//

package main

import (
	"fmt"
	"github.com/icza/gowut/gwu"
	"github.com/go-redis/redis"
)


const (
	Install = 1
	Rollback = 2
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

func DatabaseOperation(client *redis.Client) {

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
        l2 := gwu.NewLabel("Welcome to TechNinja Display Dashboard")
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

	DatabaseOperation(client)
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
	server.Start("web-ui-dashboard")
}
