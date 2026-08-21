package slotsmongo

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"serve/comm/db"
	"serve/comm/lazy"

	"github.com/samber/lo"
)

const (
	BaseDumpDir = "/data/pggames/dump"
)

func commandRun(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	//  cmd.String()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("commandRun start", "cmdline", cmd.String())
	err := cmd.Run()
	slog.Info("commandRun end", "cmdline", cmd.String(), "error", err)
	return err
}

func restoreSpinData() (err error) {
	// return nil
	// db.Collection("combine")

	backupHost, _ := lazy.RouteFile.Get("backuphost")
	if backupHost == "" {
		return nil
	}

	settingfile := fmt.Sprintf("/data/pggames/bin/config/%s_setting.yaml", lazy.ServiceName)
	_, err = os.Stat(settingfile)
	if os.IsNotExist(err) && backupHost != "" {
		err = commandRun("scp", backupHost+":"+settingfile, settingfile)
	}
	if err != nil {
		return
	}

	dumpdir := path.Join(BaseDumpDir, lazy.ServiceName)

	dirs, err := os.ReadDir(dumpdir)

	if os.IsNotExist(err) && backupHost != "" {

		os.MkdirAll(BaseDumpDir, 0755)

		// rsync -avc doudou-prod:/data/game-history /data/
		// rsync -avP root@doudou-prod:/data/pggames/dump .
		err = commandRun("rsync", "-avP", backupHost+":"+dumpdir, BaseDumpDir)

		if err != nil {
			return
		}

		dirs, err = os.ReadDir(dumpdir)
	}

	if err != nil {
		return
	}

	// combine.bson.gz  combine.metadata.json.gz  pgSpinData.bson.gz  pgSpinData.metadata.json.gz
	for _, v := range dirs {
		if !v.Type().IsRegular() {
			continue
		}

		name := v.Name()
		collname := strings.TrimSuffix(name, ".bson.gz")
		if collname == name {
			continue
		}

		count, er := db.Collection(collname).EstimatedDocumentCount(context.TODO())
		if er != nil {
			err = er
			return
		}
		if count == 0 {
			// cmd := exec.Command("mongorestore", "-d", lazy.ServiceName, "-c", collname, dumpdir)

			nsIn := fmt.Sprintf("--nsInclude=%s.%s", lazy.ServiceName, collname)
			mgoaddr := lo.Must(lazy.RouteFile.Get("mongo"))
			err = commandRun("mongorestore", "--gzip", nsIn, mgoaddr, path.Join(dumpdir, name))

			if err != nil {
				return
			}
		}

	}

	// os.Exit(0)
	return
}

// scp Not Found /data/pggames/bin/config/pg_39_setting.yaml
// 从文件中恢复 spin 数据
func RestoreSpinData() {
	return
	lo.Must0(restoreSpinData()) //需要查找/data/pggames/bin/config/%s_setting.yaml下对应的文件，暂时没有，跳过
}

func pathIsExist(pth string) bool {
	_, err := os.Stat(pth)
	lo.Must0(err == nil || os.IsNotExist(err))
	return err == nil
}

func restoreClientCache(basedir, dir string) (err error) {
	// /data/game/bin/cache/uat-wbgame.jlfafafa3.com/csh
	backupHost, _ := lazy.RouteFile.Get("backuphost")
	if backupHost == "" {
		return nil
	}

	rdir := path.Join(basedir, dir)

	if pathIsExist(rdir) {
		slog.Info("pathIsExist", "rdir", rdir)
		return
	}

	// rsync -av backuphost:/data/game/bin/cache/uat-wbgame.jlfafafa3.com/csh
	// dir = path.Clean(dir)

	err = commandRun("rsync", "-av", backupHost+":"+rdir, basedir)
	return
}

// func RestoreClientCache(dir string) {
// 	lo.Must0(restoreClientCache(dir))
// }

func RestoreClientCacheJILI(shortname string) {
	// dir := path.Join("")
	lo.Must0(restoreClientCache("/data/game/bin/cache/wbgame.bd33fgabh.com", shortname))

}

// func RestoreClientCachePG(shortname string) {
// dir := path.Join("")
// lo.Must0(restoreClientCache("/data/game/bin/cache/uat-wbgame.jlfafafa3.com", shortname))
// }
