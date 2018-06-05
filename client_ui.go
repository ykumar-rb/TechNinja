// Copyright to TechNinja Team
//
//

package main

import (
	"fmt"
	"strings"
	//"time"
	//"os/exec"
	"net"
	"github.com/go-redis/redis"
	"github.com/icza/gowut/gwu"
	"github.com/ZTP/pnp/executor"
)

const (
	Install  = 1
	Rollback = 2
	Update   = 3
	NoAction = 4
)

const (
	Success = 1
	Failure = 2
)

type SoftwareDB struct {
	Name         string
	Version      string
	AvailVersion string
	Action       string
	Status       string
	Install      string
	UnInstall    string
	Rollback     string

}


// Get preferred outbound ip of this machine
func GetSystemIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Errorf(err.Error())
		return err.Error()
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func ExecuteInstruction(cmd string)(err error) {

	cmdList := strings.Split(cmd,",")

	return executor.ExecuteServerInstructions(cmdList)
	/*
	tokens := strings.Split(cmdList,",")

	for token := range tokens{
		cmd := exec.Command("sh", "-c", token)
		retCode, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
	fmt.Printf("%s\n", retCode) */
}



func SetDataListInDB(client *redis.Client, SDBList []*SoftwareDB) {

	for _, SDB := range SDBList {

		SetDataInDB(client, SDB)
	}
}

func SetDataInDB(client *redis.Client, SDB *SoftwareDB) {
	client.HSet(SDB.Name, "Name", SDB.Name)
	client.HSet(SDB.Name, "Version", SDB.Version)
	client.HSet(SDB.Name, "AvailVersion", SDB.AvailVersion)
	client.HSet(SDB.Name, "Action", SDB.Action)
	client.HSet(SDB.Name, "Status", SDB.Status)
	client.HSet(SDB.Name, "Install", SDB.Install)
	client.HSet(SDB.Name, "UnInstall", SDB.UnInstall)
	client.HSet(SDB.Name, "Rollback", SDB.Rollback)
}

func PrepareKubernetesSetupDummyData() []*SoftwareDB {
	SDBList := []*SoftwareDB{}

	SDB1 := &SoftwareDB{}
	SDB1.Status = "Success"
	SDB1.Action = "Update"
	SDB1.Name = "kubernetes"
	SDB1.Version = "1.9.3-00"
	SDB1.AvailVersion = "1.10.3"
	SDB1.Install= "dummyInstall"
	SDB1.UnInstall= "dummyUnInstall"
	SDB1.Rollback= "dummyRollback"

	SDB2 := &SoftwareDB{}
	SDB2.Status = "Success"
	SDB2.Action = "NoAction"
	SDB2.Name = "docker-ce"
	SDB2.Version = "17.03.2~ce-0~ubuntu-xenial"
	SDB2.AvailVersion = "17.03.2~ce-0~ubuntu-xenial"
	SDB2.Install= "dummyInstall"
	SDB2.UnInstall="dummyUnInstall"
	SDB2.Rollback= "dummyRollback"

	SDBList = append(SDBList, SDB1)
	SDBList = append(SDBList, SDB2)
	return SDBList
}

func PrepareKubernetesKeyList() (keyList []string) {

	keyList = []string{"kubernetes", "docker-ce"}
	return
}

func DatabaseOperation(DBClient *redis.Client, ClientUI gwu.Window) {
	//SDPDataLIst := PrepareKubernetesSetupDummyData()
	//SetDataListInDB(DBClient, SDPDataLIst)

	// Display software details at Ninja Client UI
	keyList := PrepareKubernetesKeyList()
	DisplayAtNinjaClientUI(DBClient, ClientUI, keyList)
}

func GetDataFromDataBase(key string, client *redis.Client) (out SoftwareDB) {
	fmt.Printf("PNP Client DATABASE ")
	outSDB := SoftwareDB{}
	outSDB.Name = client.HGet(key, "Name").Val()
	outSDB.Version = client.HGet(key, "Version").Val()
	outSDB.AvailVersion = client.HGet(key, "AvailVersion").Val()
	outSDB.Action = client.HGet(key, "Action").Val()
	outSDB.Status = client.HGet(key, "Status").Val()
	outSDB.Install = client.HGet(key, "Install").Val()
	outSDB.UnInstall = client.HGet(key, "UnInstall").Val()
	outSDB.Rollback = client.HGet(key, "Rollback").Val()

	fmt.Printf("*** Software Details # Name: %s, Version: %s AvailVersion: %s Action:%v, Status: %v",
		outSDB.Name, outSDB.Version, outSDB.AvailVersion, outSDB.Action, outSDB.Status)

	data := client.HGetAll(key)
	fmt.Printf("All Value: %+v", data)
	return outSDB
}

