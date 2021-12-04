package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lukesampson/figlet/figletlib"
)

var version string = "v1.0"

type MessageCarrier struct {
	FontName string
	Messsage string
}

type RESTMessageCarrier struct {
	Fontname string `json:"fontname"`
	Message  string `json:"message"`
}

type OutFontList struct {
	Fontname []string `json:"fontname"`
}

type OutMessage struct {
	Message []string `json:"message"`
}

var SrcFolder string = "fonts"

//go:embed fonts/*.flf
var embededFiles embed.FS

func generateOutputMessage(msgIn *MessageCarrier) string {

	srcFontLocnRef := SrcFolder + "/" + msgIn.FontName
	if !strings.HasSuffix(srcFontLocnRef, ".flf") {
		srcFontLocnRef = srcFontLocnRef + ".flf"
	}

	fontBytes, err := embededFiles.ReadFile(srcFontLocnRef)
	if err != nil {
		fmt.Printf("Cannot find font : %v\n", msgIn.FontName[:len(msgIn.FontName)-4])
	}

	fontData, err := figletlib.ReadFontFromBytes(fontBytes)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return figletlib.SprintMsg(msgIn.Messsage, fontData, 180, fontData.Settings(), "left")
}

func restListAvailableFonts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var outfontlist []string
	log.Println("Received FONTLIST webcall...")
	// get list of fonts and put them into a string slice/array
	fontFileList, _ := embededFiles.ReadDir(SrcFolder)
	for _, fontFileName := range fontFileList {
		// cut off the ".flf"
		outfontlist = append(outfontlist, fontFileName.Name()[:len(fontFileName.Name())-4])
	}
	// dump the string array into the JSON holding struct
	fontEntries := &OutFontList{outfontlist}
	//push the JSON marshalled struct out the door
	json.NewEncoder(w).Encode(fontEntries)

}

func restGenerateTextOutput(w http.ResponseWriter, r *http.Request) {

	log.Println("Received GENERATE webcall...")
	var restmsg RESTMessageCarrier
	reqbody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqbody, &restmsg)
	log.Printf("Fontname: [%v] - Message: [%v]\n", restmsg.Fontname, restmsg.Message)

	msgc := &MessageCarrier{FontName: restmsg.Fontname, Messsage: restmsg.Message}
	outMessage := generateOutputMessage(msgc)

	// marshall the message into JSON and send back
	// var strs string
	// var stra []string
	// for _, y := range outMessage {
	// 	if y == 10 {
	// 		stra = append(stra, strs)
	// 		strs = ""
	// 	} else {
	// 		strs = strs + string(y)
	// 	}
	// }
	// outM := &OutMessage{stra}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(outM)

	// simple TEXT
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, outMessage)

}

func main() {

	// use an int then ensures it's a number, string text could end up with alphas
	flgPort := flag.Int("port", 8888, "Rest server port")
	flag.Parse()

	msgc := &MessageCarrier{FontName: "ghost", Messsage: "Server"}
	fmt.Println(generateOutputMessage(msgc))
	fmt.Printf("Server version %v\n", version)
	fmt.Printf("Running on port: %v\n\n", strconv.Itoa(*flgPort))

	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/v1/genmsg", restGenerateTextOutput).Methods("POST")
	myrouter.HandleFunc("/v1/getfontlist", restListAvailableFonts).Methods("GET")

	log.Println("Started FiggyServer...")
	log.Printf("http://localhost:%v/v1/genmsg      [ POST JSON ( fontname / message ) ]", strconv.Itoa(*flgPort))
	log.Printf("http://localhost:%v/v1/getfontlist [ GET ]", strconv.Itoa(*flgPort))
	log.Println("Ready...")
	http.ListenAndServe(":"+strconv.Itoa(*flgPort), myrouter)
}
