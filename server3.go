package main


import (
	"time"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "errors"
    "fmt"
    "io"
    "net"
    "os/exec"
    "os"
    //"runtime"
    //"time"
    "math"
    "strconv"
    "reflect"
    "encoding/json"
   "github.com/pion/webrtc/v3"
 //   "github.com/pion/webrtc/v3/examples/internal/signal"
   //"signal.go"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	
	"io/ioutil"
	//"os"
	"strings"
)

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// MustReadStdin blocks until input is received from stdin
func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	fmt.Println("")

	return in
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}
var port,fport,display,id  string

func aman2(sdp string,conn *websocket.Conn){
   fmt.Printf("hello world in v3")
        peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
                ICEServers: []webrtc.ICEServer{
                        {
                                URLs: []string{"stun:stun.l.google.com:19302"},
                        },
                },
        })
        if err != nil {
                panic(err)
        }
        // Open a UDP Listener for RTP Packets on port 5004
	val,err:=strconv.Atoi(fport)
	log.Println(val);
        listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004})
	//listeneraudio, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5010})
	//log.Println(listeneraudio);
        if err != nil {
                panic(err)
        }
        defer func() {
                if err = listener.Close(); err != nil {
                        panic(err);
                }
        }()
        // Create a video track
        videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	//audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "audio", "pion")
        if err != nil {
                panic(err)
        }
        rtpSender, err := peerConnection.AddTrack(videoTrack)
	//rtpSenderaudio , err :=peerConnection.AddTrack(audioTrack)
        if err != nil {
                panic(err)
        }
        // Read incoming RTCP packets
        // Before these packets are returned they are processed by interceptors. For things
        // like NACK this needs to be called.
        go func() {
                rtcpBuf := make([]byte, 1500)
                for {
                        if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
                                return
                        }

			//if _, _, rtcpErraudio := rtpSenderaudio.Read(rtcpBuf); rtcpErraudio != nil {
                          //      return
                        //}
                }
        }()
	// Set the handler for ICE connection state
        // This will notify you when the peer has connected/disconnected
        peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
                fmt.Printf("Connection State has changed %s \n", connectionState.String())
                if connectionState.String() == "closed"{
		 listener.Close();
		}
		if connectionState.String() == "disconected"{
                 listener.Close();
                }
		if connectionState.String() == "failed"{
                 listener.Close();
                }
                if connectionState == webrtc.ICEConnectionStateFailed {
                        if closeErr := peerConnection.Close(); closeErr != nil {
                                panic(closeErr)
				listener.Close();
                        }
                }
        })

        // Wait for the offer to be pasted
        offer := webrtc.SessionDescription{}
        Decode(sdp, &offer)
        log.Print("checking ")
        // Set the remote SessionDescription
        if err = peerConnection.SetRemoteDescription(offer); err != nil {
                panic(err)
        }
        // Create answer
        answer, err := peerConnection.CreateAnswer(nil)
        if err != nil {
                panic(err)
        }
        // Create channel that is blocked until ICE Gathering is complete
        gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
        // Sets the LocalDescription, and starts our UDP listeners
        if err = peerConnection.SetLocalDescription(answer); err != nil {
                panic(err)
        }
        // Block until ICE Gathering is complete, disabling trickle ICE
        // we do this because we only can exchange one signaling message
        // in a production application you should exchange ICE Candidates via OnICECandidate
        <-gatherComplete
        // Output the answer in base64 so we can paste it in browser
        //fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
        //this we have to send 
          err = conn.WriteMessage(1,[]byte(Encode(*peerConnection.LocalDescription())))
        // Read RTP packets forever and send them to the WebRTC Client
        inboundRTPPacket := make([]byte, 1600) // UDP MTU
	 for {
                n, _, err := listener.ReadFrom(inboundRTPPacket)
	//	naudio, _, err := listeneraudio.ReadFrom(inboundRTPPacket)
	//	log.Println(naudio);
                if err != nil {
                        panic(fmt.Sprintf("error during read: %s", err))
                }

                if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
                        if errors.Is(err, io.ErrClosedPipe) {
                                // The peerConnection has been closed.
                                return
                        }

                        panic(err)
                }

		//if _, err = audioTrack.Write(inboundRTPPacket[:n]); err != nil {
                  //      if errors.Is(err, io.ErrClosedPipe) {
                                // The peerConnection has been closed.
                    //            return
                      //  }

                        //panic(err)
                //}
        }
}