/*
func startPolling() {
	for {
		<-time.After(5 * time.Second)
		fmt.Println("Polling...")
	}
}
*/

func DisplayAtNinjaClientUI(DBClient *redis.Client, win gwu.Window, keyList []string) {
	// Fetching data from Database for all keys
	var sdb [2]SoftwareDB
	sdb[0] = GetDataFromDataBase(keyList[0], DBClient)
	sdb[1] = GetDataFromDataBase(keyList[1], DBClient)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(20)

	t := gwu.NewTable()
	t.Style().SetBorder2(10, gwu.BrdStyleSolid, gwu.ClrNavy)
	t.SetAlign(gwu.HACenter, gwu.VATop)
	t.Style().SetSize("1000", "400")
	t.EnsureSize(3, 3)
	t.RowFmt(0).Style().SetBackground(gwu.ClrNavy)
	t.RowFmt(0).Style().SetHeight("70")
	t.RowFmt(1).Style().SetBackground("#E6E6FA")
	t.RowFmt(1).Style().SetHeight("70")
	t.RowFmt(2).Style().SetBackground("#E6E6FA")
	t.RowFmt(2).Style().SetHeight("70")

	t.RowFmt(0).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(1).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(2).SetAlign(gwu.HADefault, gwu.VAMiddle)

	img := gwu.NewImage(fmt.Sprintf("Installed Software"), "http://www2.multilizer.com/wp-content/uploads/2014/07/tool.jpg")
	img.Style().SetSize("70", "50")
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
	for row := 1; row < 3; row++ {
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].Name)), row, 0)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].Version)), row, 1)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].AvailVersion)), row, 2)

		var statusStr string
		statusStr = sdb[row-1].Status

		t.Add(gwu.NewLabel(fmt.Sprintf("%s", statusStr)), row, 3)

		var actionStr string
		actionStr = sdb[row-1].Action
		butn1 := gwu.NewButton(fmt.Sprintf("%s", actionStr))
		butn1.Style().SetColor("white")
		butn1.Style().SetBackground(gwu.ClrGreen)
		name := "button" + sdb[row-1].Name
		butn1.SetAttr("ID", name)

		butn1.AddEHandlerFunc(func(e gwu.Event) {
			if butn1.Text() == "UPGRADE" {
				fmt.Printf("UPDATE button pressed!")
				val := butn1.Attr("ID")
				if strings.Contains(val, "kubernetes") {
					fmt.Printf("Hey last UPDATE action was for Kubernetes Software")
					data := GetDataFromDataBase(keyList[0], DBClient)
					err := ExecuteInstruction(data.UnInstall)
					if err != nil {
						fmt.Printf("Error in Uninstalling kubernetes")
					} else {
						err = ExecuteInstruction(data.Install)
						if err != nil {
							DBClient.HSet(keyList[0], "Status", "Operation failed")
							butn1.SetText("ROLLBACK")
						} else {
							DBClient.HSet(keyList[0], "Version", data.AvailVersion)
							butn1.SetText("NOACTION")
						}
					}
				} else if strings.Contains(val, "docker-ce") {
					fmt.Printf("Hey last UPDATE action was for Docker Software")
					data := GetDataFromDataBase(keyList[1], DBClient)
					err := ExecuteInstruction(data.UnInstall)
					if err != nil {
						fmt.Printf("Error in Uninstalling docker")
					} else {
						err = ExecuteInstruction(data.Install)
						if err != nil {
							DBClient.HSet(keyList[1], "Status", "Operation failed")
							butn1.SetText("ROLLBACK")
						} else {
							newVersion := DBClient.HGet(keyList[1], "AvailVersion")
							DBClient.HSet(keyList[1], "Version", newVersion)
							butn1.SetText("NOACTION")
						}
					}
				}
			} else if butn1.Text() == "INSTALL" {
				fmt.Printf("INSTALL button pressed!")
				val := butn1.Attr("ID")
				if strings.Contains(val, "kubernetes") {
					fmt.Printf("Hey last INSTALL action was for Kubernetes Software")
					data := GetDataFromDataBase(keyList[0], DBClient)
					err := ExecuteInstruction(data.Install)
					if err != nil {
							butn1.SetText("NOACTION")
						    DBClient.HSet(keyList[0],"Status","Operation failed")
						}
					DBClient.HSet(keyList[0],"Status","Operation Success")

				} else if strings.Contains(val, "docker-ce") {
					fmt.Printf("Hey last INSTALL action was for Docker Software")
					data := GetDataFromDataBase(keyList[1], DBClient)
					err := ExecuteInstruction(data.Install)
					if err != nil {
						butn1.SetText("NOACTION")
						DBClient.HSet(keyList[1], "Status", "Operation failed")
					}
					DBClient.HSet(keyList[1], "Status", "Operation Success")
				}
			} else if butn1.Text() == "ROLLBACK" {
				fmt.Printf("ROLLBACK button pressed!")
				val := butn1.Attr("ID")
				if strings.Contains(val, "kubernetes") {
					fmt.Printf("Hey last ROLLBACK action was for Kubernetes Software")
					data := GetDataFromDataBase(keyList[0],DBClient)
					err := ExecuteInstruction(data.Rollback)
					if err != nil{
						DBClient.HSet(keyList[0],"Status","Operation failed")
					}else{
						butn1.SetText("NOACTION")
						DBClient.HSet(keyList[0],"Status","Operation Success")
					}
				} else if strings.Contains(val, "docker-ce") {
					fmt.Printf("Hey last ROLLBACK action was for Docker Software")
					data := GetDataFromDataBase(keyList[1],DBClient)
					err := ExecuteInstruction(data.Rollback)
					if err != nil{
						DBClient.HSet(keyList[1],"Status","Operation failed")
					}else{
						butn1.SetText("NOACTION")
						DBClient.HSet(keyList[1],"Status","Operation Success")
					}
				}
			} else if butn1.Text() == "NOACTION" {
				fmt.Printf("NOACTION button pressed!")
			} else {
				fmt.Printf("UNKNOWN button pressed!")
			}
		}, gwu.ETypeClick)

		t.Add(butn1, row, 4)
	}

	p.Add(t)
	p.Add(btnsPanel)
	win.Add(p)

}

