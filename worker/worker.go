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

var antPropertyTemplate string

func init() {
	antPropertyTemplate = "key.store=%s\n" +
		"key.alias=%s\n" +
		"key.store.password=%s\n" +
		"key.alias.password=%s"
}

func Run(filename string) (err error) {
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
		return compileAndroid(config.Android)
	}
	return
}

func compileAndroid(config *conf.AndroidConfig) (err error) {
	err = cmd.SyncCmd("ant", []string{"clean", "-Dsdk.dir=/usr/lib/android/sdk"})
	if err != nil {
		return
	}
	var sign bool
	if config.Store != nil && config.StorePassword != nil && config.Alias != nil && config.AliasPassword != nil {
		sign = true
	}

	if sign == true {
		str := fmt.Sprintf(antPropertyTemplate, *config.Store, *config.StorePassword, *config.Alias, *config.AliasPassword)
		var file *os.File
		file, err = os.OpenFile("ant.properties", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		defer file.Close()
		if err != nil {
			return
		}
		_, err = file.WriteString(str)
		if err != nil {
			return
		}
	}

	err = cmd.SyncCmd("ant", []string{"release", "-Dsdk.dir=/usr/lib/android/sdk"})
	if err != nil {
		return
	}
	if sign == true {
		err = os.Remove("ant.properties")
		if err != nil {
			return
		}
	}
	return
}

// clone or pull the repo from remote
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
