package main
import (
	"log"
	"net/http"
	//"github.com/gorilla/websocket"
	//"errors"
	"fmt"
	//"io"
	//"net"
	"os/exec"
	//"os"
	//"runtime"
	"time"
	//"math"
	"strconv"
	//"reflect"
	"github.com/rs/cors"
	"encoding/json"
	//"reflect"
	//"github.com/pion/webrtc/v3"
	//"github.com/pion/webrtc/v3/examples/internal/signal"
)
// define  and initialize variables
var ar_goport =[3]string{"8081","8082","8083"}
var ar_ffmpegport =[9]string{"5004","5014","5024","5034","5044","5054","5064","5074","5084"}
var ar_display  =[3]string {"0","1","2"}
type Tuple struct {
    a, b , c  string
}
var mp_id = make(map[string]Tuple)  //id->goport,ffmpegport,display
var mp_ffmpegport = make(map[string]string)   // fport->id
var mp_goport =make(map[string]string)      //goport->id
var mp_display =make(map[string]string)        //display->id
var mp_id_goprocess= make(map[string]*exec.Cmd)
var mp_id_fprocess =make(map[string]*exec.Cmd)
func ffmpegdelete(fporttodelete string) {
    // do stuff
    DurationOfTime := time.Duration(15) * time.Second
    time.AfterFunc(DurationOfTime, func() {log.Print("killed ffmpeg port ",fporttodelete);delete(mp_ffmpegport,fporttodelete); })
    // do other stuff
}
func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
    if x > y {
        return y
    }
    return x
}
func main() {
	mp := map[string]string{
        "0": "601,0",
        "1": "0,0",
        "2": "0,401",
	}

        mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	// mux.Handle("/static/", fs)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

        mux.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("gfetched in get",r)
		Id, ok := r.URL.Query()["id"]
		id:=Id[0]
		width := r.URL.Query().Get("width")
		height := r.URL.Query().Get("height")
		//h:= r.URL.Query()["height"]
		log.Print("aspect ratio: ",width,"x",height)
		log.Print(id,ok)
		// now we have to get one ffmpeg port and run ffmpeg command and further get one goport and run goserver and send goport as response
		var fport string
		var gport string
		var display string
		// assign one ffmpeg port to one id for the current connection
		for i := 0; i < len(ar_ffmpegport); i++ {
		  value, ok := mp_ffmpegport[ar_ffmpegport[i]]
		  if(!ok){
		     fport=ar_ffmpegport[i];
		     _=value
		     break;
		   }
		}
		// assign one go port to the connection	
		for i := 0; i < len(ar_goport); i++ {
		   value,ok := mp_goport[ar_goport[i]]
		   if(!ok){
		      gport=ar_goport[i];
		      _=value
		      break;
		    }
		}
		// assign one display to the current connection
		for i := 0; i < len(ar_display); i++ {
		   value,ok := mp_display[ar_display[i]]
		   if(!ok){
		      display=ar_display[i];
		      _=value
		      break;
		    }
		}
		// if all ports are assigned return
		if(fport=="" || gport =="" || display==""){
		        log.Print("no empty slot found for ",id,"user")
			return
		}
		if(display != "0"){
			intw,err:=strconv.Atoi(width)
			inth,err:=strconv.Atoi(height)
		        intw=Min(intw,640);
		        inth=Min(inth,480);
		
		        width=strconv.Itoa(intw)
		        height=strconv.Itoa(inth)
			log.Print(err);
	        }
		str3:= width+"x"+height
		aratio:= width+","+height
		//str3:="1200x750";
		//aratio:="1200,750"
		// log.Print("port is ",fport,gport)
		Var, err := strconv.Atoi(display)
		Var=Var*600+500
		str:= strconv.Itoa(Var)
		log.Println("starting position", str)
		log.Println("ffmpeg running on port: ", fport)
		log.Println("map: ",mp["0"])
		// cmd :=exec.Command("ffmpeg -f x11grab -draw_mouse 0  -s 1200x821  -framerate 30 -i :0.0+0,0   -pix_fmt yuv420p -vcodec h264_nvenc -g 10 -threads 2 -preset llhq -fflags   nobuffer -f rtp -sdp_file nvidia.sdp rtp://127.0.0.1:5004?pkt_size=1200")
		// cmd, Oerr := exec.Command("ffmpeg","-f","x11grab","-draw_mouse","0", "-s","600x400","-framerate" ,"30","-i",":0.0+0,0","-pix_fmt", "yuv420p", "-vcodec", "h264_nvenc", "-g", "10","-threads","2","-preset","llhq","-fflags","nobuffer","-f","rtp","-sdp_file","nvidia.sdp","rtp://127.0.0.1:"+fport+"?pkt_size=1200" ).CombinedOutput()
			cmd:= exec.Command("ffmpeg","-f","x11grab","-r","30","-draw_mouse","0", "-s",str3,"-framerate" ,"30","-i",":0."+display+"+0,0","-pix_fmt", "yuv444p", "-vcodec", "h264_nvenc", "-g", "10","-threads","2","-preset","llhq","-fflags","nobuffer","-f","rtp","-sdp_file","nvidia.sdp","rtp://127.0.0.1:"+fport+"?pkt_size=1200" )
		// log.Print("type of cmd is ",reflect.TypeOf(cmd))
//	cmd:= exec.Command("ffmpeg","-f","x11grab","-draw_mouse","0", "-s",str3,"-framerate" ,"60","-i",":0."+display+"+0,0","-pix_fmt", "yuv420p", "-vcodec", "h264_nvenc","-g", "10","-threads","2","-preset","llhq","-fflags","nobuffer","-f","rtp","-sdp_file","nvidia.sdp","rtp://127.0.0.1:"+fport+"?pkt_size=1200" )
//	err = cmd.Start()
		//err:=error(nil)
		//log.Print("output is ",string(cmd3));
		if err != nil {
		   fmt.Printf("error in ffmpeg command%s", err)
		}

                //here have to run ffmpeg command in some particular directory
                //cmd2:=exec.Command("go", "run", "./dockerffmpeg/main2.go",str3,  fport, display,"&")
                //err2 := cmd2.Start()
                //if err2 != nil {
                  // fmt.Printf("error in docker  command%s", err2)
                //}
                //srcFolder := "./test.yaml"
                //destFolder := "/wdir/test.yaml"
                //cpCmd := exec.Command("cp", srcFolder, destFolder)
                //err = cpCmd.Run()


		log.Print("entering into second command",str3)
		//now run goserver process
		//time.Sleep(5000 * time.Millisecond)
		cmd2:=exec.Command("go", "run", "server3.go", gport, fport, display,id, aratio,"&")
		err2:= cmd2.Start()
		if err2 != nil {
		   fmt.Printf("error in goserver command%s", err2)
		}
		//have to uncomment below ///////////////////////////////////seeeeeee herererererererer
		mp_id_fprocess[id]=cmd
		mp_id_goprocess[id]=cmd2
		mp_id[id]=Tuple{gport,fport,display}
		mp_ffmpegport[fport]=id
		mp_goport[gport]=id
		mp_display[display]=id
		time.Sleep(3000 * time.Millisecond)
		// log.Print("All done")
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*"
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//(w).Header().Set("Access-Control-Allow-Origin", "*")
		//(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//setupCORS(&w, r)
		resp := make(map[string]string)
		resp["port"] = gport
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		//w.Write([]byte({"hello" :"world"}))
		w.Write(jsonResp)
		return
	})
	mux.HandleFunc("/close/", func(w http.ResponseWriter, r *http.Request) {
                Id, ok := r.URL.Query()["id"]
                id:=Id[0]
		log.Print("closing this id",id,ok)
                //both server and ffmpeg close
		//then release all maps
//		if err := mp_id_goprocess[id].Process.Kill(); err != nil {
  //      log.Fatal("failed to kill go process: ", err)
   // }
     // log.Println("killing go server  o port",mp_id[id].a)

        if err := mp_id_fprocess[id].Process.Kill(); err != nil {
        log.Fatal("failed to kill ffmpeg process: ", err)
    }
         log.Println("killing ffmpeg o port",mp_id[id].b)
               delete(mp_goport,mp_id[id].a)
	       //delete(mp_ffmpegport,mp_id[id].b)
	       delete(mp_display,mp_id[id].c)
	      ffmpegdelete(mp_id[id].b);
	       return
	})
	handler := cors.AllowAll().Handler(mux)
	log.Fatal(http.ListenAndServe(":8090", handler))
	//err := http.ListenAndServe(":8090", handler)
	//if(err != nil){
	//	log.Print("Error starting server: ",err)	//}
	//fmt.Print("listening on port 8090")
}
