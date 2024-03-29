package assets

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

var staticFiles = map[string]*staticFilesFile{
	"CSIOperator.md": {
		data:  "",
		hash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		mime:  "",
		mtime: time.Unix(1574851373, 0),
		size:  0,
	},
	"CronHPA.md": {
		data:  "",
		hash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		mime:  "",
		mtime: time.Unix(1574851373, 0),
		size:  0,
	},
	"GPUManager.md": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xbcV_S\x1aY\x16\u007f\xefOq\xb7|\x99\xa4\x00\x01\xff$\xf1-\xb3\x99\xa2\xb6R\xb3\xcb\u058c\xfb2\xb5U2͍CE\x1a\x87\x86l\xa5\xca\aT@\x10\x11fF\xa2\"QI\xfcCtlp\x8c\xd8t\v~\x98\xf4\xb9\xf7\xf2\xe4W\xd8\xea\xbeM\xab\x89\x0f\xbb/C\x95U\xed\xbd\xf7\x9c\xf3;\xbfs\xce\xefޡ!\x14\bN\xba\xbf\rI\xa1i\x1cg͏dcU\x10\x86\x86\x86\x10\xd5ӆ\xde6\xf4\x02Ջ\x82\x10\bN\"\xfb\x10)\x95\x8dޖ\xa1\xa6\f\xf5\xe8\xe9̌;\"\xb9\xff!aZM\a\x82\x93T\xa9\xd3r\x166\x1b.\x04;\x9a\xa1\xad>O\xfe\x88\xe3\x12N\xc8\xe8\x19~\x15\x11qp&9\x1d\x91H\xe9WCo\xd33\x9d\xea;\xa0l\xd3Ֆ\v\xb1\xe6\xbecm\x87в\x90\xcb\xf63E\xe8\xb4!sjhǁ\xe0\xa4\xcb\xc4KV\x96\xc8\xee\x12\xd9\xd9g\xcdw.\x04J\a6\x1b\xec\xaa\xcc\xea+\x90/\xdaH\xb6TȜ1\xa5\a{K\xb0\x94\x85\xbd%z\x92\x87\xe5\x1d\xb6\xd8u!\xb2\xd6$+\xf3t\xadAr\x17Pk\f@b\xb9\xbf\x95\xa5\xbd=C=1\xbaWt\xad\x11\bNr\x0f\x1eAp\xc0\xc1J\x06\xca\xc7p\xb0`\xa8\x05\xeepB\x10\xdc\xe8\xe1CR\xf8\x8d\xe4\u007f\xe1\x90\x1f>\xbc\xbe\xac\xf248\x11&\xea[\xfb\x03$\xd0\xfd͆a-\xb3v\x86]-\xf9\xe0r7\x10\x9c\x84b\x9dVӴ\x9a\x06m\x8d\xae\x99\x94\x96\x9a\x86\xbe\xdfO\xe5I\xe1\x03wƮ\xb6\xc9\xea>\xa9\xa5\xe0\xea\x98V\xd3\xe4M\a.Kܕ\x83\x9d\x833\xfdY\x1cZ\xc8 3\xcf\x14\x95G6Qj{\xd0*\x19ڪ\x19\x19\x8auv\x9e&Z\x99\a7t\x1d\x96\xeb.\x04\x9d6\xcf矱\uf32b\xb7\xac9\xcf\x1dspNQL\x14V]\xac8vn\x16n\xa6\\\xf5\xd7\x15C\xd5IM3]\x1f7\xa1\xf4\xfe\xab\xbe\xbe\xc1\x94=C\xd5\xc6\xc6\x1f=~@\xab\xe9\xe1(N\xc4#\xa2\xcc.\x9a\xd0K\xbb\xb8\xad\xa1j\xc1x,\x8a\x13?\xe1\xa4\xcca܉\xb8\xd6\xeeoe\a\x9c\x0e'\xe5\xd04\xe6\xf6\xdc\xdc.D\xa9ɚ:I\x1d\x9a\x94Z=C\x97ې=\xe7xy\xd7\xf7\x17\x1b\xb4\xfb\a\xd4\x1a\xbc\x13 \x9by\xe94\a4;\xec\xb4.\b\xb7w\xf9\xf9[#\x84\x9e\x86\xc3\ue604\\\bZ\xd9/OZ\xb9\x14\xbe\xf49\x87>_\x83r\x91\x1e\xb6\x90\xfd\x9bC\xf4T\x87\xed\x02\xbaY\x00]c\x8a\xd2\u007f\x97\xa6\x95M^04\x87H>\x05\xa7\xdb\u007f\x0fE\xb1<\x1b\x12\xb1\x8c\xe6\x849\xe4\xbe\xef\x87n\xaf;\xff\xa0\xbb\xa7-\xf3\xe9٤;ʳs\x87C8\x1a\x93d\x9c@s\xe8\x99\xf5\xfd\x1dN\x98hH\xb3Ė\x17\xe8B\xc7Gvտ\x9as\xea\vD \x9b\x81\x93\rēs˯\xe5\x04\x8e\"\xc7\xe7\xcf\xc9X\"\xe4\x0e\x85\xa3\x11Y\x8e\xc4$+\xa9gxv&\xf6:\x8a%3\xc0=\xae\xec\xdc?\xf3'ܕ1>\xbcP\xd3\xc8f\xd3*\xd7=\x03\xce\xe5\u009c\n\xab\u007f\xc9z{0_\xec\xf0=dϞ\xfe\x8d)'T_4U\xc3rd\xa8't\xed\x8c5/h\xe1wz\\\x80b\xbd\x9f\xda%\xb92\xa7\x9e\x9c\x1f\xb1\x8f\x1d\xb3\x15\x173\x90=\xbf\xbe\\a\xca\a\xa6ԩ\xb2\xce\x0f@o\x1dr-\xc8d \x97\x85\xdc\a\xba\xd6\xf8\x94Z\xf8\x1cx\u007f\xb3\f\xb96y[7\xf4\xb6 \xf8<\x96(Z:|WM\xf1]9\x1dH(\x94\x8e\xd8b\xd7\x16\xc2Z\x83\xab\xdc탴\x9a\xfe\xfe\xf97\xb4\x9a\xf6y|ޛ\x86\xa3\xf9\x1c\xa9\xfdnt\n\x86\xbal\xa1\xf2{\xccz:*d\xa8)Ȝ\x92Z\xde\xe7\xf5\x1a\xea\x11\x14+FהY\x9e\x99\v\x19z\x86\xc7\xf2\xba}\xe6l\xb5J\xa4҂b\xdde\xe8\xfbPZ\xb6\xd6RERi\x99\xfcT>\xf2M\xa6\xf4hW\xf1 \xb2у\x93\r\x9b\xc5\r\xb3\x00\xfe\xb1\xf1o#_\x1b\xaaf\xeaZ\xabdj\x90u\xd98q\xb9\xc0qC\v\xef\x88\a9\x82\xed\f\";\x98'\xa7\v\xce\xfcq\xd96U\xa3\xa6\xc1v\x817\xeb=5\xe0\x8eț\x0e9\xabpM\x00%\xcf\xdeg\xacz\xd0M\x1d\xba\x95\x1f\xb8z\x90Z\x11\x96\xebd\xf5\x10rm(\xb5\xfe\xfd\xd5O\x89Ĭ<1<,\xc6$96\x83=?\x8b3\xb1d\xd8#ƢÉ\x97\xd8\xff`\xc0-\xd4\x1apq`\xf4\x0e\xa1y\xc9rGd\xb7d\xa8'ח+P\xac\xc0\x92\xfe)U&\xf9\x0fpZ\xe1\x95\xff\x94\xfa\xc5즫-\xc8\xec\xdf^\xe7\xf7Q\xbf~\xde\u007f\xfbn@\x02\xbf\x1a\xfa\xb5\x14;\x98\xe7\xa8\xf9e\xe8\xf4\x97E\xc6\xf5\xe5\n]\xe8ؑ\u07b4@\xd7x\f8X\x80\xad\xde\xf5eU@\b\x99\u007f\u007f\xf9\xe1&\xa7h(\"\xd9\tE\xa2\xd3VN\xf1\xd0\u007f\x86\xf1\xe3q\xbf\xffɈ\x17\x8f\x8d>\t?\xc6\xe2\xe8\xd8#q\xcc\xe7\xf3\x8e\xfa\xfd~<\xe6\x1b\xf5\xccJ\xd3\x0fl\x1as[\xa0kT\xcf\xd2?~\x05\xed\x80#\x83\x8b}\xa3[c\x1fwX\xb7k\x1d\x1bB\xdf?\xff\xc6\xe1\x94\xdb\xfc\x99\xd4s\x8a8!w贐8zA\xabi\x87L\xbeu;\x93A=\xbe\xdc\xe1]\x1f\bN\xf2Q\xbfú0\xea\xf9_(\x17\xbd\xe3\x8f\x1f\xf9\xc3a\xd1\xfb\"\xf4\xe2\xc7\x17O\xfc#\xa3c\xa27\xfc\xe4\x85\u007f\x1c\x8fb1|C\xf9\x10z\x1d\x8a\xce\f8\xe4\xf7\xbc9\x82\xebmHm\xdai\xd5\x1a\xe6\x19C\xd5\xecW\xd4\x00\x9f}\xf9[\xc30\x98s\xb2\xab\xf2O\xc76\x8e\xe5X2.bC]\x86\xfa1d7\xa7\x12X\x12\xb1\x94\xb0о\x12\x93\xe1\x90[\x8c\xc5\xf1\x94\xebΜ\xff_\xf6Q\x1c\x8d\xc5_O\xb9\xcc\xe7\x06\xc7c\xbfRh5\x1d\x1c\x1d\xbcq\xa6\xa6\xa6\x04!4\x1b\xf9\x17\x8e\x9b\x97\xc9\x04z\xe5\x13\x84\x97\x11)<\x81\x82\xb1\xb0 x<\x1eA\x90g\xb18!\bbLJ\x84\"\x12\x8e\xcb\xd6\xcbM\nE\xf1\x84y#\t\xc2\x00\x8f\xb9q\u007f&\x13\xc8\xe7\xf5\xf2h\x038^\xcf\b\a\xe4Bc\x81\xc8\xd7<Q\xe7\xd9\xf6gB\x1b\xf1\u07b7\xc7\t\x9c@~\x8e\xfb\xbf\x01\x00\x00\xff\xff\x9fp\xa9\x90\xf3\v\x00\x00",
		hash:  "1f2c694cb54e80cc9d0bd63ff1f1362219c003f63db94cf2c73d6ba21274faf7",
		mime:  "",
		mtime: time.Unix(1574851373, 0),
		size:  3059,
	},
	"PersistentEvent.md": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xd4VKs\x1aG\x17\xdd\xf3+\xfa+m>WY\b\x900\x90\x9d\xe3\xf2\xcaI\xcaUNV\xae,\x10\x8c\x1c\x97$P\fN\xcaU,\x86\x11\x83x\nH\xf4@\x80\x8c\xb0\x8d\x84\x1f0(\x960\xcc\xf0\xf81\xf4\xednV\xfa\v\xa9\x99F\x88 \xc5e/Ê\xea\xee{\xef\xb9\xe7\x9c\xdb=ss\xe8\xa1\xf0,\xf04\x10\x14|\xc1\xfb\xbf\t\xbe S\xceH~\xdbd\x9a\x9b\x9bCT\x8b`\xad\x85\xb5$\xd5\xd2&Ӄ\xe7\xcb\xc23\x9f\x10\x14\x02\xc88\x19@\x90\x92I\xf2\x03V\xa3\x88\xef\xe9[\xa3b\x94\xf6\xdf\xd0B\x84\r\xb2\xac\x92\x82?R\x90\x8d\xd0S\x8d\x9dG\x88\x9a\xd5כ\x9b\xa0\x1e\x93M\x19\xa2\xe7\x17\xdd\x14(\x1d\xaa\x9d\x91\xc4\x1b\xac\xaa\x90˓\xfd*\xb4\xdb\xecD\x02\xa5<\t\x81L\x9e\x14ϰV\x85L\x02\x1a\x05\xdcK\x8f\xf6\x1b\xa3Wy8\xdc\"\xa58\xb4\x1b\x90x;\x14\xa5\x1f\x1f\xdcGdG!\xa90n\xabD\xaa\xd1B\x84\xc4ER\x8asL#9M{\r\xac&\xb1\xd6\xd2\xcftdH\xedA\xa2\xcc6{:\x8e\xae\bY\x85\x94>\xf0\x15\xc8n_tS\xb8[\x80f\x94\xa7\xe2Ix84^\x92\xfd\x16(]\xd8R!\xd6\xe4\x99u\xa0\xf5<H5\xfa^\x19\xa3a\x83\xfc\x18Po@wjL\uecc6\x82\xd5\x1c\xc9dq\xbfH\v\x11\xf4\xf0\xee\xddG\x88\x94Ґ\xa8\x90\xd8\x1etE\xa2fYO\xd1k(\x9d1\xd6\xf30i\x9e\xb1A\x91UR\xe4\xb5H\xcf^\rE\x89\v4ڬ\xd1\xde_P\xaaqp\x10\x95W'\"\x81\xd2a\xa7\x15\x93iz\x97\x9f\x9fQ\x1c\xdd\xf5z\xe7\xfd>t\x1bA3z\xfd4֪\xb8\x9d\xbc\x9e7\x84f\xd7 \x9b\xa6'M\x14B\xf4T\x83\x97I\x14B#-\xcf\x1ao }\xa47o\xa8\x89B\x88\xc4E8}\xf9\x83{]\bl\xb8=B\x00\x85L!4?\xfbC\xc6\x1a\x9a\xdeA3\xc7\xf4\xc0\xe0\xaa0\xbf1\xe9g^\xd0\x1b\ny\x85\x8d5\xff\x8bu\xfd\xaf\xc5l#G\xed{\x0f\u007f\xbam\xb5X\xbe\xff\x16\xa22\xd4\xf3!\x1d\xf9|\xe0E (\xac\x87t&g\x87\x80\xcb\x05%\x95\x1c(\xd3\xc6\x1f\v\x92W\xa69\x1a\xdbt\xa7\f\xb9\x1e\xe4j\xa4T\x1e\x8aa\xbeȽ>\x14\xc3Е\xa0݆\\\x82Տi=\xceݏ\xd5\x13\xbaS\xd6G\xa2фޮn\u008c\x82\xb5\xeaH,\xb0\xc1ָ֧S\x90\xabX\x8d\xb2\x93\xd7\xd3E!\x93\xe3\xb1XM\x92MY\xf7j?\t\xc7\x12\xafHkIP3\x10.\xc0\x96\xaa\xdb2\xfa\x91\xd6\xf7Hl\x8f\x892)g\xb1Z\xdc\xf0{\xb1\xb6\xcd\x12\x12\x95:\xb8]g\xafޏ\xc4\xec\xf8\xf0v\x19\x8aG\x86Ů\x04\xe6BbM\xe6\xceŃC\xba{\x80\xdb\"n\xbf\x83f\x86\xec\xb7h!\xb2:C\xd3P\x94fyE\xe3p5\n\xcd\xe8\xd4}2;\x94|\x8c\n\x11\x88\xeb\x83u5\xa3r\x985\xdaD\xaaq\x86f\x93C3:\xa1h<\xa7\x97\x13J\xa4\x1a\xdbzGJ\xf1\x99!\xbdI\xfc\xd1A\x16b-rX\xc1Z\xcbd\xb2\x9a\x114\xe2\xec\xb5|C5nl^\xf3F\xa3\xe9\xda\x1aF0J\xd9\xcc\bk2\x94jV\xb3\x93\xc6c\xa4\xf4\xc1\x98\xad\xc4?\xb8\xe3\xd9\xf8\xad\xf1/\x00\xb9;\xc9^\x87|\xdc\xe5\xf7\x00G\b\x9d\x16k\xf4u\xc2.;4\xe0\xd3\x03\rz\xbb\x8f\xa1с\x83\xda\xf8\xaa\xd9>\x81X\v2͟\xff\xffK0\xb8\x11\xf8fa\xc1\xe3\xf7\x05\xfck\x82\xf9WϚ\xff\xb9\xd7\xec\xf1\xaf/\x04W\x05ۭK\xe4P\xaa\xc1\xa7c\xdc?\x01\xa5\xcbb\xef\xc8Q\x06\xb7\xeb\xba&\xe9]\xd8҆b\x96\xc4\xdf\xc2\xe9.\u007f2\x86b\ue89bb\x83\"\xc8\xd5\xe9uڨ\xd0ltT9\x1f\x1d\xf2KlьFb\x9c$ߎJ\";\x0e\xf3>h!2+\x89A\xcaE7E\xa5θ\xda^\x134\x95ׁc\t\x8a\xfd\x8bn\xc1d\xfa\xdf㫎\xd6\xddO}\xe3v\x9e\xae?1:z\xe6\xfe}\xc1\xe6\xb1z]\x8e\xa5e\xfb\xb2}i\xd1\xe1\xf6:m\x8b\xcb\xceE\xc1-\xd8\xef\xd8\x05\x8f\xcbn\xde\xf0=\xb9e2-\x99\xd1\xcd\x0fƌ{\xe6\x10)\x9e\xe9h\xfe\xab\x9cs\xf8\x9f\xe1\x9c\x1f\x9eٽ\x12\xc2\b\xe7T]\x93\xe3\xcbI\x84R\xed\xdew\x8f&\x1c\xf1w\x8eGM\x13\xaa\u007f\x1e\f\xf6\xaf\x13\xfa\x19F=k\x81\xaf \xd4ȯ\xb7n\xb0\xc6\x1b\xba\x91\x80ɛ?\t\xb9\xe8\xa6H\xfcO\xfd\t\xe7 \x8d\x16\xf8\xdde\xd4\xff\x02o\nv\xbb\xc5\xe5X\xb2\xaf\xac\xb8\xed6\xefʢ\xcb氹\x1d\x1e\x97\xc3\xear\xb8\x9d^\xe7؛\x13\x05\x99R\x9d.\xcfE\x99\xe6\xee\xeb\xa7\xc3aq\xacح.\x8f}\xe5\x8e\xd3v\xc7auY\\\x82e\xd1jw:\x9c\xcbv\xc7\xe5t\xfc\x1d\x00\x00\xff\xff\x18\n\x0e\xb18\n\x00\x00",
		hash:  "9e1b894cd8de723613d96dcf222679e753f3c27ce166b75f66e2ba7c1576b018",
		mime:  "",
		mtime: time.Unix(1574851373, 0),
		size:  2616,
	},
	"TappController.md": {
		data:  "",
		hash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		mime:  "",
		mtime: time.Unix(1574851373, 0),
		size:  0,
	},
}

// NotFound is called when no asset is found.
// It defaults to http.NotFound but can be overwritten
var NotFound = http.NotFound

// ServeHTTP serves a request, attempting to reply with an embedded file.
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	f, ok := staticFiles[path]
	if !ok {
		if path != "" && !strings.HasSuffix(path, "/") {
			NotFound(rw, req)
			return
		}
		f, ok = staticFiles[path+"index.html"]
		if !ok {
			NotFound(rw, req)
			return
		}
	}
	header := rw.Header()
	if f.hash != "" {
		if hash := req.Header.Get("If-None-Match"); hash == f.hash {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("ETag", f.hash)
	}
	if !f.mtime.IsZero() {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
	}
	header.Set("Content-Type", f.mime)

	// Check if the asset is compressed in the binary
	if f.size == 0 {
		header.Set("Content-Length", strconv.Itoa(len(f.data)))
		io.WriteString(rw, f.data)
	} else {
		if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			header.Set("Content-Encoding", "gzip")
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			io.WriteString(rw, f.data)
		} else {
			header.Set("Content-Length", strconv.Itoa(f.size))
			reader, _ := gzip.NewReader(strings.NewReader(f.data))
			io.Copy(rw, reader)
			reader.Close()
		}
	}
}

// Server is simply ServeHTTP but wrapped in http.HandlerFunc so it can be passed into net/http functions directly.
var Server http.Handler = http.HandlerFunc(ServeHTTP)

// Open allows you to read an embedded file directly. It will return a decompressing Reader if the file is embedded in compressed format.
// You should close the Reader after you're done with it.
func Open(name string) (io.ReadCloser, error) {
	f, ok := staticFiles[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.size == 0 {
		return ioutil.NopCloser(strings.NewReader(f.data)), nil
	}
	return gzip.NewReader(strings.NewReader(f.data))
}

// ModTime returns the modification time of the original file.
// Useful for caching purposes
// Returns zero time if the file is not in the bundle
func ModTime(file string) (t time.Time) {
	if f, ok := staticFiles[file]; ok {
		t = f.mtime
	}
	return
}

// Hash returns the hex-encoded SHA256 hash of the original file
// Used for the Etag, and useful for caching
// Returns an empty string if the file is not in the bundle
func Hash(file string) (s string) {
	if f, ok := staticFiles[file]; ok {
		s = f.hash
	}
	return
}
