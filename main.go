package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func clone(url string, folder string) bool {
	zip_name := url + ".zip"
	os.Chdir(folder)
	// 下载
	fmt.Println("Start PKGBUILD Download.")
	resp, err := http.Get("https://github.com/cleararch/test_package_store/archive/refs/heads/" + zip_name)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	file, err := os.Create(zip_name)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return false
	}
	fmt.Println("Finish PKGBUILD Download.")
	// 解压缩
	zipReader, err := zip.OpenReader(zip_name)
	if err != nil {
		return false
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		fpath := filepath.Join(folder, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return false
			}

			inFile, err := f.Open()
			if err != nil {
				return false
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return false
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return false
			}
		}
	}
	return true

}

func package_install(package_name string) bool {
	clone := clone(package_name, "/tmp/")
	if clone == false {
		return false
	}
	os.Chdir("/tmp/" + "test_package_store-" + package_name)
	press_y := exec.Command("echo", "y")
	install := exec.Command("makepkg", "-fsi")
	install.Stdin, _ = press_y.StdoutPipe()
	_ = install.Start()
	_ = press_y.Run()
	// watching_install := func() {
	// 	var output, sop bytes.Buffer
	// 	install.Stdout = &sop
	// 	for {
	// 		install.Stderr = &output
	// 		if string(sop.Bytes()) != string(output.Bytes()) {
	// 			_ = exec.Command("clear").Run()
	// 		} else {
	// 			sop = output
	// 		}
	// 	}
	// }
	// go watching_install()
	// fmt.Println("Start Watching Succuesfully.")
	err := install.Wait()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func package_remove(package_name string) bool {
	press_y := exec.Command("echo", "y")
	remove_packeage := exec.Command("sudo", "pacman", "-Rscun", package_name)
	remove_packeage.Stdin, _ = press_y.StdoutPipe()
	_ = remove_packeage.Start()
	_ = press_y.Run()
	err := remove_packeage.Wait()
	if err != nil {
		return false
	} else {
		return true
	}
}
func main() {
	package_ins := flag.String("install", "foo", "安装包。")
	package_rem := flag.String("remove", "foo", "移除包。")
	frontend_use := flag.String("frontend", "0", "程序使用，无输出")
	// package_sea := flag.String("search", "foo", "Search package.")
	flag.Parse()
	if *frontend_use == "0" {
		if *package_ins != "foo" {
			if package_install(*package_ins) != true {
				fmt.Println("无法安装" + *package_ins + " 包。")
				os.Exit(1)
			} else {
				fmt.Println("安装 " + *package_ins + "成功")
				os.Exit(0)
			}
		}
		if *package_rem != "foo" {
			if package_remove(*package_rem) != true {
				fmt.Println("未能移除" + *package_rem + "包")
				os.Exit(1)
			} else {
				fmt.Println("移除 " + *package_rem + "包成功")
				os.Exit(0)
			}
		}
	} else {
		if *package_ins != "foo" {
			if package_install(*package_ins) != true {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
		if *package_rem != "foo" {
			if package_remove(*package_rem) != true {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	}
}
