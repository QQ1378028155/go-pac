package worker

import "fmt"
import "os"
import "io/ioutil"
import "encoding/json"
import "github.com/QQ1378028155/go-pac/logger"
import "github.com/QQ1378028155/go-pac/conf"
import "github.com/QQ1378028155/go-pac/cmd"
import "errors"
import "strings"
import "os/user"
import "regexp"

var antPropertyTemplate string

var (
	apkReg, _          = regexp.Compile("\\.apk$") // regular expression to find the file
	unsignedapkReg, _  = regexp.Compile("-unsigned\\.apk$")
	unalignedapkReg, _ = regexp.Compile("-unaligned\\.apk$")
	workdir, _         = os.Getwd() //record the base info
)

func init() {
	antPropertyTemplate = "key.store=%s\n" +
		"key.alias=%s\n" +
		"key.store.password=%s\n" +
		"key.alias.password=%s"
}

func Run(filename, outfile string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	config := new(conf.Config)
	err = json.Unmarshal(b, config)
	if err != nil {
		return
	}
	if config.Repository == nil {
		err = errors.New("Repository must be set.")
		return
	}
	err = fetchFromRemote(*config.Repository)
	if err != nil {
		return
	}
	if config.Android != nil {
		return compileAndroid(config.Android, outfile)
	}
	return
}

func compileAndroid(config *conf.AndroidConfig, outfile string) (err error) {
	var file *os.File

	// remove the old build
	err = cmd.SyncCmd("ant", []string{"clean", "-Dsdk.dir=/usr/lib/android/sdk"})
	if err != nil {
		return
	}
	var sign bool
	// set sign true if all infomation for signer is collected
	if config.Store != nil && config.StorePassword != nil && config.Alias != nil && config.AliasPassword != nil {
		sign = true
	}

	// generate ant.properties for ant to sign the apk while compliling and packing.

	if sign == true {
		str := fmt.Sprintf(antPropertyTemplate, *config.Store, *config.StorePassword, *config.Alias, *config.AliasPassword)
		file, err = os.OpenFile("ant.properties", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return
		}
		_, err = file.WriteString(str)
		if err != nil {
			return
		}
	}
	// use ant release to build and pack the project.
	err = cmd.SyncCmd("ant", []string{"release", "-Dsdk.dir=/usr/lib/android/sdk"})
	if err != nil {
		return
	}

	// remove the ant.properties will avoid the git conflict
	if sign == true {
		err = os.Remove("ant.properties")
		if err != nil {
			return
		}
	}

	// Anaylse the build directory and found .apk file.
	// If sign is set false. find -unsigned.apk
	// if sign is set true. find the .apk without -unsigned and -unaligned
	buildDir, err := os.Open("./bin")
	if err != nil {
		return err
	}
	fileInfos, err := buildDir.Readdir(0)
	if err != nil {
		return err
	}

	var targetApkPath string
	for i := 0; i < len(fileInfos); i++ {
		if fileInfos[i].IsDir() == true {
			continue
		}
		filename := fileInfos[i].Name()
		b0 := apkReg.Match([]byte(filename))
		b1 := unalignedapkReg.Match([]byte(filename))
		if b0 == false || b1 == true {
			continue
		}
		b2 := unsignedapkReg.Match([]byte(filename))

		if sign == b2 {
			continue
		} else {
			var wd string
			wd, err = os.Getwd()
			if err != nil {
				return
			}
			targetApkPath = wd + "/bin/" + filename
			break
		}
	}

	if targetApkPath == "" {
		err = errors.New("No apk found")
	}
	logger.Debug("Find " + targetApkPath)
	// targetApkPath records the absolute path of target apk.
	// change the working dirctory last time and copy the apk file

	logger.Debug("Enter" + workdir)
	err = os.Chdir(workdir) //recover the work directory caused by fetchFromRemote
	if err != nil {
		return
	}
	logger.Debug("copy " + targetApkPath + " to " + outfile)
	cmd.SyncCmd("cp", []string{"-R", targetApkPath, outfile})

	defer func() {
		file.Close()
		buildDir.Close()
	}()
	return
}

// clone or pull the repo from remote
// ATTENTION!!! : The working directory changed after the function is called
func fetchFromRemote(repo string) (err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	homedir := u.HomeDir
	logger.Debug("Home: " + homedir)
	owner, dir, err := getRepoDir(repo)
	if err != nil {
		return
	}
	err = os.MkdirAll(homedir+"/Library/go-pac/"+owner+"/"+dir, 0755)
	if err != nil {
		return
	}
	logger.Debug("Enter ~/Library/go-pac/" + owner + "/" + dir)
	err = os.Chdir(homedir + "/Library/go-pac/" + owner + "/" + dir)
	if err != nil {
		return
	}
	err = cmd.SyncCmd("git", []string{"init"})
	if err != nil {
		return
	}
	err = cmd.SyncCmd("git", []string{"pull", repo})
	if err != nil {
		return
	}
	return nil
}

//return the repo's owner and name.
func getRepoDir(repo string) (string, string, error) {
	strs := strings.Split(repo, "/")
	if len(strs) == 1 {
		return "", "", errors.New("Invalid Repository.")
	}
	return strs[len(strs)-2], strs[len(strs)-1], nil
}
