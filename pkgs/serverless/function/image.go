package function

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"minik8s/config"
	"minik8s/utils"
	"os"
	"os/exec"
	"strings"
)

const ImagePath = "shadownearby"

// CreateImage build image for function
func CreateImage(path string, name string) error {
	fullName := fmt.Sprintf("%s:v1", name)
	if FindImage(name) {
		log.Info("image existed")
	} else {
		srcFile, err := os.Open(path)
		if err != nil {
			log.Error("[CreateImage] open src file error: ", err)
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(utils.RootPath+"/image/func.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error("[CreateImage] open dst file error: ", err)
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			log.Error("[CreateImage] copy file error: ", err)
			return err
		}
		log.Info("[CreateImage] copy file success")
		cmd := exec.Command("docker", "build", "-t", fullName, utils.RootPath+"/image/")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Errorf("create image %s error %s", fullName, err.Error())
			return err
		}
	}
	err := saveImage(fullName)
	if err != nil {
		log.Error("[CreateImage] save image error: ", err)
		return err
	}
	return nil
}

// saveImage  save the image into the registry
func saveImage(name string) error {
	// docker tag <old-name> <new-name>
	// e.g. docker tag serverless_test:v1 localhost:5000/serverless_test:v1
	registryImgName := fmt.Sprintf("%s/%s", ImagePath, name)
	log.Infof("tagged name: %s", registryImgName)
	tagCmd := exec.Command("docker", "tag", name, registryImgName)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stdin = os.Stdin
	if err := tagCmd.Run(); err != nil {
		log.Errorf("tag image error: %s", err.Error())
		return err
	}
	cmd := exec.Command("docker", "push", registryImgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("push image error: %s", err.Error())
		return err
	}

	return nil
}

func PullImage(name string) error {
	registryImgName := fmt.Sprintf("%s/%s:v1", config.ClusterMasterIP, name)
	nerdCmd := exec.Command("nerdctl", "pull", registryImgName)
	nerdCmd.Stdout = os.Stdout
	nerdCmd.Stdin = os.Stdin
	if err := nerdCmd.Run(); err != nil {
		log.Errorf("nerdctl pull image image error: %s", err.Error())
	}
	return nil
}

// FindImage  find the image
func FindImage(name string) bool {
	cmd := exec.Command("docker", "images", name)

	// check the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("[FindImage] get output error: ", err)
		return false
	}

	result := strings.TrimSpace(string(output))
	log.Info("[FindImage] the result is: ", result)

	if strings.Contains(result, name) {
		return true
	} else {
		return false
	}
}

// DeleteImage delete the image
func DeleteImage(name string) error {
	// if the image not exist, just ignore
	imageName := fmt.Sprintf("%s/%s:v1", ImagePath, name)
	if FindImage(imageName) {
		cmd := exec.Command("docker", "rmi", imageName)
		err := cmd.Run()
		if err != nil {
			log.Error("[DeleteImage] delete first image error: ", err)
			return err
		}
	}
	return nil
}

/*RunImage to run image for function*/
// we don't need docker interface to run image
func RunImage(name string) error {
	imageName := fmt.Sprintf("%s/%s:v1", ImagePath, name)
	// 1. run the image
	cmd := exec.Command("docker", "run", "-d", "--name", name, imageName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error("[RunImage] run image error: ", err)
		return err
	}
	return nil
}
