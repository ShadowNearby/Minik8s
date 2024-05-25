package function

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"minik8s/utils"
	"os"
	"os/exec"
	"strings"
)

const ImagePath = "shadownearby"

/* build image for function*/
func CreateImage(path string, name string) error {
	// 1. create the image
	// 1.1 copy the target file to the func.py
	name = fmt.Sprintf("%s/%s:v1", ImagePath, name)
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
	// 1.2 create the image
	//sudo usermod -aG docker $USER
	cmd := exec.Command("docker", "build", "-t", name, utils.RootPath+"/image/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Error("[CreateImage] docker create image :", name, "at ", utils.RootPath, "/serverless/image/", " error: ", err)
		return err
	}
	// 2. save the image to the registry
	err = SaveImage(name)
	if err != nil {
		log.Error("[CreateImage] save image error: ", err)
		return err
	}
	log.Info("[CreateImage] create image success")
	return nil
}

/* save the image into the registry*/
func SaveImage(name string) error {
	// 2. push the image into the registry
	cmd := exec.Command("docker", "push", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error("[SaveImage] push image ", name, " error: ", err)
		return err
	}
	return nil
}

/* find the image */
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

/* delete the image */
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

	if FindImage(name) {
		cmd := exec.Command("docker", "rmi", name+":latest")
		err := cmd.Run()
		if err != nil {
			log.Error("[DeleteImage] delete second image error: ", err)
			return err
		}
	}

	log.Info("[DeleteImage] delete image finished")
	return nil
}

/*RunImage to run image for function*/
func RunImage(name string) error {
	imageName := fmt.Sprintf("%s/%s:v1", ImagePath, name)
	// 1. run the image
	cmd := exec.Command("docker", "run", "-d", "--name", name, imageName)
	err := cmd.Run()
	if err != nil {
		log.Error("[RunImage] run image error: ", err)
		return err
	}
	return nil
}
