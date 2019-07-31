// libtool package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "autoconf"
	vers     = "2.69"
	ext      = "tar.gz"
	url      = "http://ftp.gnu.org/gnu/autoconf"
	hashType = "sha256"
	hash     = "954bd69b391edc12d6a4a51a2dd1476543da5c6bbf05a95b59dc0dd6fd4c2969"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while configuring ...")
	f := "bin/autoreconf.in"

	re, err := regexp.Compile("libtoolize")
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	newContent := re.ReplaceAllString(string(content), "glibtoolize")
	err = ioutil.WriteFile(f, []byte(newContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("./configure", "--prefix="+*destdir)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func build() {
	fmt.Println("Waiting while compiling ...")
	cmd := exec.Command("make", "-j" + pkgsrc.Ncpu)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func install() {
	fmt.Println("Waiting while installing ...")
	cmd := exec.Command("make", "install")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
	install()
}