func setNoWrap(panel gwu.Panel) {
	count := panel.CompsCount()
	for i := count - 1; i >= 0; i-- {
		panel.CompAt(i).Style().SetWhiteSpace(gwu.WhiteSpaceNowrap)
	}
}

func main() {
	//  Master window
	masterWin := gwu.NewWindow("web-ui-dashboard", "TECH-NINJA CLIENT GUI !")
	masterWin.Style().SetFullSize()
	masterWin.SetAlign(gwu.HACenter, gwu.VAMiddle)

	/* Master window */
	p4 := gwu.NewPanel()
	p4.SetHAlign(gwu.HACenter)
	p4.SetCellPadding(2)
	l1 := gwu.NewLabel("Welcome to TechNinja Dashboard")
	l1.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("300%")
	l1.Style().SetColor("green")
	l1.Style().SetBackground("while")
	p4.Add(l1)
	masterWin.Add(p4)

	/* Display window for software Catalog */
	ClientWin := gwu.NewWindow("display-ui", "TECH-NINJA CLIENT GUI!")
	ClientWin.Style().SetFullWidth()
	ClientWin.SetHAlign(gwu.HACenter)
	ClientWin.SetCellPadding(2)

	ClientWin.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			fmt.Println("LOADING window:", e.Src().ID())
		case gwu.ETypeWinUnload:
			fmt.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)

	p := gwu.NewPanel()
	p.Style().SetFullWidth().SetBorderBottom2(7, gwu.BrdStyleSolid, "#cccccc")
	p.Style().SetBackground(gwu.ClrGreen)
	//p.Style().SetFullWidth().SetBorderBottom2(10, gwu.BrdStyleSolid, "#cccccc")
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)

	img1 := gwu.NewImage("riverbed", "https://www.riverbed.com/content/dam/riverbed-www/en_US/Images/" +
		"fpo/logo_riverbed_orange.png?redesign=true")

	p.Add(img1)

	l2 := gwu.NewLabel("PNP Software Catalog")
	l2.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("250%")
	l2.Style().SetColor("white")
	p.Add(l2)

	setNoWrap(p)

	p.AddHSpace(10)
	Refresh := gwu.NewLink("Refresh", "#")
	Refresh.Style().SetColor(gwu.ClrNavy)
	Refresh.SetTarget("")
	Refresh.Style().SetMarginLeftPx(5)
	Refresh.AddEHandlerFunc(func(e gwu.Event) {
		e.RemoveSess()
		e.ReloadWin("display-ui")
	}, gwu.ETypeClick)

	p.Add(Refresh)
	ClientWin.Add(p)

	redisServerEndPoint := "localhost:6389"
	// Database object creation
	DBClient := redis.NewClient(&redis.Options{
		Addr:     redisServerEndPoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := DBClient.Ping().Result()
	fmt.Println(pong, err)
	DatabaseOperation(DBClient, ClientWin)

	clientUIAddr := GetSystemIP()+":8081"
	// Adding all windows to server
	server := gwu.NewServer("techninja.com", clientUIAddr)
	server.SetText("Starting Tech Ninja!!")
	server.AddWin(ClientWin)
	server.AddWin(masterWin)

	//go startPolling()
	server.Start("display-ui")
}
