package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)
var wg sync.WaitGroup
func main(){

	start := time.Now()
	var inputDir string
	var outputName string
	flag.StringVar(&inputDir, "i", "./dvr-file/flv", "input file dir")
	flag.StringVar(&outputName, "o", "out.mp4", "output file name")
	//解析命令行参数
	flag.Parse()

	exist, err := pathExists(inputDir)
	if err != nil {
		fmt.Printf("get dir error!: %v", err)
		return
	}
	if !exist {
		inputDir = os.Args[0]
	}
	inputDir, _  = filepath.Abs(inputDir)
	fmt.Println("argv: ", inputDir, outputName)
	if err = flvs2mp4(inputDir, outputName); err != nil {
		fmt.Printf("flv to mp4 error!: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Println("Running time:", elapsed)
}

// 判断文件或目录是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 命令行调用
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	//fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

// 视频格式转换
func videoConvert(in string, out string) {
	defer wg.Done()
	//fmt.Println(in, out)
	cmdStr := fmt.Sprintf("ffmpeg -i %s -loglevel quiet -c copy -bsf:v h264_mp4toannexb -f mpegts %s", in, out)
	args := strings.Split(cmdStr, " ")
	msg, err := Cmd(args[0], args[1:])
	if err != nil {
		fmt.Printf("videoConvert failed, %v, output: %v\n", err, msg)
		return
	}
}

// 视频合成
func videoMerge(in []string, out string) {
	//fmt.Println(in, out)
	cmdStr := fmt.Sprintf("ffmpeg -i concat:%s -loglevel quiet -c copy -absf aac_adtstoasc -movflags faststart %s",
		strings.Join(in, "|"), out)
	args := strings.Split(cmdStr, " ")
	msg, err := Cmd(args[0], args[1:])
	if err != nil {
		fmt.Printf("videoMerge failed, %v, output: %v\n", err, msg)
		return
	}
}

func flvs2mp4(inDir string, outFile string)(err error) {
	tsFileDir := filepath.Join(inDir, "tsfile")
	if err = os.RemoveAll(tsFileDir); err != nil {
		return
	}
	if err = os.RemoveAll(outFile); err != nil {
		return
	}
	if err = os.Mkdir(tsFileDir,0666); err!=nil {
		return
	}

	infiles, _ := ioutil.ReadDir(inDir)
	for _, f := range infiles {
		if filepath.Ext(f.Name()) == ".flv" {
			tsfileName := filepath.Join(tsFileDir, strings.TrimSuffix(f.Name(), ".flv") + ".ts")
			wg.Add(1)
			go videoConvert(filepath.Join(inDir, f.Name()), tsfileName)
		}
	}
	wg.Wait()

	tsfiles, _ := ioutil.ReadDir(tsFileDir)
	tsfileNames := make([]string, 0, len(tsfiles))
	for _, f := range tsfiles {
		if filepath.Ext(f.Name()) == ".ts" {
			tsfileNames = append(tsfileNames, filepath.Join(tsFileDir, f.Name()))
		}
	}
	videoMerge(tsfileNames, outFile)

	return
}