func aman(sdp string,conn *websocket.Conn){
   fmt.Printf("hello world in v3")
        peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
                ICEServers: []webrtc.ICEServer{
                        {
                                URLs: []string{"stun:stun.l.google.com:19302"},
                        },
                },
        })
        if err != nil {
                panic(err)
        }
        // Open a UDP Listener for RTP Packets on port 5004
	val,err:=strconv.Atoi(fport)
	log.Println(val);
        listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: val})
	//listeneraudio, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5010})
	//log.Println(listeneraudio);
        if err != nil {
                panic(err)
        }
        defer func() {
                if err = listener.Close(); err != nil {
                        panic(err);
                }
        }()
        // Create a video track
        videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	//audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "audio", "pion")
        if err != nil {
                panic(err)
        }
        rtpSender, err := peerConnection.AddTrack(videoTrack)
	//rtpSenderaudio , err :=peerConnection.AddTrack(audioTrack)
        if err != nil {
                panic(err)
        }
        // Read incoming RTCP packets
        // Before these packets are returned they are processed by interceptors. For things
        // like NACK this needs to be called.
        go func() {
                rtcpBuf := make([]byte, 1500)
                for {
                        if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
                                return
                        }

			//if _, _, rtcpErraudio := rtpSenderaudio.Read(rtcpBuf); rtcpErraudio != nil {
                          //      return
                        //}
                }
        }()
	// Set the handler for ICE connection state
        // This will notify you when the peer has connected/disconnected
        peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
                fmt.Printf("Connection State has changed %s \n", connectionState.String())
                if connectionState.String() == "closed"{
		 listener.Close();
		}
		if connectionState.String() == "disconected"{
                 listener.Close();
                }
		if connectionState.String() == "failed"{
                 listener.Close();
                }
                if connectionState == webrtc.ICEConnectionStateFailed {
                        if closeErr := peerConnection.Close(); closeErr != nil {
                                panic(closeErr)
				listener.Close();
                        }
                }
        })

        // Wait for the offer to be pasted
        offer := webrtc.SessionDescription{}
        Decode(sdp, &offer)
        log.Print("checking ")
        // Set the remote SessionDescription
        if err = peerConnection.SetRemoteDescription(offer); err != nil {
                panic(err)
        }
        // Create answer
        answer, err := peerConnection.CreateAnswer(nil)
        if err != nil {
                panic(err)
        }
        // Create channel that is blocked until ICE Gathering is complete
        gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
        // Sets the LocalDescription, and starts our UDP listeners
        if err = peerConnection.SetLocalDescription(answer); err != nil {
                panic(err)
        }
        // Block until ICE Gathering is complete, disabling trickle ICE
        // we do this because we only can exchange one signaling message
        // in a production application you should exchange ICE Candidates via OnICECandidate
        <-gatherComplete
        // Output the answer in base64 so we can paste it in browser
        //fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
        //this we have to send 
          err = conn.WriteMessage(1,[]byte(Encode(*peerConnection.LocalDescription())))
        // Read RTP packets forever and send them to the WebRTC Client
        inboundRTPPacket := make([]byte, 1600) // UDP MTU
	 for {
                n, _, err := listener.ReadFrom(inboundRTPPacket)
	//	naudio, _, err := listeneraudio.ReadFrom(inboundRTPPacket)
	//	log.Println(naudio);
                if err != nil {
                        panic(fmt.Sprintf("error during read: %s", err))
                }

                if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
                        if errors.Is(err, io.ErrClosedPipe) {
                                // The peerConnection has been closed.
                                return
                        }

                        panic(err)
                }

		//if _, err = audioTrack.Write(inboundRTPPacket[:n]); err != nil {
                  //      if errors.Is(err, io.ErrClosedPipe) {
                                // The peerConnection has been closed.
                    //            return
                      //  }

                        //panic(err)
                //}
        }
}
func ByteSlice(b []byte) []byte { return b }
   type Message struct {
    x string
    y string
}
var upgrader = websocket.Upgrader{ReadBufferSize:   1024 ,
        WriteBufferSize: 1024 ,
         // Resolve cross-domain problems 
        CheckOrigin: func(r *http.Request) bool {
             return  true
        },}


