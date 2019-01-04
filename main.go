package main

import (
	"fmt"
	"net/http"
	"text/template"
	"webssh/asset"
	. "webssh/server"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "webssh",
	}
	index, _ := asset.Asset("www/index.html")
	t, _ := template.New("index").Parse(string(index))
	fmt.Println(t.Name())
	t.Execute(w, data)

}
func init() {
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}
	http.Handle("/home/", http.StripPrefix("/home/", http.FileServer(&fs)))
}
func main() {

	fmt.Printf("-------启动 index %s\n", Lan+":80")
	http.HandleFunc("/", index)
	http.ListenAndServe(Lan+":80", nil)

}
