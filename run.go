package main

import (
	"flag"
	"fmt"
	"net/http"
	"text/template"
	"webssh/asset"
	. "webssh/server"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Wlanip string
	}{
		Title:  "webssh",
		Wlanip: Wlan,
	}
	index, _ := asset.Asset("www/index.html")
	t, _ := template.New("index").Parse(string(index))
	fmt.Println(t.Name())
	t.Execute(w, data)

}

//跳转使用正向代理
// func Index(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "http://"+wwwhost+"/home/www", http.StatusFound)
// }

func init() {
	flag.StringVar(&Lan, "lan", "127.0.0.1", "输入内外网ip")
	flag.StringVar(&Wlan, "wlan", "127.0.0.1", "输入外网ip")
	flag.StringVar(&Logdir, "logdir", "./", "用于存储日志")
	flag.Parse()
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}
	http.Handle("/home/", http.StripPrefix("/home/", http.FileServer(&fs)))

}

func main() {
	go Startsocket()
	go Starthttp()

	fmt.Printf("-------启动 index %s\n", Lan+":8090")
	http.HandleFunc("/", Index)
	http.ListenAndServe(Lan+":8090", nil)
}