func main() {
mp := map[string]string{
	"0": "1460",
	"1": "1462",
	"2": "1464",
}
     // set the vars from args
     log.Print(os.Args)
     port= os.Args[1]
     fport=os.Args[2]
     display=os.Args[3]
     id=os.Args[4]
     aratio := os.Args[5]
     


      
     // set the output file for the logs
     file, err :=os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
     if err != nil {
         log.Fatal(err)
     }
     log.SetOutput(file)
     


      time.AfterFunc(time.Duration(3000) * time.Second, func() {resp, err := http.Get("http://localhost:8090/close/?id="+id)
                if err != nil {
                  log.Fatalln("error in close matchmaker call",err)
                }
                log.Print(resp)
                os.Exit(2)})



     // handle /sdp endpoint
     http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
        // Upgrade upgrades the HTTP server connection to the WebSocket protocol.
        conn, err := upgrader.Upgrade(w, r, nil)
        //log.Print("connection: ", conn)
        if err != nil {
            log.Print("upgrade failed: ", err)
            return
        }
        defer conn.Close()
          //err = conn.WriteMessage(1,[]byte("hello"))
        // Continuosly read and write message
        for {
            mt, message, err := conn.ReadMessage()
            if err != nil {
                log.Println("read failed:", err)
                break
            }
            log.Print("message recived  ",mt,message)
            aman(string(message),conn)
        }
    })

    // handle /link endpoint
    http.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
        // Upgrade upgrades the HTTP server connection to the WebSocket protocol.
        conn, err := upgrader.Upgrade(w, r, nil)
	// log.Print("connection: ", conn)
        if err != nil {
            log.Print("upgrade failed: ", err)
            return
        }
        defer conn.Close()
	defer func() {
                if err = conn.Close(); err == nil {
                        log.Print("websocket connection closed")
                }
        }()
        //err = conn.WriteMessage(1,[]byte("hello"))
        // Continuosly read and write message
        for {
            mt, message, err := conn.ReadMessage()
            if err != nil {
                log.Println("read failed:", err)
	
		break
            }
	    log.Print("message recived  ",mt)
            link:=string(message)
            // log.Print("link is ",link)
	    //dis:=":0."+display
	    dis:=":0."+display
	    // log.Println(dis)
            os.Setenv("DISPLAY", dis)
	    // log.Println("display is ",os.Getenv("DISPLAY"))
            // log.Println("display " ,display)
	    Var, err := strconv.Atoi(display)
            Var=Var*600+500
	    str:= strconv.Itoa(Var)
            log.Println("starting position", str)
	    log.Println("map: ",mp["0"],"link: ",link)
	    log.Println("current coordinate: ",mp[display])
	    // output,err:=exec.Command("google-chrome",  "--window-position="+mp[display] , "--window-size=600,400", "-kiosk","--use-gl=desktop", link, "--user-dir=/home/ubuntu/"+display ).CombinedOutput()
	    cmd:=exec.Command("google-chrome","--gpu" ,"--gpu-launcher" ,"--in-process-gpu" ,"--ignore-gpu-blacklist" ,"--ignore-gpu-blocklist","--ignore-gpu-blacklist","--enable-gpu-rasterization","--no-sandbox","--in-process-gpu","--enable-unsafe-webgpu", "--enable-features=Vulkan,UseSkiaRenderer","--enable-vulkan","--kiosk", "--window-position=0,0" , "--window-size="+aratio, "","--use-gl=desktop", link, "--new-window", "--user-data-dir=/home/ubuntu/go/pkg/mod/github.com/pion/webrtc/v3@v3.1.43/examples/rtp-to-webrtc/"+display,"--disable-dev-shm-usage")
	    log.Print(aratio,cmd)
	    stdoutStderr, err :=cmd.CombinedOutput()
	    log.Println("chrome run sucessfully")
	    log.Println("Chrome output for coordinate : ",mp[display]," -> ", string(stdoutStderr))
	    if err != nil {
               log.Print("failed to start chrome tab: ", err)
               return
            }

            // log.Println("Output:\n%s\n", string(output))
	    // log.Println("aman error",err)
            //if err != nil {
            //    log.Println("error in command%s", err)
            //}//else
	    log.Println("command run success")
       }
    })

    // handle /input endpoint
    http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
            conn, err := upgrader.Upgrade(w, r, nil)
	    if err != nil {
               log.Print("upgrade failed: ", err)
               return
            }
	    defer conn.Close()
	      //fmt.Printf("response in inpput is ",r)
	    for {
               mt, message, err := conn.ReadMessage()
               if err != nil {
               log.Println("read failed:", err)
	       resp, err := http.Get("http://localhost:8090/close/?id="+id)
                if err != nil {
                  log.Fatalln("error in close matchmaker call",err)
                }
                log.Print(resp)
		os.Exit(2)
               break
               }
               log.Print("message recived  type ",mt,reflect.TypeOf(message))
               var msg  map[string]interface{}
               json.Unmarshal([]byte(string(message) ), &msg)
               fmt.Println(msg)
	       var typ string
	       typ=msg["type"].(string)
	       //cmd:=exec.Command("xdotool", "windowsize",mp[display], "100%", "100%")
	       // err=cmd.Run()
               //      time.Sleep(1 * time.Second)
	       //if err != nil {
		       //  fmt.Printf("error in command%s", err)
               //}
	       if typ == "click" {
		    var x string
                    var y string
                    x=strconv.Itoa(int(msg["x"].(float64)))
                    y=strconv.Itoa(int(msg["y"].(float64)))
	            log.Print("clicking on",x,y,typ)
		    cmd2:=exec.Command("xdotool", "mousemove", "--screen", display,x,y,"click","1")
	            err=cmd2.Run()
	            //time.Sleep(2 * time.Second)
	            if err != nil {
                       fmt.Printf("error in command%s", err)
                    }
	      }else if typ == "mousemove" {
                    var x string
                    var y string
                    x=strconv.Itoa(int(msg["x"].(float64)))
                    y=strconv.Itoa(int(msg["y"].(float64)))
		    log.Print("mousemove on",x,y,typ)
                    cmd2:=exec.Command("xdotool", "mousemove", "--screen", display,x,y)
                    err=cmd2.Run()
                    //time.Sleep(2 * time.Second)
                    if err != nil {
                      fmt.Printf("error in command%s", err)
                    }
             }else if typ == "mouseup" {
                    cmd:=exec.Command("xdotool", "mouseup","1")
                    err=cmd.Run()
                    //time.Sleep(2 * time.Second)
                    if err != nil {
                      fmt.Printf("error in command%s", err)
                    }
	     }else if typ == "mousedown" {
                    cmd:=exec.Command("xdotool", "mousedown","1")
                    err=cmd.Run()
                    //time.Sleep(2 * time.Second)
                    if err != nil {
                      fmt.Printf("error in command%s", err)
                    }
	     }else if typ == "scroll" {
                    var x float64
                    var y float64
                    var z float64
		   // logs.Print("received is ",msg)
                    x=((msg["x"].(float64)))
                    y=((msg["y"].(float64)))
		    z=((msg["z"].(float64)))
                    log.Print("scroll on",x,y,z,typ)
                    var cmd2  *exec.Cmd
                    for i := 0; i < int(math.Abs(y)); i++{
                       if y > 0{
			 cmd2=exec.Command("xdotool", "click" ,"5")
                       }else{
			 cmd2=exec.Command("xdotool", "click" ,"4")
		       }
		       err=cmd2.Start()
                       //time.Sleep(2 * time.Second)
                       if err != nil {
                         fmt.Printf("error in command%s", err)
                        }
		    }
            }else if typ == "key" {
		    log.Print("received is ",msg)
		    var x string
                    x=((msg["x"].(string)))
		    if(x=="Enter"){
		       x="Return"
	            }else if x=="Backspace"{
		    x="BackSpace"
		    }else if x==" "{
                    x="space"
                    }
		      
                    //x=msg["key"].string
		    cmd:=exec.Command("xdotool", "key",x)
		    
                    err=cmd.Run()
                    //time.Sleep(2 * time.Second)
                    if err != nil {
                      fmt.Printf("error in command keyboard%s", err)
                    }
             }
          }
      })
    http.ListenAndServe(":"+port, nil)
}
